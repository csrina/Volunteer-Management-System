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
	Parents			[]*User	`json:"parents"`
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

	q := `SELECT family_id, family_name, children 
			FROM users NATURAL JOIN family 
			WHERE users.username = $1
				AND family.family_id = users.family_id`

	fdata := new(Family)
	err := db.QueryRow(q, username).Scan(
			&fdata.ID, &fdata.Name, &fdata.Children)
	if err != nil {
		logger.Println(err)
		return nil, errors.New("could not retrieve family information")
	}

	var uids []int
	q = `SELECT user_id FROM users WHERE users.family_id = $1`
	err = db.Select(&uids, q, fdata.ID)
	if err != nil {
		return fdata, err
	}
	// Make slice for parents
	fdata.Parents = make([]*User, len(uids), len(uids) + 1)
	for _, uid := range uids {
		u := new(User)
		err = u.init(uid)
		if err != nil {
			logger.Println("Error creating user from uid in getFamilyViaRequest:  " + err.Error())
			continue
		}
		// Add user who belongs to family to slice
		fdata.Parents = append(fdata.Parents, u)
	}
	return fdata, nil
}

/*
 * User sans password field
 */
type User struct {
	UserID 		int
	Role 		int
	Username 	string
	FirstName 	string
	LastName 	string
	Email 		string
	Phone 		string
	FamilyID 	int
	Bonus		float64
	BonusNote	string
}

/*
 * Initializes reciever struct based on the given UID, a user from the db.
 */
func (u *User) init(UID int) error {
	q := `SELECT 	user_id, user_role, username, first_name, last_name, 
					email, phone_number, family_id, bonus_hours, bonus_note
			FROM users 
			WHERE user_id = $1`
	err := db.QueryRow(q, UID).Scan(&u.UserID, &u.Role, &u.FirstName, &u.LastName,
									&u.Email, &u.Phone, &u.FamilyID, &u.Bonus, &u.BonusNote)
	if err != nil {
		return err
	}
	return nil
}


func getUID(r *http.Request) (UID int) {
	// get session
	sesh, _ := store.Get(r, "loginSession")
	username, ok := sesh.Values["username"].(string)
	if !ok {
		logger.Println("Invalid user token: ", username)
		return -1
	}
	q := `SELECT user_id FROM users WHERE username = $1`
	err := db.QueryRow(q, username).Scan(&UID)
	if err != nil {
		return -1
	}
	return UID
}

/* Given a UID, get the FID which the user belongs to */
func getUsersFID(userID int) (int, error) {
	FID := -1
	q := `SELECT family_id FROM users WHERE users.user_id = $1 `
	err := db.QueryRow(q, userID).Scan(&FID)
	if err != nil {
		logger.Println("error retrieving fid")
		return -1, err
	}
	return FID, nil
}
