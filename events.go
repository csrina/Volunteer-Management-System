package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"strings"

	"github.com/gorilla/mux"
)

type Response struct {
	Msg      string `json:"msg"`      // message
	UID      int    `json:"userId"`   // users ID for adding to bookings listing client side
	UserName string `json:"userName"` // user's full name for adding to bookings listing client side
	Colour   string `json:"color"`    // color of room
	BID      int    `json:"bookId"`   // booking id
	ID       int    `json:"id"`       // block/event id
}

func (r *Response) setMsg(msg string) *Response {
	r.Msg = msg
	return r
}

func (r *Response) setBID(bid int) *Response {
	r.BID = bid
	return r
}

func (r *Response) setID(bid int) *Response {
	r.ID = bid
	return r
}

func (r *Response) send(w http.ResponseWriter) {
	enc := json.NewEncoder(w) // encoder for sending response
	enc.Encode(r)
}

// Expected time formats from calendar
const (
	isoTimeShort = "2006-01-02"
	isoTimeFull  = "2006-01-02T15:04:05"
)

/*
 * Performs auth checks, analyzes request and directs appropriately
 */
func eventPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	role, err := getRoleNum(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dest := mux.Vars(r)["target"] // Determine POST destination from URL
	if dest == "book" {           // add/remove booking corresponding to event data
		bookingHandler(w, r, role)
	} else if role == ADMIN {
		adminPostHandler(w, r, dest)
	} else if role == TEACHER {
		teacherPostHandler(w, r, dest)
	}
}

