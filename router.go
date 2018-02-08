// 395 Project Team Gold

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// tmp user struct just holds username and password
type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func logging(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.URL.Path)
		f.ServeHTTP(w, r)
	})
}

func errorMessage(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errResp := context.Get(r, "error")
		if errResp != nil {
			fmt.Println("CONTEXT PASSED")
			w.WriteHeader(errResp.(int))
		}
		f.ServeHTTP(w, r)
	})
}

func createRouter() (*mux.Router, error) {
	r := mux.NewRouter()
	r.Use(logging)
	r.Use(errorMessage)
	r.StrictSlash(true)
	// static file handling (put assets in views folder)
	r.PathPrefix("/views/").Handler(http.StripPrefix("/views/", http.FileServer(http.Dir("./views/"))))

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/", baseRoute)
	s.HandleFunc("/login/", loginHandler).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", calSetup).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", undoSetup).Methods("DELETE")
	s.HandleFunc("/events", getEvents).Methods("GET") // will be the blocks + bookings as a json stream
	//s.HandleFunc("/events", addBooking).Methods("POST") // Update block bookings

	v := r.PathPrefix("/app").Subrouter()
	// need redirect for '/' -> '/dashboard'
	v.HandleFunc("/dashboard/", renderDashboard).Methods("GET")
	v.HandleFunc("/schedule/", renderCalendar).Methods("GET")

	return r, nil
}

func baseRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Base route to Caraway API")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Login stuffs here")
	decoder := json.NewDecoder(r.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		logger.Fatal(err)
	}
	defer r.Body.Close()
	logger.Println("login request for user " + u.Username + u.Password)
	auth(u.Username, u.Password)
}
