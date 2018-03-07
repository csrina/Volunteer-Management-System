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
	fmt.Printf("User ID found %d\n\n\n", user)
	start := startOfWeek(time.Now())
	var q string
	if role == 1 {
		q = `SELECT username, s.user_id, block_start, block_end, children 
FROM (
SELECT username, user_id, children
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
	err = encoder.Encode(friendly)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
