// 395 Project Team Gold

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func createRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/", baseRoute)

	return r, nil
}

func baseRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Base route to Caraway API")
}
