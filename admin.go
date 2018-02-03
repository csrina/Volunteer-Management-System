package main

import (
	"net/http"
	"time"
)

func calSetup(w http.ResponseWriter, r *http.Request) {
	//blocks have to be filled by either json decode or from forms
	//start / end also need to be pulled from the request body
	//blocks := []Timeblock{}
}

func applyTemplate(startDate time.Time, stopDate time.Time, blocks []TimeBlock) {
	Tx, err := db.Begin()
	if err != nil {
		logger.Fatal("Could not start database transaction")
	}
	for startDate.Before(stopDate) {
		for _, block := range blocks {
			switch block.Start.Weekday() {
			case time.Monday:
				tempDate := startDate.AddDate(0, 0, 0)
				block.setDay(tempDate)
				err := block.insertBlock()
				if err != nil {
					Tx.Rollback()
				}
				break
			case time.Tuesday:
				tempDate := startDate.AddDate(0, 0, 1)
				block.setDay(tempDate)
				err := block.insertBlock()
				if err != nil {
					Tx.Rollback()
				}
				break
			case time.Wednesday:
				tempDate := startDate.AddDate(0, 0, 2)
				block.setDay(tempDate)
				err := block.insertBlock()
				if err != nil {
					Tx.Rollback()
				}
				break
			case time.Thursday:
				tempDate := startDate.AddDate(0, 0, 3)
				block.setDay(tempDate)
				err := block.insertBlock()
				if err != nil {
					Tx.Rollback()
				}
				break
			case time.Friday:
				tempDate := startDate.AddDate(0, 0, 4)
				block.setDay(tempDate)
				err := block.insertBlock()
				if err != nil {
					Tx.Rollback()
				}
				break
			}
		}
		startDate = startDate.AddDate(0, 0, 7)
	}
	Tx.Commit()
}

func (tb *TimeBlock) setDay(startDate time.Time) {
	tb.Start = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), tb.Start.Hour(), tb.Start.Minute(), 0, 0, tb.Start.Location())
	tb.End = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), tb.End.Hour(), tb.End.Minute(), 0, 0, tb.Start.Location())
}

// DO NOT USE THIS UNLESS 100% REQUIRED
func undoSetup() {
}
