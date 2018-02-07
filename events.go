package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Contains the implementation details for events json streaming

// Models record from Booking table in DB
type Booking struct {
	block_id int
}

// An Event is a time block + a booking array + other details needed by calendar
type Event struct {
	// Block info + a title (required field for calendar)
	ID    int       `db:"block_id" json:"id"`
	Title string    `db:"note" json:"title"`
	Start time.Time `db:"block_start" json:"start"`
	End   time.Time `db:"block_end" json:"end"`
	Room  string    `db:"room_name" json:"color"` // fullCalendar will make blocks colour of room
	// bookings data
	Bookings []Booking `json:"bookings"`
	// description
	Note []string `json:"note"`
}

// Expected time formats from calendar
const (
	iso_time_short = "2006-01-02"
	iso_time_full  = "2006-01-02T15:04:05-0700"
)

// Using url encoded params, responds with a json event stream
func getEvents(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query() // Get the params from url as a {key : value} string map
	start, err1 := parseDate(params["start"][0])
	end, err2 := parseDate(params["end"][0])
	if err1 != nil || err2 != nil {
		logger.Fatal("Could not parse dates")
	}

	logger.Println("Start Date: " + start.String() + "\nEnd Date: " + end.String())
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
	}
	/* Add note to event */
	e.Note = append(e.Note, b.Note)
	/* Get room and bookings for the Event */
	db.Select(e, `SELECT room_name FROM room WHERE $1 = room_id`, b.Room) // Get the room name
	db.Select(e, `SELECT count(*) FROM booking WHERE $1 = block_id`, b.ID)
	/* Set title */
	e.Title = e.Room + "Facilitation"

	/* Debug logging */
	logger.Println("New event created: ", e)
	logger.Println("From block: ", b)

	return e
}
