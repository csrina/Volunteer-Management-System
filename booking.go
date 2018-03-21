package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
)

// bookingBlock - struct matching a booking record
//noinspection GoUnusedType
type bookingBlock struct {
	BookingID int       `db:"booking_id" json:"bookingId"`
	BlockID   int       `db:"block_id" json:"blockId"`
	FamilyID  int       `db:"family_id" json:"familyId"`
	UserID    int       `db:"user_id" json:"userID"`
	Start     time.Time `db:"booking_start" json:"start"`
	End       time.Time `db:"booking_end" json:"end"`
}

/*
 * Initialize a Booking struct given a BookingID
 */
func (b *bookingBlock) init(BID int) error {
	q := `SELECT booking_id, block_id, family_id, user_id, COALESCE(LOCALTIMESTAMP, booking_start), COALESCE(LOCALTIMESTAMP, booking_end) FROM booking WHERE booking_id = $1`
	err := db.QueryRow(q, BID).Scan(&b.BookingID, &b.BlockID, &b.FamilyID, &b.UserID, &b.Start, &b.End)
	if err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

// inserts a booking into the database
func (b *bookingBlock) insertBooking() error {
	q := `INSERT INTO booking (block_id, family_id, user_id, 
			booking_start, booking_end) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING booking_id`

	err := db.QueryRow(q, b.BlockID, b.FamilyID, b.UserID, b.Start, b.End).Scan(&b.BookingID)
	if err != nil {
		return err
	}
	return nil
}

// updates an existing booking in the database
func (b *bookingBlock) updateBooking() error {
	q := `UPDATE booking SET booking_start = $2, booking_end = $3
			WHERE booking_id = $3`

	results, err := db.Exec(q, b.BlockID, b.Start, b.End)
	if err != nil {
		return err
	}
	count, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		err = errors.New("booking not updated")
		return err
	}
	return nil
}

// deletes an existing booking in the database
func (b *bookingBlock) deleteBooking() error {
	q := `DELETE FROM booking WHERE booking_id = $1`

	results, err := db.Exec(q, b.BookingID)
	if err != nil {
		logger.Println(err)
		return err
	}
	count, err := results.RowsAffected()
	if err != nil {
		logger.Println(err)
		return err
	}
	if count != 1 {
		err = errors.New("booking not deleted")
		logger.Println(err)
		return err
	}
	return nil
}

// -------------------------------------------------
// -------------------------------------------------
func bookingHandler(w http.ResponseWriter, r *http.Request, role int) {
	booking, err := bookingFromJSON(r)
	if err != nil {
		logger.Println(err)
		http.Error(w, "could not resolve booking data, contact your administrator if the problem persists", http.StatusBadRequest)
		return
	}
	response := new(Response) // Response data
	usr := UserShort{UserID: booking.UserID}
	response.UserName, _ = usr.getFullName()
	response.UID = booking.UserID;

	// book/unbook based on BookingID (if 0 value, then booking doesn't exist yet and must be created)
	if booking.BookingID > 0 {
		err = booking.unbook(role)
		response.Msg = "Successfully cancelled booking"
	} else {
		err = booking.book(role)
		response.Msg = "Successfully created booking"
	}
	// We may have an error condition after book/unbook
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.setBID(booking.BookingID).send(w)
}

/* Joined relation of time_block and booking_block */
type Booking struct {
	BookingID int       `db:"booking_id" json:"bookingId"`
	BlockID   int       `db:"block_id" json:"blockId"`
	FamilyID  int       `db:"family_id" json:"familyId"`
	UserID    int       `db:"user_id" json:"userID"`
	Start     time.Time `db:"block_start" json:"endBlock"`
	End       time.Time `db:"block_end" json:"endBlock"`
	Room      string    `json:"color"`
	RoomID    int       `db:"room_id" json:"room_id"`
	Modifier  int       `db:"modifier" json:"modifier"`
	Note      string    `db:"note" json:"note"`
}


/*
 * Reads json request, and creates a partially filled booking struct
 */
