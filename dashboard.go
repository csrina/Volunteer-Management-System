package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/now"
)

/*
 * dashboard.js expects the following struct in json format
 */
 type DashData struct {
 	 HoursGoal   		float64  			`json:"hoursGoal"`
	 HoursBooked 		float64  			`json:"hoursBooked"`
	 HoursDone   		float64  			`json:"hoursDone"`
	 History1     		ChartDataSet 		`json:"history1"` // historical hours completed/week for interval
	 History2			ChartDataSet		`json:"history2"` // same for parent two of family
	 StartOfPeriod   	time.Time       	`json:"startOfPeriod"` // start date for chart
	 EndOfPeriod    	time.Time       `json:"endOfPeriod`  // end date for chart
}

func (dd *DashData) setHoursGoal(numKids int) {
	if numKids == 1 {
		dd.HoursGoal = ONE_CHILD_HOURS_GOAL
	} else {
		dd.HoursGoal = DEFAULT_HOURS_GOAL
	}
}
/* START OF SCHOOL YEAR FOR HISTORY TRACKING */
const (
	PERIOD_LENGTH = 3 // months
	ONE_CHILD_HOURS_GOAL = 5.00
	DEFAULT_HOURS_GOAL = 7.50
)

 /*
  * Corresponding to a single object, in the datasets array, of a chart.js chart.
  * We use it for the historical hourly attendance for a family.
  */
 type ChartDataSet struct {
	Label 		string         `json:"label"`       // Dataseries name
	Data  		[]DurationPoint `json:"data"`        // array of data points
	Fill    	bool            `json:"fill"`        // do we fill area under line (or within the bars)?
	BorderColor string          `json:"borderColor"` // really the colour colour
	Tension		float64      	`json:"lineTension"` // 0 = no curvyness (no interp.) 1.00 max curvyness
	Stepped     string          `json:"steppedLine"`
	SpanGaps    bool            `json:"spanGaps"`
}

// For charting
type DurationPoint struct {
	X 	time.Time		`json:"t"`
	Y 	float64     	`json:"y"`
}

func (c ChartDataSet) configureAsHistoricalHours(label, colour string, fill bool, tension float64) ChartDataSet {
	c.Label = label
	c.Fill = fill
	c.BorderColor = colour
	c.Tension = tension
	c.Stepped = "after"
	c.SpanGaps = false
	return c // for chaining
}

/*
 *  Parents may have multiple facillitations in a day, which when charting -- leads to weirdness
 *  this function will scan to see if the x-value (date) exists before appending the new point.
 *
 *   If the x-value exits, the y-value (duration) is simply added to.
 */
func (c *ChartDataSet) addDurationPoint(p DurationPoint) *ChartDataSet {
	for i, point := range c.Data {
		// Add duration to existing point's duration (Y) value
		if p.X.YearDay() == point.X.YearDay() {
			point.Y += p.Y
			c.Data[i] = point;
			p.Y = 0  // Flag to prevent adding another point for this date
			break
		}
	}
	// If p.Y has been set to 0, dont add the point
	if p.Y > 0 {
		c.Data = append(c.Data, p)
	}
	return c
}

 /* Replacement for dashboardData
  * which delegates most responsibility to functions
  */
