package main

import (
	"errors"
	"time"

	_ "github.com/lib/pq"
	"database/sql"
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
	Note     string    `db:"note" json:"note"`
}

func getTimeBlockByID(id int) (*TimeBlock, error) {
	tb := new(TimeBlock)
	q := `SELECT * FROM time_block WHERE time_block.block_id = $1`
	err := db.QueryRow(q, id).Scan(&tb.ID, &tb.Start, &tb.End, &tb.Room, &tb.Modifier, &tb.Note)
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

	q := `INSERT INTO time_block (block_start, block_end, room_id, modifier, note)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING block_id`
	err := db.QueryRow(q, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Note).Scan(&tb.ID)
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
			SET(block_id, block_start, block_end, room_id, modifier, note)
			= ($1, $2, $3, $4, $5, $6)
		WHERE (time_block.block_id = $1)`

	_, err := db.Exec(q, tb.ID, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Note)
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
