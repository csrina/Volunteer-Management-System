package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

/* Role constants */
const (
	_ = iota // We start our roles at 1 -- therefore ignore 0
	FACILITATOR
	TEACHER
	ADMIN
)

// tmp user struct just holds username and password
type TmpUser struct {
	Username string `json:"username" db:"username"`
	Password []byte `json:"password" db:"password"`
	UserID   int    `db:"user_id"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u TmpUser
	err := decoder.Decode(&u)
	if err != nil {
		logger.Println(err)
	}
	defer r.Body.Close()
	logger.Println("login request for user " + u.Username)
	session, err := store.New(r, "loginSession")
	if err != nil {
		logger.Println(err)
		return
	}
	session.Values["username"] = u.Username

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
	session.Values["role"] = role
	session.Save(r, w)
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
	var users []TmpUser
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
	fmt.Println(users[0].Password)
	fmt.Println(password)
	if err := bcrypt.CompareHashAndPassword(users[0].Password, password); err != nil {
		logger.Println(err)
		logger.Printf("'%v'\n", string(password))
		logger.Printf("'%v'\n", string(users[0].Password))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	logger.Printf("Password is correct for user %v\n", username)
}

func checkPassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u TmpUser
	err := decoder.Decode(&u)
	if err != nil {
		logger.Println(err)
	}
	defer r.Body.Close()
	session, err := store.Get(r, "loginSession")
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		session.Values["username"] = nil
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Retrieve our struct and type-assert it
	uname, ok := session.Values["username"].(string)
	if !ok {
		logger.Printf("error occured while retriving sesion value")
		return
	}
	role, ok := session.Values["role"].(int)
	if !ok {
		logger.Printf("error occured while retriving sesion value")
		return
	}
	auth(w, uname, u.Password, role)
}

func updatePassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u TmpUser
	err := decoder.Decode(&u)
	if err != nil {
		logger.Println(err)
	}
	defer r.Body.Close()
	session, err := store.Get(r, "loginSession")
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		session.Values["username"] = nil
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Retrieve our struct and type-assert it
	uname, ok := session.Values["username"].(string)
	if !ok {
		logger.Printf("error occured while retriving sesion value")
		return
	}
	passUpdate(w, uname, u.Password)
}

func passUpdate(w http.ResponseWriter, username string, password []byte) {
	encrypt_password, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not encrypt new password\n"))
		return
	}
	q := `update users
			SET password = $1
			WHERE username = $2`
	if _, err := db.Exec(q, string(encrypt_password), username); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not update password\n"))
		logger.Println(err)
		return
	}
	logger.Printf("Password updated")
	w.WriteHeader(http.StatusOK)
	logger.Printf("Password upadated for user %v\n", username)

}

func adminPassUpdate(w http.ResponseWriter, userID int, password []byte) {
	encrypt_password, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not encrypt new password\n"))
		return
	}

	q := `update users
			SET password = $1
			WHERE user_id = $2`
	if _, err := db.Exec(q, string(encrypt_password), userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not update password\n"))
		logger.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	logger.Printf("Password upadated for user %v\n", userID)

}
