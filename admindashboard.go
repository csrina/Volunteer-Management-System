package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type roomFull struct {
	RoomID     int           `json:"roomId" db:"room_id"`
	RoomName   string        `json:"roomName" db:"room_name"`
	TeacherID  int           `json:"teacherId" db:"teacher_id"`
	Children   sql.NullInt64 `json:"children" db:"children"`
	RoomNumber string        `json:"roomNum" db:"room_num"`
}

type roomDetailed struct {
	RoomID     int           `json:"roomId" db:"room_id"`
	RoomName   string        `json:"roomName" db:"room_name"`
	Teacher    string        `json:"teacher" db:"teacher"`
	Children   sql.NullInt64 `json:"children" db:"children"`
	RoomNumber string        `json:"roomNum" db:"room_num"`
}

type roomShort struct {
	RoomID   int    `json:"roomId" db:"room_id"`
	RoomName string `json:"roomName" db:"room_name"`
}

type userFull struct {
	UserID     int    `json:"userId" db:"user_id"`
	UserRole   int    `json:"userRole" db:"user_role"`
	UserName   string `json:"userName" db:"username"`
	Password   []byte `json:"password" db:"password"`
	FirstName  string `json:"firstName" db:"first_name"`
	LastName   string `json:"lastName" db:"last_name"`
	Email      string `json:"email" db:"email"`
	Phone      string `json:"phoneNumber" db:"phone_number"`
	BonusHours int    `json:"bonusHours" db:"bonus_hours"`
	BonusNote  string `json:"bonusNote" db:"bonus_note"`
}

type UserShort struct {
	UserID   int    `db:"user_id" json:"userId"`
	UserName string `db:"username" json:"userName"`
}

/*
 * Retrieves the user's first + last name and returns it
 */
func (u *UserShort) getFullName() (name string, err error) {
	q := `SELECT first_name || ' ' || last_name FROM users WHERE user_id = $1`
	if u.UserID > 0 {
		err = db.QueryRow(q, u.UserID).Scan(&name)
	} else if u.UserName != "" {
		err = db.QueryRow(q, u.UserName).Scan(&name)
	} else {
		// struct is empty
		return "", errors.New("Cannot retrieve name of unidentifiable user -- need UID or UName")
	}

	return // returns name, error via magical named return values
}

type familyFull struct {
	FamilyID   int    `json:"familyId" db:"family_id"`
	FamilyName string `json:"familyName" db:"family_name"`
	Children   int    `json:"children" db:"children"`
	Parents    []int  `json:"parents"`
	Dropped    []int  `json:"dropped"`
}

type familyDetailed struct {
	FamilyID   int         `json:"familyId" db:"family_id"`
	FamilyName string      `json:"familyName" db:"family_name"`
	Children   int         `json:"children" db:"children"`
	Parents    []UserShort `json:"parents"`
}

