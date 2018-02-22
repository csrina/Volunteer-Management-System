package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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
}

type FriendlyFormat struct {
	Eventlist   []string `json:"eventlist"`
	HoursBooked float64  `json:"hoursBooked"`
	HoursDone   float64  `json:"hoursDone"`
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
	user := 4

	start := startOfWeek(time.Now())

	q := `SELECT username, s.user_id, block_start, block_end 
FROM (
SELECT username, user_id
FROM family INNER JOIN users
ON user_id = parent_one or user_id = parent_two) r
INNER JOIN
(SELECT user_id, block_start, block_end
FROM booking b INNER JOIN time_block t 
ON b.block_id = t.block_id) s
ON r.user_id = s.user_id
WHERE r.user_id = $1 AND block_start > $2 AND block_start < $3`

	//need to make this query only look to the end of the week!

	var bookings []WeeklyBooking

	err := db.Select(&bookings, q, user, start, start.AddDate(0, 0, 6))
	if err != nil {
		fmt.Printf("%v", err)
	}
	var friendly FriendlyFormat
	i := 0
	hoursDone := 0.0
	hoursBooked := 0.0
	layout := "Mon Jan 2 15:04"
	for each := range bookings {
		friendly.Eventlist = append(friendly.Eventlist, bookings[each].BlockStart.Format(layout)+" to ")
		friendly.Eventlist[i] = friendly.Eventlist[i] + bookings[each].BlockEnd.Format(layout)
		if bookings[each].BlockEnd.Before(time.Now()) {
			hoursDone += bookings[each].BlockEnd.Sub(bookings[each].BlockStart).Hours()
		}
		hoursBooked += bookings[each].BlockEnd.Sub(bookings[each].BlockStart).Hours()
		i++
	}
	friendly.HoursBooked = hoursBooked
	friendly.HoursDone = hoursDone

	encoder := json.NewEncoder(w)
	err = encoder.Encode(friendly)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