/* not implemented yet */
func teacherPostHandler(w http.ResponseWriter, r *http.Request, dest string) {
	e, err := EventFromJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var response *Response
	switch dest {
	case "update":
		// check for room_name
		response, err = e.update()
	default:
		err = errors.New("invalid target")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.send(w)
	return
}

/*
 * Handles posts which are not bookings (ergo, must be admin.)
 */
func adminPostHandler(w http.ResponseWriter, r *http.Request, dest string) {
	e, err := EventFromJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response *Response
	switch dest {
	case "add":
		response, err = e.add()
	case "update":
		response, err = e.update()
	case "delete":
		response, err = e.delete()
	default:
		err = errors.New("invalid target")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.send(w)
}

func mapEventJSON(r *http.Request) (map[string]interface{}, error) {
	/* Read the json */
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// need to unmarshall -> map prior to selectively adding fields to an empty booking
	var evInterface interface{}
	json.Unmarshal(body, &evInterface)
	ev, ok := evInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("couldn't map request")
	}
	logger.Println("Ev Map: ", ev)
	return ev, nil
}

// An Event is a time block + a booking array + other details needed by calendar
type Event struct {
	// Block info + a title (required field for calendar)
	ID     int       `db:"block_id" json:"id"`
	Title  string    `db:"note" json:"title"`
	Start  time.Time `db:"block_start" json:"start"`
	End    time.Time `db:"block_end" json:"end"`
	RoomID int       `db:"room_id" json:"roomId"`
	Room   string    `db:"room_name" json:"room"` // fullCalendar will make blocks colour of room
	Colour string    `json:"color"`               // color code for event rendering (corresponds to the room name)
	// booking ids for lookup
	Capacity     int         `json:"capacity"`
	BookingCount int         `json:"bookingCount"`
	Booked       bool        `json:"booked"`
	Bookings     []UserShort `json:"bookings"` // we store their actual names in the username field
	// description
	Modifier int    `db:"modifier" json:"modifier"`
	Note     string `json:"note"`
}

/*
 * Deletes the underlying timeblock db entry of the reciever event.
 */
func (e *Event) delete() (*Response, error) {
	tb, err := getTimeBlockByID(e.ID)
	if err != nil {
		return nil, err
	}
	err = tb.delete()
	if err != nil {
		return nil, err
	}
	return &Response{Msg: "Event removed"}, nil
}

/*
 * Update the underlying timeblock db entry of a reciever event (e)
 * to reflect the start and end times of the reciever event.
 */
func (e *Event) update() (*Response, error) {
	tb, err := getTimeBlockByID(e.ID)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	tb.Start = e.Start
	tb.End = e.End
	tb.Note = e.Note
	tb.Title = e.Title
	tb.Modifier = e.Modifier
	tb.Capacity = e.Capacity

	err = tb.update()
	if err != nil {
		return nil, err
	}
	return &Response{Msg: "Successfully updated time block"}, nil
}

/*
 * Add a new time block to DB, for the event reciever.
 * Return a response struct or error
 */
func (e *Event) add() (*Response, error) {
	tb := e.getTimeBlock()
	var err error
	e.ID, err = tb.insert() // insert block & set e.ID to correspond to the tbID returned
	if err != nil {
		return nil, err
	}
	e.updateColourCode() // update color of event
	// create and return response
	response := new(Response)
	response.Msg = "Successfully added time block"
	response.Colour = e.Colour
	response.setID(tb.ID)
	return response, nil
}

func EventFromJSON(r *http.Request) (*Event, error) {
	e := new(Event)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(e)
	return e, err
}

func (e *Event) getTimeBlock() *TimeBlock {
	tb := new(TimeBlock)
	tb.ID = e.ID
	tb.Start = e.Start
	tb.End = e.End
	tb.Note = e.Note
	tb.Modifier = e.Modifier
	tb.Room = e.RoomID
	if e.Title == "" {
		e.Title = "Facilitation"
	}
	tb.Title = e.Title
	tb.Capacity = e.Capacity

	return tb
}

/* Lists the events requested */
func listEvents(r *http.Request) ([]*Event, error) {
	/* obtain the blockz in range */
	params := r.URL.Query() // Get the params from url as a {key : value} string map

	start := params.Get("start")
	end := params.Get("end")
	if start == "" || end == "" {
		return nil, errors.New("date(s) couldn't be resolved")
	}
	logger.Println("Start Date: " + start + "\tEnd Date: " + end)
	/* Get time blocks in range */
	blocks, err1 := getBlocksWithMoments(start, end)
	if err1 != nil {
		panic(err1)
	}

	/* Make the eventz */
	uid := getUID(r)
	if uid < 0 {
		return nil, errors.New("uid unresolved")
	}

	evs := makeEvents(blocks)
	for _, e := range evs {
		e.setBookingStatus(uid)
	}

	/* If target given, and target is dash --> only return events for which the user is booked */
	dest := mux.Vars(r)["target"]
	if dest == "dash" {
		var evs2 []*Event
		for _, e := range evs {
			if e.Booked {
				evs2 = append(evs2, e) // add to evs2 if booked
			}
		}
		return evs2, nil // return booked events only
	}
	return evs, nil // return all events in range
}

/* Using a url encoded params, responds with a json event stream */
func getEvents(w http.ResponseWriter, r *http.Request) {
	evs, err := listEvents(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//context.Set(r, "error", http.StatusBadRequest)
		return
	}

	if role, _ := getRoleNum(r); role == TEACHER {
		roomsTaught, err := getTaughtRooms(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		for _, e := range evs {
			for _, rm := range roomsTaught {
				if rm == e.Room {
					e.Booked = true // set flag (we repurpose the booked flag to mean 'teaches'
					break
				}
			}
		}
	}

	serveEventJSON(w, evs)
}

/*
 * Serves the event json stream via the io writer (generic io writer need for testing)
 */
func serveEventJSON(w io.Writer, events []*Event) {
	/* Create and serve JSON*/
	eventJSON, err := json.Marshal(events)
	if err != nil {
		logger.Println(err)
		return
	} else {
		w.Write(eventJSON)
	}
}

/*
 * Creates Events from slice of time blocks
 */
func makeEvents(blocks []TimeBlock) []*Event {
	/* Get events from blocks */
	var events []*Event
	for _, b := range blocks {
		events = append(events, NewEvent(&b))
	}
	return events
}

/* NewEvent creates and returns ptr to a corresponding event representation */
func NewEvent(b *TimeBlock) *Event {
	/* Init w/ directly transferable properties */
	e := &Event{
		ID:       b.ID,
		Start:    b.Start,
		End:      b.End,
		Note:     b.Note,
		Title:    b.Title,
		Modifier: b.Modifier,
		Capacity: b.Capacity,
	}
	err := e.initBookings()
	if err != nil {
		logger.Println(err)
	}
	/* Get room and bookings for the Event */
	err = db.QueryRow(`SELECT room_name FROM room WHERE $1 = room_id`, b.Room).Scan(&e.Room) // Get the room name
	if err != nil {
		logger.Println(err)
	}
	e.updateColourCode()
	return e
}

/*
 * Initializes the bookings list from the DB
 *
 * requires the eventID be set to a valid blockID
 */
func (e *Event) initBookings() error {
	tb := e.getTimeBlock()
	bookings, err := tb.bookings()
	if err != nil {
		if err != sql.ErrNoRows {
			e.BookingCount = 0
			return nil
		}
		return err
	}

	e.BookingCount = len(bookings)
	for _, b := range bookings {
		u := new(UserShort)
		u.UserID = b.UserID
		u.UserName, err = u.getFullName()
		if err != nil {
			logger.Println(err)
			return err
		}
		e.Bookings = append(e.Bookings, *u)
	}
	return nil
}

func (e *Event) setBookingStatus(uid int) (*Event, error) {
	/* Set the booked flag based on the provided uid */
	booked := 0
	q := `SELECT count(*) FROM booking WHERE (block_id = $1 AND user_id = $2)`
	db.QueryRow(q, e.ID, uid).Scan(&booked)
	if booked == 0 {
		e.Booked = false
	} else {
		e.Booked = true
	}
	/* Get the booked count */
	q = `SELECT count(*) FROM booking WHERE block_id = $1`
	err := db.QueryRow(q, e.ID).Scan(&e.BookingCount)
	if err != nil {
		logger.Println(err)
		return e, err
	}
	return e, nil
}

/* Prety coloour plalalalette */
//noinspection GoUnusedConst,GoUnusedConst,GoUnusedConst
const (
	RED      = "#F44336"
	PINK     = "#E91E63"
	PURPLE   = "#9C27B0"
	BLUE     = "#2196F3"
	DGREEN   = "#4CAF50"
	LGREEN   = "#76FF03"
	LIME     = "#AEEA00"
	YELLOW   = "#FAD201"
	ORANGE   = "#FF9800"
	GREY     = "#9E9E9E"
	BLUEGREY = "#607D8B"
)

/* CHanges the color code to correspond to the room name of the event */
func (e *Event) updateColourCode() {
	switch e.Room {
	case "red":
		e.Colour = RED
	case "pink":
		e.Colour = PINK
	case "purple":
		e.Colour = PURPLE
	case "grey":
		e.Colour = GREY
	case "blue":
		e.Colour = BLUE
	case "green":
		e.Colour = DGREEN
	case "orange":
		e.Colour = ORANGE
	case "yellow":
		e.Colour = YELLOW
	default:
		e.Colour = strings.Split(e.Room, " ")[0] // take only first string to avoid breaking the world
	}
}

