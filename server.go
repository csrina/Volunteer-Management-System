package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// open config file
	config_file, err := os.Open(".config")
	if err != nil {
		panic(err)
	}
	read, err := config_file.Read(b)
	if err != nil {
		panic(err)
	}
	config_file.Close()

	//pull out connection string
	//an example of mine is:
	//dbname=gold user=postgres host=localhost port=5454 sslmode=disable
	b := make([]byte, 100)
	psqlInfo := string(b[:read])
	fmt.Println(read)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	//this is only in to make sure it is correctly adding a
	//table to the specified database
	_, err = db.Query("create table class(id serial, name text)")
	if err != nil {
		fmt.Println("table already created, go make psql load file")
		//panic(err)
	}
}
