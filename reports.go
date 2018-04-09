//395 project team gold
//API functions to create admin reports

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/now"
)

/* START OF SCHOOL YEAR FOR HISTORY TRACKING */
const (
	PERIOD_LENGTH        = 3 // months
	ONE_CHILD_HOURS_GOAL = 2.50
	DEFAULT_HOURS_GOAL   = 5.00
)

/* Find start/end contstraints within the month */
func setWeekConstraint(time time.Time) (start, end string) {
	check := now.New(time)
	start = check.BeginningOfWeek().Format("Mon Jan 2")
	end = check.EndOfWeek().Format("Mon Jan 2")

	return start, end
}

func getHourGoal(children int) float64 {
	if children == 1 {
		return ONE_CHILD_HOURS_GOAL
	}
	return DEFAULT_HOURS_GOAL
}

type familyReport struct {
}

func monthlyReport(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("date")
	req := strings.Split(raw, "-")

	addMonth, _ := strconv.Atoi(req[1])

	q := `SELECT family_id, family_name, children
			FROM family ORDER BY UPPER(family_name)`

	families := []familyShort{}
	err := db.Select(&families, q)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	layout := "2006"
	month := []monthReport{}

	start, err := time.Parse(layout, req[0])
	if err != nil {
		fmt.Printf("problem parsing date in query")
	}

	copy := start.AddDate(0, addMonth-1, 0)
	for i, fam := range families {
		goal := getHourGoal(fam.Children)
		start = copy
		end := now.New(start)
		month = append(month, monthReport{})
		for start.Before(end.EndOfMonth()) {
			begin, finish := setWeekConstraint(start)
			hours := familyHoursBooked(fam.FamilyID, start,
				start.AddDate(0, 0, 7))
			month[i].Weeks = append(month[i].Weeks, weekReport{
				Start: begin,
				End:   finish,
				Total: hours - goal,
			})
			start = start.AddDate(0, 0, 8)
			start = now.New(start).BeginningOfWeek()
		}
		month[i].FamilyID = fam.FamilyID
		month[i].FamilyName = fam.FamilyName
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(month)
}

func defaultReport(w http.ResponseWriter, r *http.Request) {
	q := `SELECT family_id, family_name
			FROM family ORDER BY UPPER(family_name)`

	families := []familyShort{}

	err := db.Select(&families, q)
	if err != nil {
		logger.Println(err)

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
