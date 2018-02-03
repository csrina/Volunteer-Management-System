package main

import (
	"time"

	_ "github.com/lib/pq"
)

/*
 * Corresponds to a tuple in the time_block table.
 */
type TimeBlock struct {
	Id       int       `db:"block_id"`
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
	q := `INSERT INTO time_block (block_id, block_start, block_end, room, modifier, note)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING block_id`
	id, err := db.Exec(q, tb.Id, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Note)
	if err {
		return err
	}
	tb.Id = id
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

	_, err := db.Exec(q, tb.Id, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Note)

	if err {
		return err
	}
}

/*
 * Retrieve records for block(s) from table in range (min inclusive, max exclusive).
 */
func getBlocks(start time.Time, end time.Time) ([]TimeBlock, error) {

}
