package main

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
}

func renderCalendar(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "F Schedule"}
	t, _ := template.ParseFiles("views/calendar.gohtml")
	t.Execute(w, p)
}
