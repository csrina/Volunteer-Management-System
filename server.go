package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	b := make([]byte, 100)
	config_file, err := os.Open(".config")

	if err != nil {
		panic(err)
	}
	read, err := config_file.Read(b)
	if err != nil {
		panic(err)
	}
	config_file.Close()
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

	_, err = db.Query("create table class(id serial, name text)")
	if err != nil {
		fmt.Println("table already created, go make psql load file")
		//panic(err)
	}
}
