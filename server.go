package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var logger *log.Logger

func main() {

	logger = log.New(os.Stdout, "status: ", log.LstdFlags)
	// open config file
	logger.Println("reading config file...")

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
	logger.Println("opening database...")
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()
}
