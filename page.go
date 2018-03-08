package main /* Stores data for filling templates */
import (
	"errors"
	"net/http"
)

type Page struct {
	PageName string
	Role     string
	Username string
	/* Dependency flags for templates */
	Calendar bool // page has calendar --> set flag to true
	Chart bool // page requires chart.js
}

/*
 * loads a Page struct with data from the request & returns ptr to it
 * --> requires user to be logged in
 */
func loadPage(pn string, r *http.Request) (*Page, error) {
	data := &Page{
		PageName: pn,
	}
	/* get user's role */
	role, err := getRoleNum(r)
	if err != nil {
		return nil, err
	}

	switch role {
	case FACILITATOR:
		data.Role = "Facilitator"
	case TEACHER:
		data.Role = "Teacher"
	case ADMIN:
		data.Role = "Admin"
	default:
		return nil, errors.New("insufficient access rights")
	}

	/* Get user name for filling in template too */
	sesh, _ := store.Get(r, "loginSession")
	uname, ok := sesh.Values["username"].(string)
	if !ok {
		return nil, errors.New("invalid username")
	}
	data.Username = uname
	return data, nil
}

/* Retrieves the role of the requesting party */
func getRoleNum(r *http.Request) (int, error) {
	sesh, err := store.Get(r, "loginSession")
	if err != nil {
		return -1, err
	}
	uname, ok := sesh.Values["username"].(string)
	if !ok {
		return -1, errors.New("username of session invalid type")
	}
	/* Get and return the role */
	var role int
	q := `SELECT user_role FROM users WHERE (username = $1)`
	err = db.QueryRow(q, uname).Scan(&role)
	if err != nil {
		return -1, err
	}
	return role, nil
}


/* CompareName returns true when PageName == comparator */
func (p Page) CompareName(comparator string) bool {
	return p.PageName == comparator
}