func bookingFromJSON(r *http.Request) (*Booking, error) {
	booking := new(Booking)
	// Now extract data from json request body
	ev, err := mapEventJSON(r)
	if err != nil {
		return nil, err
	}
	eid, ok := ev["id"].(float64)
	if !ok {
		eid = -1
	}
	uid := getUID(r)
	if uid < 1 {
		return nil, errors.New("uid unresolved -- please re-login")
	}
	if uidTest, ok := ev["userId"].(float64); ok {
		uid = int(uidTest)
	} else if uName, ok := ev["userId"].(string); ok {
		uid, err = getUIDFromName(uName)
		if err != nil {
			logger.Println("Failed to resolve: ", ev["userId"]);
			return nil, err
		}
	}
	booking.UserID = uid
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

func (b *Booking) getTimeBlock(role int) (*TimeBlock, error) {
	if role != ADMIN {
		return nil, errors.New("Insufficient priviledge")
	}
	tb := new(TimeBlock)
	tb.ID = b.BlockID
	tb.Start = b.Start
	tb.End = b.End
	tb.Note = b.Note
	tb.Modifier = b.Modifier
	if b.RoomID < 1 {
		q := `SELECT room_id FROM room WHERE room_name = $1`
		err := db.QueryRow(q, b.Room).Scan(&tb.Room)
		if err != nil {
			return nil, err
		}
	}
	return tb, nil
}

/*
 * Returns the booking id for the event/user if it exists,
 * Side-EFfect: If the bookingID is unset --> checks db for booking
 *      if bookingBlock exists in DB: returns the id AND sets the struct bookingID
 */
func (b *Booking) getBookingID() (int, error) {
	if b.BookingID > 0 {
		return b.BookingID, nil
	}
	q := `SELECT booking_id FROM booking WHERE (block_id = $1 AND user_id = $2)`
	err := db.QueryRow(q, b.BlockID, b.UserID).Scan(&b.BookingID)
	if err != nil {
		return -1, err
	}
	return b.BookingID, nil
}

/*
 * The returned map has the keys: "start" and "end"
 *
 * Side-Effect: Like getBookingID, getTimesMap will actively update b.Start and b.End,
 * if they are zeroed (default new() state).
 */
func (b *Booking) getTimesMap() (map[string]time.Time, error) {
	// Test if dates are set (note events which are infinite, will always go thru DB)
	if b.Start.IsZero() && b.End.IsZero() {
		q := `SELECT block_start, block_end FROM time_block 
				WHERE block_id = $1`
		err := db.QueryRow(q, b.BlockID).Scan(&b.Start, &b.End)
		if err != nil {
			return nil, err
		}
	}
	b.Start = time.Date(b.Start.Year(), b.Start.Month(), b.Start.Day(),
		b.Start.Hour(), b.Start.Minute(), 0, 0, time.Local)
	b.End = time.Date(b.End.Year(), b.End.Month(), b.End.Day(),
		b.End.Hour(), b.End.Minute(), 0, 0, time.Local)
	times := make(map[string]time.Time)
	times["start"] = b.Start
	times["end"] = b.End
	return times, nil
}

/*
 * Determines if a booking is legal: that is, the maximum capacity is not reached,
 * and/or the booking is not in past.
 */
func (b *Booking) isLegal() (bool, string) {
	if b.Start.Before(time.Now()) && b.End.Before(time.Now()) {
		return false, "event has passed" // Cannot book if the block is in the past!
	}

	var bids []int
	q := `SELECT booking_id FROM booking WHERE booking.block_id = $1`
	db.Select(&bids, q, b.BlockID)
	if len(bids) >= 3 {
		return false, "event is full"
	}
	// Check if overlaps with existing booking of user
	q = `SELECT booking_id FROM booking NATURAL JOIN time_block
					WHERE booking.user_id = $1 AND booking.block_id = time_block.block_id 
						AND (
								(time_block.block_start <= $2 AND $2 <= time_block.block_end)
							OR
								(time_block.block_start <= $3 AND $3 <= time_block.block_end)
						)`

	db.Select(&bids, q, b.UserID, b.Start, b.End)
	return len(bids) == 0, "conflicts with an existing booking"
}

/*
 * Creates a bookingBlock in the DB for a Booking struct
 * which has not yet been stored. The struct's bookingID
 * is set upon successful booking.
 */
func (b *Booking) book(role int) error {
	if role == FACILITATOR {
		if ok, reason := b.isLegal(); !ok {
			return errors.New(reason)
		}
	}
	// Dont update booking_start and booking_end in DB --> these are the ACTUAL start/end times
	q := `INSERT INTO booking (block_id, user_id, family_id) 
			VALUES ($1, $2, $3)
			RETURNING booking_id`

	err := db.QueryRow(q, b.BlockID, b.UserID, b.FamilyID).Scan(&b.BookingID)
	if err != nil {
		logger.Println("Error creating booking: ", err, "\nbooking data: ", b)
		return err
	}
	logger.Println("New Booking created!\nBooking id: ", b.BookingID)
	return nil
}

func (b *Booking) unbook(role int) error {
	if b.Start.Before(time.Now()) && b.End.Before(time.Now()) {
		logger.Println("Now: ", time.Now(), "Start: ", b.Start, "  local: ", b.Start.Local())
		return errors.New("cannot unbook from completed event")
	}

	if b.Start.Before(time.Now()) || b.End.Before(time.Now()) && role == FACILITATOR {
		return errors.New("only administration may remove bookings in progress")
	}

	q := `DELETE FROM booking WHERE booking_id = $1`
	_, err := db.Exec(q, b.BookingID)
	if err != nil {
		logger.Println("Error deleting booking (id = ", b.BookingID, "): ", err)
		return err
	}
	b.BookingID = 0
	return nil
}

/*
 * Booking is the usual struct for calendar based interactions with the schedule
 * (in combination with Event)
 *
 * This method updates the underlying timeblock of a Booking.
 *       E.g. when administrator resizes an event on the calendar
 */
func (b *Booking) updateTimeBlock(role int) error {
	if role != ADMIN {
		return errors.New("insufficient permission for changing time block duration")
	}
	/* Update DB reference */
	q := `UPDATE time_block 
			SET (block_start, block_end) = ($1, $2)
			WHERE (block_id = $3)`

	_, err := db.Exec(q, b.Start, b.End, b.BlockID)
	return err
}

/*
 *  Given a blockID, get the number of bookings created
 */
func getBookingCount(blockID int) int {
	cnt := 0
	q := `SELECT count(*) FROM booking WHERE block_id = $1`
	db.QueryRow(q, blockID).Scan(&cnt)
	return cnt
}

/* Ayyy */
func getUserBookings(start time.Time, end time.Time, UID int) ([]Booking, error) {
	/* Get all bookings in range start-now  (start > block_start & end > blocK_end) */
	q := `SELECT booking_id, block_id, family_id, user_id, block_start, block_end, room_id, modifier
			FROM booking NATURAL JOIN time_block WHERE (
					time_block.block_id = booking.block_id
					AND booking.user_id = $1
					AND time_block.block_start >= $2 AND time_block.block_start < $3
					AND time_block.block_end > $2 AND  time_block.block_end <= $3
			) ORDER BY block_start`

	var bookBlocks []Booking
	err := db.Select(&bookBlocks, q, UID, start, end)
	logger.Println("Selected blocks: ", bookBlocks)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return bookBlocks, nil
}

/* Like get bookings for a family*/
func (f *Family) getFamilyBookings(start time.Time, end time.Time) ([]Booking, error) {
	/* Get all bookings in range start-now  (start > block_start & end > blocK_end) */
	q := `SELECT booking_id, block_id, family_id, user_id, block_start, block_end, room_id, modifier
			FROM booking NATURAL JOIN time_block WHERE (
					time_block.block_id = booking.block_id
					AND booking.family_id = $1
					AND time_block.block_start >= $2 AND time_block.block_start < $3
					AND time_block.block_end > $2 AND  time_block.block_end <= $3
			) ORDER BY block_start`

	var bookBlocks []Booking
	err := db.Select(&bookBlocks, q, f.ID, start, end)
	logger.Println("Selected blocks: ", bookBlocks)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Ensure locale is set
	for i, b := range bookBlocks {
		bookBlocks[i].Start = time.Date(b.Start.Year(), b.Start.Month(), b.Start.Day(),
			b.Start.Hour(), b.Start.Minute(), 0, 0, time.Local)
		bookBlocks[i].End = time.Date(b.End.Year(), b.End.Month(), b.End.Day(),
			b.End.Hour(), b.End.Minute(), 0, 0, time.Local)

	}
	return bookBlocks, nil
}
