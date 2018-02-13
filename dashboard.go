package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/*
join family and users

join booking and time block

pull out username, user_id, block_start, block_end
*/
type WeeklyBooking struct {
	UserName string `db:"username"`
	UserId   int    `db:"user_id"`
	//	BlockId    int       `db:"block_id"`
	BlockStart time.Time `db:"block_start"`
	BlockEnd   time.Time `db:"block_end"`
}

func dashboardLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}
	vars := mux.Vars(r)
	user := vars["userId"]
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
WHERE r.user_id = $1 AND block_start > TIMESTAMP 'now'`

	//need to make this query only look to the end of the week!

	var bookings []WeeklyBooking

	err := db.Select(&bookings, q, user)
	if err != nil {
		fmt.Printf("%v", err)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(bookings)
}
