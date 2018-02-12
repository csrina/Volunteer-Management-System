package main

import (
	"fmt"
	"time"
)

/*
join family and users

join booking and time block

pull out username, user_id, block_start, block_end
*/
type TempStruct struct {
	UserName string `db:"username"`
	UserId   int    `db:"user_id"`
	//	BlockId    int       `db:"block_id"`
	BlockStart time.Time `db:"block_start"`
	BlockEnd   time.Time `db:"block_end"`
}

func dashBoard() {
	q := `SELECT username, s.user_id, block_start, block_end 
FROM (
SELECT username, user_id
FROM family INNER JOIN users
ON user_id = parent_one or user_id = parent_two) r
INNER JOIN
(SELECT user_id, block_start, block_end
FROM booking b INNER JOIN time_block t 
ON b.block_id = t.block_id) s
ON r.user_id = s.user_id`

	var bookings []TempStruct

	err := db.Select(&bookings, q)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%#v", bookings[0])
}