func getDashData(w http.ResponseWriter, r *http.Request) {
	family, err := getFamilyViaRequest(r)
	if err != nil {
		logger.Println("Failed to retrieve family")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dd := new(DashData)
	dd, err = dd.updateHoursData(*family, time.Now()) // Get the relevant data
	if err != nil {
		logger.Println("Failed to update hours using family")
		logger.Println("Family: ", family)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dd.History1 = dd.History1.configureAsHistoricalHours("P1", "#FAD201", false, 0.0)
	dd.History2 = dd.History2.configureAsHistoricalHours("P2", "#201FAD", false, 0.0)
	logger.Println("DASHDATA: ", dd)
	logger.Println("P1 CDS: ", dd.History1)
	logger.Println("P2 CDS: ", dd.History2)


	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	encoder.Encode(dd)
}

/*
 * Creates a DashData struct and fills it with data given a family, and the reference date: today
 * Get the hours data relative to the day passed as today.
 * The history will be tracked from the FIRST_MONTH and FIRST_DAY
 * of today's year.
 */
func (dd *DashData) updateHoursData(fam Family, today time.Time) (*DashData, error) {
	dd.setHoursGoal(fam.Children)
	if today.Weekday() == time.Sunday {
		today = today.AddDate(0, 0, 1) // move to monday so we reference next week
	} else if today.Weekday() == time.Saturday {
		today = today.AddDate(0, 0, 2)
	}
	now.FirstDayMonday = true
	nowToday := now.New(today)
	startOfWeek := nowToday.BeginningOfWeek()
	dd.EndOfPeriod = nowToday.EndOfWeek()
	logger.Println("T: ", today, "SOW: ", startOfWeek, " EOW: ", dd.EndOfPeriod)

	dd.StartOfPeriod = today.AddDate(0, -PERIOD_LENGTH, 0) // Go back 4 months
	/* Get all bookings relevant */
	bookings, err := fam.getFamilyBookings(dd.StartOfPeriod, dd.EndOfPeriod)
	if (err != nil) {
		return nil, err
	}
	for _, b := range bookings {
		duration := (b.End.Sub(b.Start).Hours() * float64(b.Modifier))
		// historical bookings must be separated by parent
		if b.Start.Before(startOfWeek) {
			if (b.UserID == fam.ParentOneID) {
				dd.History1.addDurationPoint(DurationPoint{Y: duration, X: now.New(b.Start).BeginningOfDay()})
			} else {
				dd.History2.addDurationPoint(DurationPoint{Y: duration, X: now.New(b.Start).BeginningOfDay()})
			}
		// Family hours are conglomerated in the totals
		} else if b.Start.Before(nowToday.Time) && b.End.After(startOfWeek) {
			dd.HoursDone += duration
			dd.HoursBooked += duration // Even though they're done, theyre still booked 4 this week
		} else {
			dd.HoursBooked += duration // Time is after today, but before week end --> booked hours
		}
	}

	return dd, nil
}

/*
join family and users

join booking and time block

pull out username, user_id, block_start, block_end
*/
type WeeklyBooking struct {
	UserName   string    `db:"username"`
	UserId     int       `db:"user_id"`
	BlockStart time.Time `db:"block_start"`
	BlockEnd   time.Time `db:"block_end"`
	Children   int       `db:"children"`
}

type FriendlyFormat struct {
	Eventlist   []string `json:"eventlist"`
	HoursBooked float64  `json:"hoursBooked"`
	HoursDone   float64  `json:"hoursDone"`
	Children    int      `json:"children"`
}

func startOfWeek(current time.Time) time.Time {
	layoutDay := "Mon"
	check := current.Format(layoutDay)
	for check != "Mon" {
		current = current.AddDate(0, 0, -1)
		check = current.Format(layoutDay)
	}
	return current
}

func dashboardData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}
	// make session stuff
	user := getUID(r)
	role, err := getRoleNum(r)
	logger.Printf("User ID found %d\n\n\n", user)
	start := startOfWeek(time.Now())
	var q string
	if role == 1 {
		q = `SELECT username, s.user_id, block_start, block_end, children 
FROM (
SELECT username, user_id, children
FROM family INNER JOIN users
ON user_id = parent_one or user_id = parent_two) r
INNER JOIN
(SELECT user_id, block_start, block_end
FROM booking b INNER JOIN time_block t 
ON b.block_id = t.block_id) s
ON r.user_id = s.user_id
WHERE r.user_id = $1 AND block_start > $2 AND block_start < $3`
	} else {
		fmt.Fprintln(w, "non facilitator doesnt have dashboard right now")
		return
	}

	var bookings []WeeklyBooking

	err = db.Select(&bookings, q, user, start, start.AddDate(0, 0, 6))
	if err != nil {
		fmt.Printf("%v", err)
	}
	var friendly FriendlyFormat
	i := 0
	hoursDone := 0.0
	hoursBooked := 0.0
	layout := "Mon Jan 2 15:04"
	if len(bookings) == 0 {
		fmt.Fprintln(w, "no results")
		return
	}
	for each := range bookings {

		if bookings[each].BlockEnd.Before(time.Now()) {
			hoursDone += bookings[each].BlockEnd.Sub(bookings[each].BlockStart).Hours()
		}
		hoursBooked += bookings[each].BlockEnd.Sub(bookings[each].BlockStart).Hours()
		if bookings[each].BlockEnd.After(time.Now()) {
			friendly.Eventlist = append(friendly.Eventlist, bookings[each].BlockStart.Format(layout)+" to ")
			friendly.Eventlist[i] = friendly.Eventlist[i] + bookings[each].BlockEnd.Format(layout)
			i++
		}
	}
	friendly.HoursBooked = hoursBooked
	friendly.HoursDone = hoursDone
	friendly.Children = bookings[0].Children

	encoder := json.NewEncoder(w)
	logger.Println("FRIENDLY:  ", friendly)
	err = encoder.Encode(friendly)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
