package main

import (
	"net/http"
)

func loadAdminDash(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("admindashboard", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("admindashboard.tmpl")
	s.ExecuteTemplate(w, "admindashboard", pg)
}

func loadAdminUsers(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("adminusers", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("adminusers.tmpl")
	s.ExecuteTemplate(w, "adminusers", pg)
}

func loadAdminCalendar(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("admincalendar", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("admincalendar.tmpl")
	s.ExecuteTemplate(w, "admincalendar", pg)
}

func loadAdminReports(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("adminreports", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("adminreports.tmpl")
	s.ExecuteTemplate(w, "adminreports", pg)
}

func loadAdminClasses(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("adminclasses", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("adminclasses.tmpl")
	s.ExecuteTemplate(w, "adminclasses", pg)
}
