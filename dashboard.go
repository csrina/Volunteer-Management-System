package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// For charting
type DurationPoint struct {
	X time.Time `json:"t"`
	Y float64   `json:"y"`
}

/* for sorting duration points by time */
type DurationPoints []DurationPoint
func (s DurationPoints) Less(i, j int) bool { return s[i].X.Before(s[j].X) }
func (s DurationPoints) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s DurationPoints) Len() int           { return len(s) }

/*
 * Corresponding to a single object, in the datasets array, of a chart.js chart.
 * We use it for the historical hourly attendance for a family.
 */
type ChartDataSet struct {
	Label       string         `json:"label"`       // Dataseries name
	Data        DurationPoints `json:"data"`        // array of data points
	Fill        bool           `json:"fill"`        // do we fill area under line (or within the bars)?
	BorderColor string         `json:"borderColor"` // really the colour colour
	Tension     float64        `json:"lineTension"` // 0 = no curvyness (no interp.) 1.00 max curvyness
	Stepped     string         `json:"steppedLine"`
	SpanGaps    bool           `json:"spanGaps"`
}

/* default dashboard chart settings */
func (c *ChartDataSet) configureAsHistoricalHours(label, colour string, fill bool, tension float64) *ChartDataSet {
	c.Label = label
	c.Fill = fill
	c.BorderColor = colour
	c.Tension = tension
	c.Stepped = "before"
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

	fd := new(FamilyData)
	err = fd.init(family, time.Now()) // initialize dashboard for family
	if err != nil {
		logger.Println("Failed to initialize family data")
		logger.Println("Family: ", family)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fd.spanHistoryGaps() // Fill in some 0 days for better charting

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	encoder.Encode(fd)
}
