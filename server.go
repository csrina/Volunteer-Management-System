package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var logger *log.Logger

type args struct {
	ServerPort string
	LogName    string
}

// Args used to hold any command line arguments
var Args args

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.StringVar(&Args.LogName, "l", "os.Stderr", "set logfile name")
	flag.StringVar(&Args.ServerPort, "p", ":8080", "set webserver port")

	flag.Parse()
}

func startDb() error {
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
	return err
}

func main() {

	if Args.ServerPort != "os.Stderr" {
		f, err := os.OpenFile(Args.LogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		logger = log.New(f, "status: ", log.LstdFlags)
	} else {
		logger = log.New(os.Stderr, "status: ", log.LstdFlags)
	}
	// open config file
	logger.Println("Reading config file...")
	err := startDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logger.Println("Databse opened and pinged.....")

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
