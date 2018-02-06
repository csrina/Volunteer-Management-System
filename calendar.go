package main

import (
	"html/template"
	"net/http"
)

/* Will eventually be needed */
type Page struct {
	Title string
}

/* Renders the calendar found in views/calendar.gohtml */
func renderCalendar(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "F Schedule"}
	t, _ := template.ParseFiles("views/calendar.gohtml")
	t.Execute(w, p)
}
