package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

/* TimeBlock ...
 * Corresponds to a tuple in the time_block table.
 */
type TimeBlock struct {
	ID       int       `db:"block_id" json:"blockID"`
	Start    time.Time `db:"block_start" json:"start"`
	End      time.Time `db:"block_end" json:"end"`
	Room     int       `db:"room_id" json:"room"`
	Modifier int       `db:"modifier" json:"modifier"`
	Capacity int       `db:"capacity" json:"capacity"`
	Title    string    `db:"title" json:"title"`
	Note     string    `db:"note" json:"note"`
}

func getTimeBlockByID(id int) (*TimeBlock, error) {
	tb := new(TimeBlock)
	q := `SELECT * FROM time_block WHERE time_block.block_id = $1`
	err := db.QueryRow(q, id).Scan(&tb.ID, &tb.Start, &tb.End, &tb.Room, &tb.Modifier, &tb.Title, &tb.Note)
	return tb, err
}

/*
 * Returns the bookings associated with a given timeblock
 */
func (tb *TimeBlock) bookings() ([]bookingBlock, error) {
	q := `SELECT booking_id FROM booking WHERE booking.block_id = $1`
	var bids []int
	err := db.Select(&bids, q, tb.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Dont log this but pass it up chain incase its important to know
		}
		logger.Println(err)
		return nil, err
	}
	// Initialize the booking structs using the ids and return them in a slice
	var bookings []bookingBlock
	for _, bid := range bids {
		b := new(bookingBlock)
		err = b.init(bid)
		if err != nil {
			logger.Println("bookingBlock creation failed: ", err)
			return nil, err
		}
		bookings = append(bookings, *b)
	}
	return bookings, nil
}

/*
 * Inserts tb into the database, returns the id (and actively sets the structs ID field in place)
 */
func (tb *TimeBlock) insert() (int, error) {

	q := `INSERT INTO time_block (block_start, block_end, room_id, modifier, capacity, title, note)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING block_id`
	err := db.QueryRow(q, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Capacity, tb.Title, tb.Note).Scan(&tb.ID)
	if err != nil {
		return -1, err
	}
	return tb.ID, nil

}

/*
 * Saves the state of the block (tb)
 * to the db. Where tb is an existing block in the db
 */
func (tb *TimeBlock) update() error {
	q := `UPDATE time_block
			SET(block_start, block_end, room_id, modifier, capacity, title, note)
			= ($2, $3, $4, $5, $6, $7, $8)
		WHERE (time_block.block_id = $1)`

	_, err := db.Exec(q, tb.ID, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Capacity, tb.Title, tb.Note)
	return err
}

/*
 * Deletes the recieving block's (tb) entry in the db.
 */
func (tb *TimeBlock) delete() error {
	q := `DELETE FROM time_block WHERE time_block.block_id = $1`
	results, err := db.Exec(q, tb.ID)
	if err != nil {
		logger.Println(err)
		return err
	}
	count, err := results.RowsAffected()
	if err != nil {
		logger.Println(err)
		return err
	} else if count != 1 {
		err = errors.New("time block not deleted")
		logger.Println(err)
		return err
	}
	return nil
}

/*
 * Retrieve records for block(s) from table in range (start inclusive, end exclusive).
 */
func getBlocks(start time.Time, end time.Time) ([]TimeBlock, error) {
	// Retrieve blocks w/in date range
	q := `SELECT * FROM time_block
			   WHERE block_start >= $1 AND block_end <= $2`
	var blocks []TimeBlock
	err := db.Select(&blocks, q, start, end)

	if err != nil {
		return nil, err
	} else if len(blocks) > 0 {
		return blocks, nil
	} else {
		return nil, nil
	}
}

func getBlocksWithMoments(start string, end string) ([]TimeBlock, error) {
	// Retrieve blocks w/in date range
	q := `SELECT * FROM time_block
			   WHERE block_start >= $1 AND block_end <= $2`
	var blocks []TimeBlock
	err := db.Select(&blocks, q, start, end)

	if err != nil {
		return nil, err
	} else if len(blocks) > 0 {
		return blocks, nil
	} else {
		return nil, nil
	}
}

/* Setter for day field of timeblock */
func (tb *TimeBlock) setDay(startDate time.Time) {
	tb.Start = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), tb.Start.Hour(), tb.Start.Minute(), 0, 0, tb.Start.Location())
	tb.End = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), tb.End.Hour(), tb.End.Minute(), 0, 0, tb.Start.Location())
}

/* For templating calendar aka schedule builder mode */
const (
	WEEKLY = iota
	MONTHLY
)

