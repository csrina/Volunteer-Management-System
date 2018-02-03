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
	s.HandleFunc("/", baseRoute)
	s.HandleFunc("/login/", login)

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
