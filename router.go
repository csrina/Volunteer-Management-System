// 395 Project Team Gold

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.URL.Path)
		f(w, r)
	}
}

func createRouter() (*mux.Router, error) {
	r := mux.NewRouter()
	r.StrictSlash(true)
	// static file handling (put assets in views folder)
	r.PathPrefix("/views/").Handler(http.StripPrefix("/views/", http.FileServer(http.Dir("./views/"))))
	r.PathPrefix("/login/").Handler(http.StripPrefix("/login/", http.FileServer(http.Dir("./views/login/"))))
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/", logging(baseRoute))
	s.HandleFunc("/login/facilitator/", logging(loginHandler)).Methods("POST")
	s.HandleFunc("/login/teacher/", logging(loginHandler)).Methods("POST")
	s.HandleFunc("/login/admin/", logging(loginHandler)).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", calSetup).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", undoSetup).Methods("DELETE")
	s.HandleFunc("/events", getEvents).Methods("GET") // will be the blocks + bookings as a json stream
	//s.HandleFunc("/events", addBooking).Methods("POST") // Update block bookings

	v := r.PathPrefix("/app").Subrouter()
	// need redirect for '/' -> '/dashboard'
	v.HandleFunc("/dashboard/", logging(renderDashboard)).Methods("GET")
	v.HandleFunc("/schedule/", logging(renderCalendar)).Methods("GET")

	return r, nil
}

func baseRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Base route to Caraway API")
}

func loadLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/login/facilitatorLogin.html")
}
