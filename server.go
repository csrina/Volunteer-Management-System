package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// open config file
	config_file, err := os.Open(".config")
	if err != nil {
		panic(err)
	}
	b := make([]byte, 100)
	read, err := config_file.Read(b)
	if err != nil {
		panic(err)
	}
	config_file.Close()

	//pull out connection string
	//an example of mine is:
	//dbname=gold user=postgres host=localhost port=5454 sslmode=disable
	psqlInfo := string(b[:read])
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}