func createFamily(w http.ResponseWriter, r *http.Request) {
	family := familyFull{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&family)
	tx, err := db.Begin()
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := `INSERT INTO family (family_name, children)
			VALUES ($1, $2)
			RETURNING family_id`

	q2 := `UPDATE users
			SET family_id = $2
			WHERE user_id = $1`

	err = tx.QueryRow(q, family.FamilyName,
		family.Children).Scan(&family.FamilyID)

	if err != nil {
		tx.Rollback()
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, user := range family.Parents {
		_, err := tx.Exec(q2, user, family.FamilyID)
		if err != nil {
			tx.Rollback()
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	tx.Commit()
	w.WriteHeader(http.StatusCreated)
}

//TODO: update all users in recieved list
func updateFamily(w http.ResponseWriter, r *http.Request) {
	family := familyFull{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&family)
	tx, err := db.Begin()

	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := `UPDATE family
			SET family_name = $2, children = $3
			WHERE family_id = $1`

	q2 := `UPDATE users
			SET family_id = $2
			WHERE user_id = $1`

	q3 := `UPDATE users
			SET family_id = NULL
			WHERE family_id = $2
			AND user_id = $1`

	_, err = tx.Exec(q, family.FamilyID, family.FamilyName,
		family.Children)

	if err != nil {
		tx.Rollback()
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, user := range family.Parents {
		_, err := tx.Exec(q2, user, family.FamilyID)
		if err != nil {
			tx.Rollback()
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	for _, user := range family.Dropped {
		fmt.Println(user)
		_, err := tx.Exec(q3, user, family.FamilyID)
		if err != nil {
			tx.Rollback()
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

func basicRoomList(w http.ResponseWriter, r *http.Request) {
	rooms := []roomShort{}
	q := `SELECT room_id, room_name
			FROM room;`

	err := db.Select(&rooms, q)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rooms)
}

//gets all users not currently linked to a family
func lonelyFacilitators(w http.ResponseWriter, r *http.Request) {
	users := []UserShort{}
	q := `SELECT user_id, username
			FROM users
			WHERE family_id IS NULL
			AND user_role = 1`

	err := db.Select(&users, q)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(users)
}

func getTeachers(w http.ResponseWriter, r *http.Request) {
	users := []UserShort{}
	q := `SELECT user_id, username
			FROM users
			WHERE user_role = 2`
	err := db.Select(&users, q)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(users)
}

//does not return admin in this list
func getUserList(w http.ResponseWriter, r *http.Request) {
	options := r.URL.Query()
	userID, err := strconv.Atoi(options.Get("u"))
	//indicates we didnt have the flag or bad value
	if err != nil {
		q := `SELECT user_id, user_role, last_name, first_name, username, email, phone_number
				FROM users
				WHERE user_role != 3`
		userList := []userFull{}
		err := db.Select(&userList, q)

		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(userList)
	} else {
		q := `SELECT user_id, user_role, last_name, first_name, username, email, phone_number, bonus_hours, bonus_note
				FROM users
				WHERE user_id = ($1)`
		user := userFull{}
		err := db.QueryRowx(q, userID).StructScan(&user)

		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(user)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	newUser := userFull{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&newUser)
	fmt.Printf("%#v", newUser)

	newPass, err := bcrypt.GenerateFromPassword(newUser.Password, bcrypt.DefaultCost)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := `INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number, bonus_hours, bonus_note)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = db.Exec(q, newUser.UserRole, newUser.UserName, newPass,
		newUser.FirstName, newUser.LastName, newUser.Email, newUser.Phone,
		newUser.BonusHours, newUser.BonusNote)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	user := userFull{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)

	q := `UPDATE users 
			SET username = $2, first_name = $3, last_name = $4, email = $5, phone_number = $6, bonus_hours = $7, bonus_note = $8
			WHERE user_id = $1`

	_, err := db.Exec(q, user.UserID, user.UserName, user.FirstName, user.LastName, user.Email, user.Phone, user.BonusHours, user.BonusNote)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func removeFromFamily(w http.ResponseWriter, r *http.Request) {
	user := userFull{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)

	q := `UPDATE users
			SET family_id = NULL
			WHERE user_id = $1`

	_, err := db.Exec(q, user.UserID)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getFamilyList(w http.ResponseWriter, r *http.Request) {
	options := r.URL.Query()
	familyID, err := strconv.Atoi(options.Get("f"))
	if err != nil {
		q := `SELECT family_id, family_name, children
				FROM family`

		familyList := []familyFull{}
		err := db.Select(&familyList, q)

		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(familyList)
	} else {
		q2 := `SELECT family_id, family_name, children
				FROM family
				WHERE family_id = $1`
		q3 := `SELECT user_id, username
				FROM users
				WHERE family_id = $1`
		family := familyDetailed{}
		err := db.QueryRowx(q2, familyID).StructScan(&family)

		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.Select(&family.Parents, q3, familyID)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(family)
	}
}

func getClassInfo(w http.ResponseWriter, r *http.Request) {
	q := `SELECT room.room_id, room.room_name, users.username AS teacher, room.room_num
			FROM room, users
			WHERE room.teacher_id = users.user_id`
	classes := []roomDetailed{}

	err := db.Select(&classes, q)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(classes)
}

func createClass(w http.ResponseWriter, r *http.Request) {
	class := roomFull{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&class)

	q := `INSERT INTO room (room_name, teacher_id, room_num)
			VALUES ($1, $2, $3)`

	_, err := db.Exec(q, class.RoomName, class.TeacherID, class.RoomNumber)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateClass(w http.ResponseWriter, r *http.Request) {
	class := roomFull{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&class)

	q := `UPDATE room
			SET room_name = $2, teacher_id = $3, room_num = $4
			WHERE room_id = $1`

	_, err := db.Exec(q, class.RoomName, class.TeacherID, class.RoomNumber)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func loadAdminDash(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("admindashboard", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("admindashboard.tmpl")
	pg.DotJS = true
	s.ExecuteTemplate(w, "admindashboard", pg)
}

func loadAdminUsers(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("adminusers", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("adminusers.tmpl")
	pg.DotJS = true
	pg.MultiSelect = true
	s.ExecuteTemplate(w, "adminusers", pg)
}

func loadAdminCalendar(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("calendar", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("admincalendar.tmpl")
	pg.Calendar = true
	pg.DotJS = true
	pg.Toaster = true
	s.ExecuteTemplate(w, "admincalendar", pg)
}

func loadAdminReports(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("adminreports", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("adminreports.tmpl")
	pg.DotJS = true
	pg.Chart = true
	s.ExecuteTemplate(w, "adminreports", pg)
}

func loadAdminClasses(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("adminclasses", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s := tmpls.Lookup("adminclasses.tmpl")
	pg.DotJS = true
	s.ExecuteTemplate(w, "adminclasses", pg)
}
