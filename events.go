package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// An Event is a time block + a booking array + other details needed by calendar
type Event struct {
	// Block info + a title (required field for calendar)
	ID       int       `db:"block_id" json:"id"`
	Title    string    `db:"note" json:"title"`
	Start    time.Time `db:"block_start" json:"start"`
	End      time.Time `db:"block_end" json:"end"`
	Room     string    `db:"room_name" json:"color"` // fullCalendar will make blocks colour of room
	Modifier int       `db:"Modifier" json:"value"`
	// booking ids for lookup
	BookingCount int `json:"bookings"`
	// description
	Note string `json:"note"`
}

// Expected time formats from calendar
const (
	iso_time_short = "2006-01-02"
	iso_time_full  = "2006-01-02T15:04:05"
)

/*
 * Update the time block of an event upon a POST request from calendar
 * --> NEW DURATION AND/OR START/END TIMES WILL BE UPDATED
 * TODO: Add modifier and Note to be updated if changed as well
 */
func updateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	/* Read the json posted from calendar */
	body, _ := ioutil.ReadAll(r.Body)
	logger.Println(string(body))

	/*
	 * Unmarshal json into string:string map intermediate.
	 * 		We require the intermediate, because go doesn't
	 *      parse the timestamps correctly. However, psql
	 *      parses the date string just fine, we don't even
	 *      need to use parseDate to insert!
	 */
	var evInterface interface{}
	json.Unmarshal(body, &evInterface)

	ev := evInterface.(map[string]interface{})
	logger.Println(ev)

	id, ok := ev["id"].(float64)
	if ok != true {
		logger.Println("Invalid id given")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	/* Update DB reference */
	q := `UPDATE time_block 
			SET (block_start, block_end) = ($2, $3)
			WHERE (block_id = $1)`
	_, err := db.Exec(q, int(id), ev["start"], ev["end"])
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Write(body)
	}
}

// Using url encoded params, responds with a json event stream
func getEvents(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query() // Get the params from url as a {key : value} string map
	start, err1 := parseDate(params.Get("start"))
	end, err2 := parseDate(params.Get("end"))
	if err1 != nil || err2 != nil {
		logger.Println("Could not parse dates")
		w.WriteHeader(http.StatusBadRequest)
		//context.Set(r, "error", http.StatusBadRequest)
		return
	}

	logger.Println("Start Date: " + start.String() + "\tEnd Date: " + end.String())
	/* Get time blocks in range */
	blocks, err1 := getBlocks(start, end)
	if err1 != nil {
		panic(err1)
	}
	serveEventJSON(w, makeEvents(blocks))
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
func parseDate(date string) (time.Time, error) {
	// try parsing long-form
	d, err := time.Parse(iso_time_full, date)
	if err == nil {
		return d, nil
	}
	// try short form
	d, err = time.Parse(iso_time_short, date)
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
	err = db.QueryRow(`SELECT count(*) FROM booking WHERE $1 = block_id`, b.ID).Scan(&e.BookingCount)
	if err != nil {
		logger.Println(err.Error())
	}
	logger.Println("COUNT: ", e.BookingCount)

	/* Debug logging */
	//logger.Println("New event created: ", e.Title)
	//logger.Println("From block: ", b)

	return e
}
