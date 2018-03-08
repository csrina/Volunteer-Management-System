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