type IntervalData struct {
	Repeats int // See constant declarations above
	Delta   int // e.g. 2 ==> every 2nd week/month, 1 == every week/month, 3 == every 3rd week/month
	// SubDelta may b eused for monthly repeaters
	SubDeltas []int // E.g. Every 2nd monday (on a 2nd month Delta) --> 2nd monday of every 2nd month;; is a slice to allow for things like (1, 3) --> 1st and 3rd Day of DeltaREpeatingMonth
}

/*
 * Events from templater have interval fields included
 */
type BuilderEvent struct {
	Event
	Interval IntervalData
}

// Advance by major delta (e.g. next month)
func (be *BuilderEvent) Increment() (be *BuilderEvent) {
	switch be.Interval.Repeats {
	case WEEKLY:
		be.Start = be.Start.AddDate(0, 0, 7 * be.Interval.Delta)
		be.End = be.End.AddDate(0, 0, 7 * be.Interval.Delta)
	case MONTHLY:
		be.Start = be.Start.Add(0, be.IntervalDelta, 0)
		be.End = be.End.Add(0, be.IntervalDelta, 0)
	}
}

/*
 * List of builderEvents for creating schedule
 * Period start = inclusive start date for period to build
 * Period end is exclusive moment for period to build in
 */
type ScheduleBuilderData struct {
	events      []BuilderEvent `json:"events"`
	periodStart time.Time      `json:"periodStart"`
	periodEnd   time.Time      `json:"periodEnd"`
}

// For handling template build/destroy/whateverelse POSTS
func schedulePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	role, err := getRoleNum(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if role != ADMIN {
		http.Error(w, "You must be an admin to do this action (your auth cookie may have expired)", http.StatusForbidden)
		return
	}

	dest := mux.Vars(r)["target"] // Determine POST destination from URL
	switch dest {
	case "build":
		buildRequestHandler(w, r)
		return
	case "destroy":
		// If theres time, maybe a takedown mode I dunno
		http.Error(w, "Not implemented yet", http.StatusBadGateway)
		return
	default:
		http.Error(w, "Invalid destination specified", http.StatusBadGateway)
	}
}

func buildRequestHandler(w http.ResponseWriter, r *http.Request) {
	/* Read the json */
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	builderData := new(ScheduleBuilderData)
	json.Unmarshal(body, builderData)
	bdRoomMap := make(map[string][]BuilderEvent)
	for _, e := range builderData.events {
		// only days are relevant coming from the builder (set the year/month to the start of the period
		e.Start = time.Date(builderData.periodStart.Year(), builderData.periodStart.Month(), e.Start.Day(), e.Start.Hour(), e.Start.Minute(), e.Start.Second(), 0, time.Local)
		e.End = time.Date(builderData.periodStart.Year(), builderData.periodStart.Month(), e.End.Day(), e.End.Hour(), e.End.Minute(), e.End.Second(), 0, time.Local)
		bdRoomMap[e.Room]= append(bdRoomMap[e.Room], e) // Split up by rooms (can be parallelized)
	}

	/*
	 * Building schedule is complex; quadratic best case, most likely cubic or worse on average.
	 * We use go routines for each room, as their blocks are independent, so can be parallelized easily.
	 * Additionally, we know adding events to DB will have significant I/O time, so we can swap in other room's routines
	 * rather than wait for each write before proceeding. Going over board with the routines would be bad in otherways,
	 * gut feeling says rooms should be sufficient as a divisor.
	 */
	for _, evs := range bdRoomMap { // ignore roomnames [key is blanked (i.e. _)]
		go BuildSchedule(evs, builderData.periodStart, builderData.periodEnd)
	}
}

/* Creates blocks to satisfy the interval conditions  for each bd in evs, within the period indicated by sop/eop (startofendofperiod */
func BuildSchedule(evs []BuilderEvent, sop, eop time.Time) {
	for _, be := range evs { // foreach be => BuilderEvent given via template calendar
		// Before start of period (period starts  mid week for example), add weeks until in ranger
		for be.Start.Before(sop) {
			// dun dun dun duh duuuuuuuuun de da da dun de da dooo da da bump ba bum bah da (*james bond theme plays*)
			be.Start.AddDate(0,0,7)
			be.End.AddDate(0,0,7)
		}

		for be.Start.Before(eop) { // event in range --> add it to calendar and increment
			switch be.Interval.Repeats {
			case MONTHLY:
				for _, i := range be.Interval.SubDeltas { // For each sub delta value --> add events to calendar for the month
					subEv := be
					subEv.Start.AddDate(0, 0, 7 * i)         // add subDelta to start/end
					subEv.End.AddDate(0, 0, 7 * i)
					if subEv.Start.Before(eop) {             // if still in range add, else continue (no guarentee list is sorted, must keep going!)
						subEv.add()           // Add the event if still in range
					}
				}
			case WEEKLY:
				be.add() // Add to db
			}
			be.Increment() // increment by Delta
		}
	}
}

