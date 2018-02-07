package main

import (
	"html/template"
	"net/http"
)

/* Will eventually be needed */
type Page struct {
	Title string
}

func renderView(w http.ResponseWriter, view string, p *Page) {
	t, _ := template.ParseFiles(view + ".gohtml")
	t.Execute(w, p)
}

/* Renders the calendar found in views/calendar.gohtml */
func renderCalendar(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "F Schedule"}
	renderView(w, "views/calendar", p)
}

/* renders the dashboard for a user */
func renderDashboard(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "F Dashboard"}
	renderView(w, "views/dashboard", p)
}
