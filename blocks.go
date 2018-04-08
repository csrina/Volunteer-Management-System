package main

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/jmoiron/sqlx"
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
	Title    string    `db:"title" json:"title"`
	Note     string    `db:"note" json:"note"`
}

type MsgWording struct {
	Start time.Time `db:"block_start" json:"start"`
	End   time.Time `db:"block_end" json:"end"`
	Title string    `db:"title" json:"title"`
	Room  string    `db:"room_name" json:"roomname"`
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

	q := `INSERT INTO time_block (block_start, block_end, room_id, modifier, title, note)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING block_id`
	err := db.QueryRow(q, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Title, tb.Note).Scan(&tb.ID)
	if err != nil {
		return -1, err
	}
	return tb.ID, nil

}

func getOldTimeBlockValues(tx *sqlx.Tx, tid int) (MsgWording, error) {
	var msg MsgWording
	q := "select time_block.block_start, time_block.block_end, time_block.title, room.room_name from time_block, room where time_block.block_id = $1 AND room.room_id = time_block.room_id"
	stmt, err := tx.Preparex(q)
	if err != nil {
		logger.Println(err)
		return msg, err
	}
	defer stmt.Close()

	if err := stmt.Get(&msg, tid); err != nil {
		tx.Rollback()
		logger.Println(err)
		return msg, err
	}
	return msg, nil
}

func getUsersInTimeBlockValues(tx *sqlx.Tx, tid int) ([]int, error) {
	var usrs []int
	logger.Println(tid)
	q := "SELECT user_id FROM booking WHERE block_id = $1"
	stmt, err := tx.Preparex(q)
	if err != nil {
		logger.Println(err)
		return usrs, err
	}
	defer stmt.Close()

	if err := stmt.Select(&usrs, tid); err != nil {
		tx.Rollback()
		logger.Println(err)
		return usrs, err
	}
	logger.Println(usrs)
	return usrs, nil

}

func updateTimeBlockQuery(tx *sqlx.Tx, tb *TimeBlock) error {
	q := "UPDATE time_block SET(block_id, block_start, block_end, room_id, modifier, title, note) = ($1, $2, $3, $4, $5, $6, $7) WHERE (time_block.block_id = $1)"
	stmt, err := tx.Preparex(q)
	if err != nil {
		logger.Println(err)
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(tb.ID, tb.Start, tb.End, tb.Room, tb.Modifier, tb.Title, tb.Note); err != nil {
		tx.Rollback()
		logger.Println(err)
		return err
	}
	return nil
}

func deleteTimeBlockQuery(tx *sqlx.Tx, tid int) error {
	q := "DELETE FROM time_block WHERE block_id = $1"
	stmt, err := tx.Preparex(q)
	if err != nil {
		logger.Println(err)
		return err
	}
	logger.Println("stmt")
	defer stmt.Close()
	logger.Println("stmt")
	if _, err := stmt.Exec(tid); err != nil {
		tx.Rollback()
		logger.Println(err)
		return err
	}
	logger.Println("stmt")
	return nil
}

func deleteExistingBookings(tx *sqlx.Tx, tid int) error {
	q := "DELETE from booking where block_id = $1"
	stmt, err := tx.Preparex(q)
	if err != nil {
		logger.Println(err)
		return err
	}
	logger.Println("stmt")
	defer stmt.Close()
	logger.Println("stmt")
	if _, err := stmt.Exec(tid); err != nil {
		tx.Rollback()
		logger.Println(err)
		return err
	}
	logger.Println("stmt")
	return nil
}

func createNewMessage(tx *sqlx.Tx, msg MsgWording, change string) (int, error) {
	var msgID int
	const layout = "Jan 2, 2006 at 3:04pm"

	q := "INSERT INTO notifications (msg) values ($1) RETURNING msg_id"
	newMsg := fmt.Sprintf("The event '%v' you were booked on %v-%v in room '%v' has been %v and you have been unbooked.", msg.Title, msg.Start.Format(layout), msg.End.Format(layout), msg.Room, change)
	stmt, err := tx.Preparex(q)
	if err != nil {
		logger.Println(err)
		return 0, err
	}
	defer stmt.Close()

	if err := stmt.Get(&msgID, newMsg); err != nil {
		tx.Rollback()
		logger.Println(err)
		return 0, err
	}
	return msgID, nil
}

func createMessageForUsers(tx *sqlx.Tx, u int, msgID int) error {
	q := "INSERT INTO notify (user_id, msg_id) values ($1, $2)"
	stmt, err := tx.Preparex(q)
	if err != nil {
		logger.Println(err)
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(u, msgID); err != nil {
		tx.Rollback()
		logger.Println(err)
		return err
	}
	return nil
}

/*
 * Saves the state of the block (tb)
 * to the db. Where tb is an existing block in the db
 */
func (tb *TimeBlock) update() error {
	return blockChanges(tb, true)
}

/*
 * Deletes the recieving block's (tb) entry in the db.
 */
func (tb *TimeBlock) delete() error {

	return blockChanges(tb, false)
}

func blockChanges(tb *TimeBlock, isUpdate bool) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	var msg MsgWording
	var usrs []int
	var msgID int
	change := "updated"
	if !isUpdate {
		change = "removed"
	}
	{
		msg, err = getOldTimeBlockValues(tx, tb.ID)
		if err != nil {
			return err
		}
		logger.Println(msg)
	}
	{
		usrs, err = getUsersInTimeBlockValues(tx, tb.ID)
		if err != nil {
			return err
		}
	}
	{
		err = deleteExistingBookings(tx, tb.ID)
		if err != nil {
			return err
		}
	}
	if isUpdate {
		{
			err = updateTimeBlockQuery(tx, tb)
			if err != nil {
				return err
			}
		}
	} else {
		{
			err = deleteTimeBlockQuery(tx, tb.ID)
			if err != nil {
				return err
			}
		}
	}
	{
		msgID, err = createNewMessage(tx, msg, change)
		if err != nil {
			return err
		}
	}
	for _, u := range usrs {
		{
			err = createMessageForUsers(tx, u, msgID)
			if err != nil {
				logger.Println(err)
				return err
			}
		}
	}
	return tx.Commit()
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
