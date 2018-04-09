package main /* Stores data for filling templates */
import (
	"net/http"
	"strings"
)

type Message struct {
	MsgID int    `db:"msg_id" json:"msgid"`
	Msg   string `db:"msg" json:"msg"`
}

type Page struct {
	PageName string
	Role     string
	Username string
	Room     string // For teachers, potentially others
	Messages []Message

	/* Dependency flags for templates */
	Calendar    	bool // page has calendar --> set flag to true
	Chart       	bool // page requires chart.js
	Dashboard   	bool //is it the dahsboard?
	DotJS      		bool
	Toaster     	bool
	MultiSelect 	bool // page requires multiselectc css and js
}

/*
 * loads a Page struct with data from the request & returns ptr to it
 * --> requires user to be logged in
 */
func loadPage(pn string, r *http.Request) (*Page, error) {
	data := &Page{
		PageName: pn,
	}

	/* Get user name for filling in template too */
	sesh, _ := store.Get(r, "loginSession")
	uname, ok := sesh.Values["username"].(string)
	if !ok {
		return nil, &ClientSafeError{Msg: "invalid username"}
	}
	data.Username = uname
	/* get user's role */
	role, err := getRoleNum(r)
	if err != nil {
		return nil, err
	}
	switch role {
	case FACILITATOR:
		data.Role = "Facilitator"
		data.Messages, err = getNotifications(data.Username)
		if err != nil {
			return nil, err
		}
	case TEACHER:
		data.Role = "Teacher"
	case ADMIN:
		data.Role = "Admin"
	default:
		return nil, &ClientSafeError{Msg: "insufficient access rights"}
	}

	return data, nil
}

func getNotifications(uname string) ([]Message, error) {
	var msgs []Message
	q := `SELECT notify.msg_id, notifications.msg
			from notify, notifications, users
			where users.user_id = notify.user_id 
				AND notify.viewed = 'f' 
				AND users.username= $1 
				AND notifications.msg_id = notify.msg_id`

	err := db.Select(&msgs, q, uname)
	if err != nil {
		logger.Println(err)
		return msgs, err
	}

	return msgs, nil

}

/* Retrieves the role of the requesting party */
func getRoleNum(r *http.Request) (int, error) {
	sesh, err := store.Get(r, "loginSession")
	if err != nil {
		return -1, err
	}
	uname, ok := sesh.Values["username"].(string)
	if !ok {
		return -1, &ClientSafeError{Msg: "username of session invalid type"}
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

/*** PAGE LOADS FOR GET REQUESTS ***/
// Initial select role form
func loadMainLogin(w http.ResponseWriter, r *http.Request) {
	s := tmpls.Lookup("mainLogin.tmpl")
	p := &Page{
		PageName: "mainLogin",
		Role:     "",
		Username: "",
		Calendar: false,
	}
	s.ExecuteTemplate(w, "content", p)
}

// Actual login form
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
	p := &Page{
		PageName: "login",
		Role:     title,
		Username: "",
	}
	s := tmpls.Lookup("login.tmpl")
	s.ExecuteTemplate(w, "loginForm", p)
}

// Facilitator dash
func loadDashboard(w http.ResponseWriter, r *http.Request) {
	role, err := getRoleNum(r)
	if err != nil {
		if _, ok := err.(*ClientSafeError); ok {
			http.Error(w, err.Error(), http.StatusBadGateway)
		} else {
			http.Error(w, "Something funny happened, sorry. Please try again ", http.StatusInternalServerError)
		}
		return
	}

	if role == TEACHER {
		http.Redirect(w, r, "/teacher", http.StatusFound)
	} else if role == ADMIN {
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
	}

	pg, err := loadPage("dashboard", r) // load page
	if err != nil {
		if _, ok := err.(*ClientSafeError); ok {
			http.Error(w, err.Error(), http.StatusBadGateway)
		} else {
			http.Error(w, "Something funny happened, sorry. Please try again ", http.StatusInternalServerError)
		}
		return
	}
	s := tmpls.Lookup("dashboard.tmpl")
	// dependency flags for dashboard
	pg.Calendar = true
	pg.Chart = true
	pg.Dashboard = true
	pg.Toaster = true
	s.ExecuteTemplate(w, "dashboard", pg) // include page struct
}

// change pw
func loadPassword(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("password", r) // load page
	if err != nil {
		if _, ok := err.(*ClientSafeError); ok {
			http.Error(w, err.Error(), http.StatusBadGateway)
		} else {
			http.Error(w, "Something funny happened, sorry. Please try again ", http.StatusInternalServerError)
		}
		return
	}
	pg.Toaster = true
	s := tmpls.Lookup("password.tmpl")
	s.ExecuteTemplate(w, "password", pg) // include page struct
}

// Calendar (Facilitator)
func loadCalendar(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("calendar", r)
	if err != nil {
		if _, ok := err.(*ClientSafeError); ok {
			http.Error(w, err.Error(), http.StatusBadGateway)
		} else {
			http.Error(w, "Something funny happened, sorry. Please try again ", http.StatusInternalServerError)
		}
		return
	}
	s := tmpls.Lookup("calendar.tmpl")
	// calendar dependency flag
	pg.Calendar = true
	pg.Toaster = true
	s.ExecuteTemplate(w, "calendar", pg)
}

// Teacher dashboard
func loadTeacher(w http.ResponseWriter, r *http.Request) {
	pg, err := loadPage("teacher", r)
	if err != nil {
		if _, ok := err.(*ClientSafeError); ok {
			http.Error(w, err.Error(), http.StatusBadGateway)
		} else {
			http.Error(w, "Something funny happened, sorry. Please try again ", http.StatusInternalServerError)
		}
		return
	}
	s := tmpls.Lookup("teacher.tmpl")
	pg.Calendar = true
	pg.Toaster = true
	s.ExecuteTemplate(w, "teacher", pg)
}
