// 395 Project Team Gold

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// tmp user struct just holds username and password
type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.URL.Path)
		f(w, r)
	}
}

func createRouter() (*mux.Router, error) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/", logging(baseRoute))
	s.HandleFunc("/login/", logging(loginHandler)).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", calSetup).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", undoSetup).Methods("DELETE")

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
