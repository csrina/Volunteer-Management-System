package main

import (
	"encoding/json"
	"net/http"
)

type userFull struct {
	UserID    int    `json:"userId" db:"user_id"`
	UserRole  int    `json:"userRole" db:"user_role"`
	UserName  string `json:"userName" db:"username"`
	Password  string `json:"password" db:"password"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phoneNumber" db:"phone_number"`
}

func getUserList(w http.ResponseWriter, r *http.Request) {
	q := `SELECT user_id, user_role, last_name, first_name, username, email, phone_number
				FROM users`

	userList := []userFull{}
	err := db.Select(&userList, q)

	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(userList)
}

func loadAdminDash(w http.ResponseWriter, r *http.Request) {
	q := `SELECT family_id, family_name
			FROM family`

	families := []familyShort{}

	err := db.Select(&families, q)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, fam := range families {
		hours := familyHoursBooked(fam.FamilyID)
		families[i].WeekHours = hours
	}
	s := tmpls.Lookup("admindashboard.tmpl")
	err = s.ExecuteTemplate(w, "admindashboard", families)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
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
