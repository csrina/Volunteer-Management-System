//395 project team gold
//API functions to create admin reports

package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jinzhu/now"
)

type familyShort struct {
	FamilyID   int     `json:"familyId" db:"family_id"`
	FamilyName string  `json:"familyName" db:"family_name"`
	WeekHours  float64 `json:"weekHours"`
}

func defaultReport(w http.ResponseWriter, r *http.Request) {
	q := `SELECT family_id, family_name
			FROM family`

	families := []familyShort{}

	err := db.Select(&families, q)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, fam := range families {
		hours := familyHoursBooked(fam.FamilyID,
			now.BeginningOfWeek(), now.EndOfWeek())
		families[i].WeekHours = hours
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(families)
}

func monthlyReport(w http.ResponseWriter, r *http.Request) {
	q := `SELECT family_id, family_name
			FROM family`

	families := []familyShort{}

	err := db.Select(&families, q)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, fam := range families {
		hours := familyHoursBooked(fam.FamilyID,
			now.BeginningOfMonth(), now.EndOfMonth())
		families[i].WeekHours = hours
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(families)
}

func familyHoursBooked(FID int, start time.Time, end time.Time) float64 {
	bookBlocks, err := getFamilyBookings(start, end, FID)
	if err != nil {
		return -1
	}
	return getHoursBookingSlice(bookBlocks)
}

func getHoursBookingSlice(bks []Booking) float64 {
	duration := 0.00
	for _, b := range bks {
		duration += (b.End.Sub(b.Start).Hours() * float64(b.Modifier))
	}
	return duration
}

func getFamilyBookings(start time.Time, end time.Time, FID int) ([]Booking, error) {
	/* format dates for psql */
	logger.Println("Preformat --> Start ", start, "\tEnd ", end)

	/* Get all bookings in range start-now  (start > block_start & end > blocK_end) */
	q := `SELECT booking_id, block_id, family_id, user_id, block_start, block_end, room_id, modifier
			FROM booking NATURAL JOIN time_block WHERE (
					time_block.block_id = booking.block_id
					AND booking.family_id = $1
					AND time_block.block_start >= $2 AND time_block.block_start < $3
					AND time_block.block_end > $2 AND  time_block.block_end <= $3
			) ORDER BY block_start`

	var bookBlocks []Booking
	err := db.Select(&bookBlocks, q, FID, start, end)
	logger.Println("Selected blocks: ", bookBlocks)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return bookBlocks, nil
}
