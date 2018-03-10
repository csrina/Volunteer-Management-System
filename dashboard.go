package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/now"
)


/*
 * dashboard.js expects the following struct in json format/
 */
 type DashData struct {
 	 HoursGoal   	float64  	`json:"hoursGoal"`
	 HoursBooked 	float64  	`json:"hoursBooked"`
	 HoursDone   	float64  	`json:"hoursDone"`
	 History     	[]float64 	`json:"history"` // historical hours completed/week for interval
 }

 /* Replacement for dashboardData
  * which delegates most responsibility to functions
  */
func getDashData(w http.ResponseWriter, r *http.Request) {
	UID 	:= getUID(r)
	goal 	:= getHoursGoal(UID)
	booked 	:= getHoursBooked(UID)
	done 	:= getHoursDone(UID)
	history := getHoursHistory(UID, time.Now())

	dd := &DashData{
		HoursGoal: goal,
		HoursBooked: booked,
		HoursDone: done,
		History: history,
	}
	logger.Println("DASHDATA: ", dd)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	encoder.Encode(dd)
}

/* Get the weekly goal hours for a user */
func getHoursGoal(UID int) float64 {
	numKids := 0
	q := `SELECT children FROM family
			WHERE (parent_one = $1 OR parent_two = $1)`
	err := db.QueryRow(q, UID).Scan(&numKids)
	if err != nil {
		logger.Println(err)
		logger.Println(">> ERROR: gethoursgoal")
		return -1
	}
	if numKids > 1 {
		return 7.50 // or however much multikid families must pay to educate their kids
	}
	return 5.00
}

/* Get the booked hours for a user */
func getHoursBooked(UID int) float64 {
	logger.Println("getting hours booked")
	start := now.BeginningOfWeek()
	end := now.EndOfWeek()
	t := time.Now().Weekday()
	if t == 6 || t == 0{ //its a weekend
		start = start.AddDate(0,0,8)
		end = end.AddDate(0,0,8)
		logger.Println(">> Start: ", start, " end: ", end)
	}
	bookBlocks, err := getUserBookings(start, end, UID)
	if err != nil {
		logger.Println(">> ERROR: getHoursBooked")
		return -1
	}
	return getHoursBookingSlice(bookBlocks)
}

/* Gets hours completed this week */
func getHoursDone(UID int) float64 {
	logger.Println("getting hours done")
	start := now.BeginningOfWeek()
	end := time.Now()
	t := end.Weekday()
	if t == 6 || t == 0{ //its a weekend
		start = end 
	}
	bookBlocks, err := getUserBookings(start, end, UID)
	if err != nil {
		logger.Println(">> ERROR: getHoursDone")
		return -1
	}
	return getHoursBookingSlice(bookBlocks)
}

func getHoursBookingSlice(bookBlocks []Booking) float64 {
	logger.Println("getHoursFromBookingSlice\nBLOCKSGIVEN: ", bookBlocks)
	var duration float64
	duration = 0.00
	for _, bb := range bookBlocks {
		duration += (bb.End.Sub(bb.Start).Hours() * float64(bb.Modifier))
	}
	return duration
}

/* Returns historical hours/week for past 3 months */
/* Does not work yet. Do not expect good results yet */
func getHoursHistory(UID int, curr time.Time) []float64 {
	start := now.New(curr).BeginningOfWeek().AddDate(0,-3, 0) // 3 months prior to beginning of week
	bookBlocks, err := getUserBookings(start, time.Now(), UID)
	if err != nil {
		logger.Println(err)
		return nil
	}
	var history []float64
	var duration float64
	duration = 0.00
	for i, bb := range bookBlocks {
		duration = (bb.End.Sub(bb.Start).Hours() * float64(bb.Modifier))
		if i % 5 == 0 { // Need to give this correct logic for deducing weeks
			history = append(history, duration)
			duration = 0.00 // reset to 0, week end
		}
	}
	return history
}

/*
join family and users

join booking and time block

pull out username, user_id, block_start, block_end
*/
type WeeklyBooking struct {
	UserName   string    `db:"username"`
	UserId     int       `db:"user_id"`
	BlockStart time.Time `db:"block_start"`
	BlockEnd   time.Time `db:"block_end"`
	Children   int       `db:"children"`
}

type FriendlyFormat struct {
	Eventlist   []string `json:"eventlist"`
	HoursBooked float64  `json:"hoursBooked"`
	HoursDone   float64  `json:"hoursDone"`
	Children    int      `json:"children"`
}

func startOfWeek(current time.Time) time.Time {
	layoutDay := "Mon"
	check := current.Format(layoutDay)
	for check != "Mon" {
		current = current.AddDate(0, 0, -1)
		check = current.Format(layoutDay)
	}
	return current
}

func dashboardData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}
	// make session stuff
	user := getUID(r)
	role, err := getRoleNum(r)
	logger.Printf("User ID found %d\n\n\n", user)
	start := startOfWeek(time.Now())
	var q string
	if role == 1 {
		q = `SELECT username, s.user_id, block_start, block_end, children 
FROM (
SELECT distinct username, user_id, children
FROM family INNER JOIN users
ON user_id = parent_one or user_id = parent_two) r
INNER JOIN
(SELECT user_id, block_start, block_end
FROM booking b INNER JOIN time_block t 
ON b.block_id = t.block_id) s
ON r.user_id = s.user_id
WHERE r.user_id = $1 AND block_start > $2 AND block_start < $3`
	} else {
		fmt.Fprintln(w, "non facilitator doesnt have dashboard right now")
		return
	}

	var bookings []WeeklyBooking

	err = db.Select(&bookings, q, user, start, start.AddDate(0, 0, 6))
	if err != nil {
		fmt.Printf("%v", err)
	}
	var friendly FriendlyFormat
	i := 0
	hoursDone := 0.0
	hoursBooked := 0.0
	layout := "Mon Jan 2 15:04"
	if len(bookings) == 0 {
		fmt.Fprintln(w, "no results")
		return
	}
	for each := range bookings {
		if bookings[each].BlockEnd.Before(time.Now()) {
			hoursDone += bookings[each].BlockEnd.Sub(bookings[each].BlockStart).Hours()
		}
		hoursBooked += bookings[each].BlockEnd.Sub(bookings[each].BlockStart).Hours()
		if bookings[each].BlockEnd.After(time.Now()) {
			friendly.Eventlist = append(friendly.Eventlist, bookings[each].BlockStart.Format(layout)+" to ")
			friendly.Eventlist[i] = friendly.Eventlist[i] + bookings[each].BlockEnd.Format(layout)
			i++
		}
	}
	friendly.HoursBooked = hoursBooked
	friendly.HoursDone = hoursDone
	friendly.Children = bookings[0].Children

	encoder := json.NewEncoder(w)
	logger.Println("FRIENDLY:  ", friendly)
	err = encoder.Encode(friendly)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
