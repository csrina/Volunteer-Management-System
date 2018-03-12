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

// An Event is a time block + a booking array + other details needed by calendar
type Event struct {
	// Block info + a title (required field for calendar)
	ID     int       `db:"block_id" json:"id"`
	Title  string    `db:"note" json:"title"`
	Start  time.Time `db:"block_start" json:"start"`
	End    time.Time `db:"block_end" json:"end"`
	Room   string    `db:"room_name" json:"room"` // fullCalendar will make blocks colour of room
	Colour string    `json:"color"`               // color code for event rendering (corresponds to the room name)
	// booking ids for lookup
	BookingCount int  `json:"bookingCount"`
	Booked       bool `json:"booked"`
	// description
	Note string `json:"note"`
}

type Response struct {
	Msg string `json:"msg"`
	BID int    `json:"bookId"`
}

func (r *Response) setMsg(msg string) *Response {
	r.Msg = msg
	return r
}

func (r *Response) setBID(bid int) *Response {
	r.BID = bid
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
	booking, err := bookingFromJSON(r)
	if err != nil {
		logger.Println(err)
		http.Error(w, "could not resolve booking data, contact your administrator if the problem persists", http.StatusBadRequest)
		return
	}
	response := new(Response)     // Response data
	dest := mux.Vars(r)["target"] // Determine POST destination from URL
	if dest == "book" {
		// book/unbook based on BookingID (if 0 value, then booking doesn't exist yet and must be created)
		if booking.BookingID > 0 {
			err = booking.unbook(role)
			response.Msg = "Successfully cancelled booking"
		} else {
			err = booking.book(role)
			response.Msg = "Successfully created booking"
		}
	} else if dest == "update" {
		err = booking.updateTimeBlock(role)
		response.Msg = "Successfully updated time block"
	} else {
		// Invalid target
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// We may have an error condition after book/unbook
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.setBID(booking.BookingID).send(w)
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

/*
 * Reads json request, and creates a partially filled booking struct
 */
func bookingFromJSON(r *http.Request) (*Booking, error) {
	booking := new(Booking)
	uid := getUID(r)
	if uid < 1 {
		return nil, errors.New("uid unresolved")
	}
	booking.UserID = uid
	// Now extract data from json request body
	ev, err := mapEventJSON(r)
	if err != nil {
		return nil, err
	}
	eid, ok := ev["id"].(float64)
	if !ok {
		return nil, errors.New("event id NaN")
	}
	booking.BlockID = int(eid)
	booking.FamilyID, err = getUsersFID(booking.UserID)
	if err != nil {
		return nil, errors.New("couldn't resolve FID from UID")
	}
	booking.getTimesMap() // Will set our Start/End for the event with DB values
	// get the bookingID (if this booking doesn't exist, the ID will remain as 0
	booking.BookingID, err = booking.getBookingID()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	// We dont need roomID nor modifer to make a bookingBlock entry
	return booking, nil
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
	if strings.ContainsAny(start, ";") || strings.ContainsAny(end, ";") {
		return nil, errors.New("';' in date, scary")
	}
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
	serveEventJSON(w, evs)
}

/*
 * Serves the event json stream via the io writer (generic io writer need for testing)
 */
func serveEventJSON(w io.Writer, events []*Event) {
	/* Create and serve JSON*/
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	enc.Encode(events)
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
	logger.Println(events)
	return events
}

/*
 * Parses the ISO6801 date string passed in url by fullCalendar.
 * Attempts to get long-form, then short-form;
 * upon failure it returns the current datetime and the parsing error
 */
func _(date string) (time.Time, error) {
	// try parsing long-form

	d, err := time.Parse(isoTimeShort, date)
	if err == nil {
		return d, nil
	}
	// try short form
	d, err = time.Parse(isoTimeFull, date)
	if err == nil {
		return d, nil
	}
	return time.Now(), err
}

/* NewEvent creates and returns ptr to a corresponding event representation */
func NewEvent(b *TimeBlock) *Event {
	/* Init w/ directly transferable properties */
	e := &Event{
		ID:    b.ID,
		Start: b.Start,
		End:   b.End,
		Note:  b.Note,
		Title: "Facilitation",
	}
	/* Get room and bookings for the Event */
	err := db.QueryRow(`SELECT room_name FROM room WHERE $1 = room_id`, b.Room).Scan(&e.Room) // Get the room name
	if err != nil {
		logger.Println(err)
	}
	e.updateColourCode()
	return e
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
		e.Colour = LGREEN
	case "orange":
		e.Colour = ORANGE
	case "yellow":
		e.Colour = YELLOW
	default:
		e.Colour = e.Room
	}
}
