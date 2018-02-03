// 395 Project Team Gold

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func createRouter() (*mux.Router, error) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/", baseRoute).Methods("GET")
	s.HandleFunc("/login/", login)
	s.HandleFunc("/admin/calendar/setup/", calSetup).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", undoSetup).Methods("DELETE")

	return r, nil
}

func baseRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Base route to Caraway API")
	logger.Println("Call to baseRoute")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login stuffs here")
	logger.Println("Call to login")
}
