package main

import (
	"time"

	_ "github.com/lib/pq"
)

/* TimeBlock ...
 * Corresponds to a tuple in the time_block table.
 */
type TimeBlock struct {
	ID       int       `db:"block_id"`
	Start    time.Time `db:"block_start"`
	End      time.Time `db:"block_end"`
	Room     int       `db:"room_id"`
	Modifier int       `db:"modifier"`
	Note     []string  `db:"note"`
}

/*
 * Inserts tb into the database.
 */
func (tb *TimeBlock) insertBlock() error {
	q := `INSERT INTO time_block (block_start, block_end, room, modifier, note)
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
			SET block_id = $1
			SET block_start = $2
			SET block_end = $3
			SET room_id = $4
			SET modifier = $5
			SET note = $6
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
			   WHERE ($1 <= block_start AND $2 > block_end)`

	blocks := []TimeBlock{}
	err := db.Select(&blocks, q, start, end)

	if err != nil {
		return nil, err
	} else if len(blocks) > 0 {
		return blocks, nil
	} else {
		return nil, nil
	}
}
