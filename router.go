// 395 Project Team Gold

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func logging(f http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		logger.Printf("'%v' %v", r.Method, r.URL.Path)
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

func checkSession(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "loginSession")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve our struct and type-assert it
		val := session.Values["username"]
		if val != nil {
			logger.Println(val)
			f.ServeHTTP(w, r)
		} else {
			if strings.Contains(r.URL.Path, "login") {
				f.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/login", http.StatusFound)
			}
		}
		return
	})
}

func createRouter() (*mux.Router, error) {
	r := mux.NewRouter()
	r.Use(logging)
	r.Use(errorMessage)

	r.StrictSlash(true)
	// static file handling (put assets in views folder)
	r.PathPrefix("/views/").Handler(http.StripPrefix("/views/", http.FileServer(http.Dir("./views/"))))

	// api routes+calls set up
	apiRoutes(r)

	r.HandleFunc("/", baseRoute)
	// load login pages html tmplts
	r.HandleFunc("/login", loadMainLogin)
	r.HandleFunc("/logout", handleLogout)
	l := r.PathPrefix("/login").Subrouter()
	l.HandleFunc("/facilitator", loadLogin)
	l.HandleFunc("/teacher", loadLogin)
	l.HandleFunc("/admin", loadLogin)

	r.Use(checkSession)

	//load dashboard and calendar pages
	r.HandleFunc("/dashboard", loadDashboard)
	r.HandleFunc("/calendar", loadCalendar)

	return r, nil
}

func apiRoutes(r *mux.Router) {
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/admin/calendar/setup/", calSetup).Methods("POST")
	s.HandleFunc("/admin/calendar/setup/", undoSetup).Methods("DELETE")
	s.HandleFunc("/dashboard", getDashData).Methods("GET")

	/* Events JSON routes for scheduler system */
	s.HandleFunc("/events/{target}", getEvents).Methods("GET")
	s.HandleFunc("/events/{target}", eventPostHandler).Methods("POST")
	l := s.PathPrefix("/login").Subrouter()
	l.HandleFunc("/facilitator/", loginHandler).Methods("POST")
	l.HandleFunc("/teacher/", loginHandler).Methods("POST")
	l.HandleFunc("/admin/", loginHandler).Methods("POST")
}

//noinspection ALL
func baseRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Base route to Caraway API, redirecting to main page")
	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}

func loadDashboard(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("dashboard", r) // load page
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("dashboard.tmpl")
	// dependency flags for dashboard
	pg.Calendar = true
	pg.Chart = true
	s.ExecuteTemplate(w, "dashboard", pg) // include page struct
}

func loadCalendar(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("calendar", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("calendar.tmpl")
	// calendar dependency flag
	pg.Calendar = true
	logger.Println(pg)
	logger.Println(s.ExecuteTemplate(w, "calendar", pg))
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "loginSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Retrieve our struct and type-assert it
	session.Values["username"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
