package main

import (
	_ "github.com/lib/pq"
)

func auth(username string, password string) {
	q := `SELECT username, password
			FROM users
			where username = $1`
	logger.Println("Checking if user exits")
	logger.Println(q + " " + username)
	users := []User{}
	if err := db.Select(&users, q, username); err != nil {
		logger.Println(err)
	}
	count := len(users)
	if count == 0 {
		logger.Printf("User %v not found.\n", username)
	} else if count == 1 {
		logger.Printf("User %v found\n", username)
		if password == users[0].Password {
			logger.Printf("Password is correct\n")
		} else {
			logger.Printf(users[0].Password)
			logger.Printf("Password is wrong\n")
		}
	} else {
		logger.Printf("Unknown error occurd in auth()\n")
	}
	logger.Printf("number of users found %v\n", len(users))
}
