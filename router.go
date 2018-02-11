// 395 Project Team Gold

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

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
	r.PathPrefix("/login/").Handler(http.StripPrefix("/login/", http.FileServer(http.Dir("./views/login/"))))
	s := r.PathPrefix("/api/v1").Subrouter()
<<<<<<< HEAD
	s.HandleFunc("/", logging(baseRoute))
	s.HandleFunc("/login/facilitator/", logging(loginHandler)).Methods("POST")
	s.HandleFunc("/login/teacher/", logging(loginHandler)).Methods("POST")
	s.HandleFunc("/login/admin/", logging(loginHandler)).Methods("POST")
=======
	s.HandleFunc("/", baseRoute)
	s.HandleFunc("/login/", loginHandler).Methods("POST")
>>>>>>> fc6eb61ff15395ea2b5b73be1f1d970b281f9e58
	s.HandleFunc("/admin/calendar/setup/", calSetup).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", undoSetup).Methods("DELETE")
	/* Events JSON routes for scheduler system */
	s.HandleFunc("/events", getEvents).Methods("GET")
	s.HandleFunc("/events", updateEvent).Methods("POST") // Changes made to schedule,  update block

	v := r.PathPrefix("/app").Subrouter()
	// need redirect for '/' -> '/dashboard'
	v.HandleFunc("/dashboard/", renderDashboard).Methods("GET")
	/* Calendar requests */
	v.HandleFunc("/schedule/", renderCalendar).Methods("GET")

	return r, nil
}

func baseRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Base route to Caraway API")
}

func loadLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/login/facilitatorLogin.html")
}
