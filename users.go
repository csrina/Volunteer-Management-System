package main

import (
	"net/http"
	"errors"
)

/*
 * Model of a row in the family table
 */
type Family struct {
	ID				int		`db:"family_id" json:"familyID"`
	Name			string	`db:"family_name" json:"familyName"`
	ParentOneID 	int		`db:"parent_one" json:"p1_id"`
	ParentTwoID 	int		`db:"parent_two" json:"p2_id"`
	Children    	int		`db:"children" json:"numKids"`
}

/*
 *  Retrieve family via userID contained in the request.
 */
func getFamilyViaRequest(r *http.Request) (*Family, error) {
	// get session
	sesh, _ := store.Get(r, "loginSession")
	username, ok := sesh.Values["username"].(string)
	if !ok {
		logger.Println("Invalid user token: ", username)
		return nil, errors.New("invalid token")
	}

	q := `SELECT family_id, family_name, COALESCE(parent_one, -1), COALESCE(parent_two, -1), children 
			FROM users NATURAL JOIN family 
			WHERE users.username = $1
				AND (family.parent_one = users.user_id 
						OR family.parent_two = users.user_id)`

	fdata := new(Family)
	err := db.QueryRow(q, username).Scan(
			&fdata.ID, &fdata.Name, &fdata.ParentOneID, &fdata.ParentTwoID, &fdata.Children)

	if err != nil {
		logger.Println(err)
		return nil, errors.New("could not retrieve family information")
	}
	logger.Println("Retrieved family: ", fdata)
	return fdata, nil
}

func getUID(r *http.Request) int {
	// get session
	sesh, _ := store.Get(r, "loginSession")
	username, ok := sesh.Values["username"].(string)
	if !ok {
		logger.Println("Invalid user token: ", username)
		return -1
	}

	q := `SELECT user_id FROM users WHERE username = $1`
	var uid int
	err := db.QueryRow(q, username).Scan(&uid)
	if err != nil {
		return -1
	}
	return uid
}
