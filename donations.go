package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Donor interface {
	GiveCharity(donee *Family, amount float64) (*Donation, error)
	GetID() int
}

/*
 * Donations for things that satisfy the interfaces
 */
type Donation struct {
	ID        int       `db:"donation_id" json:"id"`
	Sender    Donor     `db:"donor_id" json:"donor"`
	Recipient Donor     `db:"donee_id" json:"donee"`
	Amount    float64   `db:"amount" json:"amount"`
	Date      time.Time `db:"date_sent" json:"date"`
}

func (d *Donation) isLegal() (legal bool, err error) {
	// check Donation.DAte for weekend --> if yes, modify to use last week's hoursDone
	// else use the current week hoursDone
	// We have the hours now
	// hoursDone < donation amount ===> bail
	// donor id OR donee id invalid ==> bail
	return
}

// save donation in db, update with returned id
func (d *Donation) save() error {
	q := `INSERT INTO donation (donor_id, donee_id, amount)
				VALUES ($1, $2, $3) RETURNING donation_id, date_sent`

	return db.QueryRow(q, strconv.Itoa(d.Sender.GetID()), strconv.Itoa(d.Recipient.GetID()), d.Amount).Scan(&d.ID, &d.Date)
}

type Donations []Donation

/* Get the total net amount for all of the donations in this slice */
func (gifts Donations) netAmount() (netAmount float64) {
	for _, g := range gifts {
		netAmount += g.Amount
	}
	return netAmount
}

func donatePostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := r.GetBody()
	if err != nil {
		logger.Println(err)
		http.Error(w, "Improperly formatted request body", http.StatusBadRequest)
	}
	dec := json.NewDecoder(body)
	donation := new(Donation)
	err = dec.Decode(donation)
	if err != nil {
		http.Error(w, "Couldn't parse request", http.StatusBadRequest)
	}
	// ensure donor family has the funds
	ok, err := donation.isLegal()
	if err != nil {
		http.Error(w, "Couldn't complete donation", http.StatusBadRequest)
		return
	} else if !ok {
		http.Error(w, "Insufficient hours to meet donation requirement", http.StatusBadRequest)
		return
	}
	donation.save()
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	enc.Encode(donation) // send back with id for toaster
}
