package main

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/jinzhu/now"
)

/*
 * dashboard.js expects the following struct in json format
 */
type DashData struct {
	HoursGoal     float64      `json:"hoursGoal"`
	HoursBooked   float64      `json:"hoursBooked"`
	HoursDone     float64      `json:"hoursDone"`
	History1      ChartDataSet `json:"history1"`      // historical hours completed/week for interval
	History2      ChartDataSet `json:"history2"`      // same for parent two of family
	StartOfPeriod time.Time    `json:"startOfPeriod"` // start date for chart
	EndOfPeriod   time.Time    `json:"endOfPeriod"`   // end date for chart
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
	PERIOD_LENGTH        = 3 // months
	ONE_CHILD_HOURS_GOAL = 2.50
	DEFAULT_HOURS_GOAL   = 5.00
)

/*
 * Corresponding to a single object, in the datasets array, of a chart.js chart.
 * We use it for the historical hourly attendance for a family.
 */
type ChartDataSet struct {
	Label       string          `json:"label"`       // Dataseries name
	Data        []DurationPoint `json:"data"`        // array of data points
	Fill        bool            `json:"fill"`        // do we fill area under line (or within the bars)?
	BorderColor string          `json:"borderColor"` // really the colour colour
	Tension     float64         `json:"lineTension"` // 0 = no curvyness (no interp.) 1.00 max curvyness
	Stepped     string          `json:"steppedLine"`
	SpanGaps    bool            `json:"spanGaps"`
}

// For charting
type DurationPoint struct {
	X time.Time `json:"t"`
	Y float64   `json:"y"`
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
			c.Data[i] = point
			p.Y = 0 // Flag to prevent adding another point for this date
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
	todaySaved := now.New(today)        // If its a weekend, we need this saved value for later
	if today.Weekday() == time.Sunday { // weekend days must be shifted to monday
		today = today.AddDate(0, 0, 1) // move to monday so we reference next week
	} else if today.Weekday() == time.Saturday {
		today = today.AddDate(0, 0, 2)
	}
	now.Monday()
	nowToday := now.New(today)                // We use this to determine start of week, so it should be the adjusted today
	startOfWeek := nowToday.BeginningOfWeek() // if weekend, this is next week's monday
	dd.EndOfPeriod = nowToday.EndOfWeek()

	dd.StartOfPeriod = today.AddDate(0, -PERIOD_LENGTH, 0) // Go back 4 months
	/* Get all bookings relevant */
	bookings, err := fam.getFamilyBookings(dd.StartOfPeriod, dd.EndOfPeriod)
	if err != nil {
		return nil, err
	}
	for _, b := range bookings {
		duration := (b.End.Sub(b.Start).Hours() * float64(b.Modifier))
		// historical bookings must be separated by parent
		if b.Start.Before(startOfWeek) {
			if b.UserID == fam.ParentOneID {
				dd.History1.addDurationPoint(DurationPoint{Y: duration, X: now.New(b.Start).BeginningOfDay()})
			} else {
				dd.History2.addDurationPoint(DurationPoint{Y: duration, X: now.New(b.Start).BeginningOfDay()})
			}
			// Family hours are conglomerated in the totals
		} else if b.Start.Before(todaySaved.Time) && b.End.After(startOfWeek) {
			dd.HoursDone += duration
			dd.HoursBooked += duration // Even though they're done, theyre still booked 4 this week
		} else {
			dd.HoursBooked += duration // Time is after today, but before week end --> booked hours
		}
	}
	dd.HoursBooked = math.Trunc(dd.HoursBooked*100) / 100
	dd.HoursDone = math.Trunc(dd.HoursDone*100) / 100
	return dd, nil
}
