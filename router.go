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
	r.PathPrefix("/tmp/").Handler(http.StripPrefix("/tmp/", http.FileServer(http.Dir("./public/"))))
	// r.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("public"))))

	r.HandleFunc("/login", loadMainLogin)
	l := r.PathPrefix("/login").Subrouter()
	l.HandleFunc("/facilitator", loadLogin)
	l.HandleFunc("/teacher", loadLogin)
	l.HandleFunc("/admin", loadLogin)

	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/admin/calendar/setup/", calSetup).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", undoSetup).Methods("DELETE")
	/* Events JSON routes for scheduler system */
	s.HandleFunc("/events", getEvents).Methods("GET")
	s.HandleFunc("/events/{target}", eventPostHandler).Methods("POST")
	v := r.PathPrefix("/app").Subrouter()
	// need redirect for '/' -> '/dashboard'
	v.HandleFunc("/dashboard/", renderDashboard).Methods("GET")
	/* Calendar requests */
	v.HandleFunc("/schedule/", renderCalendar).Methods("GET")

	l = s.PathPrefix("/login").Subrouter()
	l.HandleFunc("/facilitator/", loginHandler).Methods("POST")
	l.HandleFunc("/teacher/", loginHandler).Methods("POST")
	l.HandleFunc("/admin/", loginHandler).Methods("POST")

	return r, nil
}

func baseRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Base route to Caraway API")
}
