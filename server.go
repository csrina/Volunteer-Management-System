package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var logger *log.Logger

func main() {

	logger = log.New(os.Stdout, "status: ", log.LstdFlags)
	// open config file
	logger.Println("Reading config file...")

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
	//dbname=caraway user=postgres host=localhost port=5454 sslmode=disable
	psqlInfo := string(b[:read])
	logger.Println("Opening database...")
	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logger.Println("Databse opened and pinged.....")
	defer db.Close()

	r, err := createRouter()
	if err != nil {
		log.Fatal("Could not create router")
	}
	logger.Println("Routes created")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Could not start server")
	}
	logger.Println("Server running......")
}
