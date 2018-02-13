package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// tmp user struct just holds username and password
type User struct {
	Username string `json:"username" db:"username"`
	Password []byte `json:"password" db:"password"`
	UserID   int    `db:"user_id"`
}

func loadMainLogin(w http.ResponseWriter, r *http.Request) {
	s := tmpls.Lookup("mainLogin.tmpl")
	s.ExecuteTemplate(w, "content", nil)
}
func loadLogin(w http.ResponseWriter, r *http.Request) {
	var title string
	cur := r.URL.Path
	if strings.Contains(cur, "facilitator") {
		title = "Facilitator "
	} else if strings.Contains(cur, "teacher") {
		title = "Teacher "
	} else if strings.Contains(cur, "admin") {
		title = "Admin "
	}
	title = title + "Login"
	s := tmpls.Lookup("login.tmpl")
	s.ExecuteTemplate(w, "loginForm", title)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		logger.Println(err)
	}
	defer r.Body.Close()
	logger.Println("login request for user " + u.Username)
	cur := r.URL.Path
	var role int
	if strings.Contains(cur, "facilitator") {
		role = 1
	} else if strings.Contains(cur, "teacher") {
		role = 2
	} else if strings.Contains(cur, "admin") {
		role = 3
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	auth(w, u.Username, u.Password, role)
}

/* genreate password hash to work with the auth below
   userPassword1 := "some user-provided password"
   // Generate "hash" to store from user password
   hash, err := bcrypt.GenerateFromPassword([]byte(userPassword1), bcrypt.DefaultCost)
   if err != nil {
       // TODO: Properly handle error
       log.Fatal(err)
   }
   fmt.Println("Hash to store:", string(hash))
   // Store this "hash" somewhere, e.g. in your database
*/

func auth(w http.ResponseWriter, username string, password []byte, role int) {
	q := `SELECT username, password, user_id
			FROM users
			where username = $1 AND user_role = $2`
	logger.Println("Checking if user exits")
	logger.Println(q + " " + username)
	users := []User{}
	if err := db.Select(&users, q, username, role); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Println(err)
		return
	}
	count := len(users)
	if count == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Printf("User %v not found.\n", username)
		return
	}
	logger.Printf("User %v found\n", username)
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword(users[0].Password, password); err != nil {
		logger.Printf("'%v'\n", string(password))
		logger.Printf("'%v'\n", string(users[0].Password))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    "SessionToken",
		Value:   strconv.Itoa(users[0].UserID),
		Expires: expire,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusAccepted)
	logger.Printf("Password is correct for user %v\n", username)
}
