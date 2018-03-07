package main

import (
	"time"

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
	Note     string    `db:"note" json:"note"`
}

/*
 * Inserts tb into the database.
 */
func (tb *TimeBlock) insertBlock() error {
	q := `INSERT INTO time_block (block_start, block_end, room_id, modifier, note)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING block_id`
	err := db.QueryRow(q, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Note).Scan(&tb.ID)
	if err != nil {
		return err
	}
	return nil
}

/*
 * Saves the state of the block (tb)
 * to the db. Where tb is an existing block in the db
 */
func (tb *TimeBlock) updateBlock() error {
	q := `UPDATE time_block
			SET(block_id, block_start, block_end, room_id, modifier, note)
			= ($1, $2, $3, $4, $5, $6)
		WHERE (time_block.block_id = $1)`

	_, err := db.Exec(q, tb.ID, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Note)

	return err
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

/* Setter for day field of timeblock */
func (tb *TimeBlock) setDay(startDate time.Time) {
	tb.Start = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), tb.Start.Hour(), tb.Start.Minute(), 0, 0, tb.Start.Location())
	tb.End = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), tb.End.Hour(), tb.End.Minute(), 0, 0, tb.Start.Location())
}
