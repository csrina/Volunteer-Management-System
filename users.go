package main

import (
	"math"
	"net/http"
	"sort"
	"time"

	"database/sql"

	"github.com/jinzhu/now"
	"github.com/lucasb-eyer/go-colorful"
)

/*
 * Model of a row in the family table
 */
type Family struct {
	ID       int     `db:"family_id" json:"familyID"`
	Name     string  `db:"family_name" json:"familyName"`
	Parents  []*User `json:"parents"`
	Children int     `db:"children" json:"numKids"`
}

/* Short-form, used in reporting.go */
type familyShort struct {
	FamilyID   int     `json:"familyId" db:"family_id"`
	FamilyName string  `json:"familyName" db:"family_name"`
	WeekHours  float64 `json:"weekHours"`
	Children   int     `json:"children" db:"children"`
}

/* For monthly reporting */
type familyMonth struct {
	FamilyID   int       `json:"familyId" db:"family_id"`
	FamilyName string    `json:"familyName" db:"family_name"`
	Weeks      []float64 `json:"weeks"`
	Month      float64   `json:"month"`
	Children   int       `json:"children" db:"children"`
}

// Data for a family's dashboard or for other repurposing
type FamilyData struct {
	family       *Family
	HoursGoal    float64   				`json:"hoursGoal"`   // weekly goal
	HoursBooked  float64   				`json:"hoursBooked"` // hours to booked in
	HoursDone    float64   				`json:"hoursDone"`   // actually completed hrs
	// Where historical keys are the parent UID
	History       map[int]*ChartDataSet `json:"history"`       // datasets are facillitators + a donation dataset
	StartOfPeriod time.Time             `json:"startOfPeriod"` // start date for chart
	EndOfPeriod   time.Time             `json:"endOfPeriod"`   // end date for chart
	StartOfWeek   time.Time 			`json:"startMoment"`				  //  start of current week (potentially adjusted to next monday if current time is a weekend)
}

func (fd *FamilyData) GetDonee() Donee { return fd.family }
func (fd *FamilyData) GetHours() (hours float64) { return fd.HoursDone }
func (fd *FamilyData) GetStartDate() time.Time { return fd.StartOfPeriod }
func (fd *FamilyData) GetStartWeek() time.Time { return fd.StartOfWeek }
func (fd *FamilyData) GetEndDate() time.Time { return fd.EndOfPeriod }
func (fd *FamilyData) SetHours(hours float64) { fd.HoursDone = hours }

func (fd *FamilyData) setHoursGoal(numKids int) {
	fd.HoursGoal = getHourGoal(numKids)
}

/*
 * Get the hours which can be donated (i.e. hoursDone - hoursGoal), adjusting for weekends.
 * Unlike the dashboard default (if weekend monday is next monday); Donations defaults to last week,
 * until monday has arrived. Completed hours differs from Hours Done only if the current time (@ instance access)
 * falls on a weekend.
 */
func (fd *FamilyData) GetAvailableHours() (hours float64, err error) {
	/* Get all relevant bookings for the last week (or this week if !weekend) */
	today := time.Now()
	sow := time.Now()
	if fd.StartOfWeek.Before(today) {
		return math.Trunc((fd.HoursDone-fd.HoursGoal) * 1000) / 1000, nil // use this weeks data
	}
	// today is before the listed monday --> it is therefore a weekend
	now.Monday() // Make monday beginning of the week (now defaults to sunday for each closure)
	sow = now.New(today).BeginningOfWeek() // get last monday

	/* Get hours from bookings completed */
	bookings, err := fd.family.getFamilyBookings(sow, today) // get the bookings for the relevant week
	if err != nil {
		logger.Println("Error getting booking data (getavailablehours): ", err)
		return -1, err
	}
	for _, b := range bookings {
		hours += (b.End.Sub(b.Start).Hours() * float64(b.Modifier))
	}

	/* Account for donations */
	donations, err := fd.family.getShortDonations(sow, today)
	logger.Println(donations)
	if err != nil {
		return math.Trunc((hours - fd.HoursGoal) * 1000) / 1000, nil
	}
	hours += donations.netAmount(fd.family.ID)
	logger.Println(hours)

	return math.Trunc((hours - fd.HoursGoal) * 1000) / 1000, nil
}

func (fd *FamilyData) setHoursData(startOfWeek time.Time, today now.Now) error {
	/* Get all relevant bookings for the period */
	bookings, err := fd.family.getFamilyBookings(fd.StartOfPeriod, fd.EndOfPeriod)
	if err != nil {
		logger.Println("Error getting hours data: ", err)
		return err
	}
	/* Create the historical datasets + get this week's totals */
	for _, b := range bookings {
		duration := (b.End.Sub(b.Start).Hours() * float64(b.Modifier))
		// historical bookings are mapped by their parent userID
		if b.Start.Before(startOfWeek) {
			fd.History[b.UserID].addDurationPoint(DurationPoint{Y: duration, X: now.New(b.Start).BeginningOfDay()})
		} else if b.Start.Before(today.Time) && b.End.After(startOfWeek) {
			fd.HoursDone += duration
			fd.HoursBooked += duration // Even though they're done, theyre still booked 4 this week
		} else {
			fd.HoursBooked += duration // Time is after today, but before week end --> booked hours
		}
	}

	donations, err := fd.family.getShortDonations(startOfWeek, fd.EndOfPeriod)
	if err != nil {
		logger.Println(err)
		return err
	}

	fd.HoursDone += donations.netAmount(fd.family.ID)
	fd.HoursBooked = math.Trunc(fd.HoursBooked*1000) / 1000
	fd.HoursDone = math.Trunc(fd.HoursDone*1000) / 1000
	return err
}

// Add 0 values to correct charting behaviour inbetween gaps of data
func (fd *FamilyData) spanHistoryGaps() {
	for _, parent := range fd.family.Parents {
		fd.History[parent.UserID].configureAsHistoricalHours(parent.FirstName, colorful.FastHappyColor().Hex(), false, 0.00)
		data := fd.History[parent.UserID].Data
		var zeroData DurationPoints
		// fill in some 0 y-points to span gaps better
		for i, _ := range data {
			if i > 0 {
				dayDelta := data[i].X.YearDay() - data[i-1].X.YearDay()
				if dayDelta > 5 {
					// prepend a 0 value point with x = median day
					medianDay := data[i-1].X.AddDate(0, 0, dayDelta/2)
					medianPt := DurationPoint{X: medianDay, Y: 0}
					zeroData = append(zeroData, medianPt)
				}
			}
		}
		fd.History[parent.UserID].Data = append(fd.History[parent.UserID].Data, zeroData...)
		sort.Sort(fd.History[parent.UserID].Data)
	}
}

/*
 * Creates a FamilyData struct and fills it with data given a family, and the reference date: today
 * Get the hours data relative to the day passed as today.
 * The history will be tracked from the FIRST_MONTH and FIRST_DAY
 * of today's year.
 */
func (fd *FamilyData) init(fam *Family, today time.Time) error {
	fd.setHoursGoal(fam.Children)
	fd.family = fam // set fd to hold family given
	fd.History = make(map[int]*ChartDataSet)
	for _, prnt := range fam.Parents {
		fd.History[prnt.UserID] = new(ChartDataSet) // create datasets for each parent
	}

	/* Parse start/end for history retrieval */
	realToday := now.New(today)         // need this incase we alter today value to account for weekend viewing
	if today.Weekday() == time.Sunday { // weekend days must be shifted to monday
		today = today.AddDate(0, 0, 1) // move to monday so we reference next week
	} else if today.Weekday() == time.Saturday {
		today = today.AddDate(0, 0, 2) // move to monday
	}

	now.Monday()                                           // first day of week
	nowToday := now.New(today)                             // We use this to determine start of week, so it should be the adjusted today
	fd.StartOfWeek = nowToday.BeginningOfWeek()             // if weekend, this is next week's monday
	fd.EndOfPeriod = nowToday.EndOfWeek()                  // set end of period
	fd.StartOfPeriod = today.AddDate(0, -(PERIOD_LENGTH), 0) // Go back 4 months --> set startperiod
	// create the history
	err := fd.setHoursData(fd.StartOfWeek, *realToday)
	if err != nil {
		return err
	}
	return err
}

// Like get user bookings but for a family
func (f *Family) getFamilyBookings(start time.Time, end time.Time) ([]Booking, error) {
	/* Get all bookings in range start-now  (start > block_start & end > blocK_end) */
	q := `SELECT booking_id, block_id, family_id, user_id, block_start, block_end, room_id, modifier
			FROM booking NATURAL JOIN time_block WHERE (
					time_block.block_id = booking.block_id
					AND booking.family_id = $1
					AND time_block.block_start >= $2 AND time_block.block_start < $3
					AND time_block.block_end > $2 AND  time_block.block_end <= $3
			) ORDER BY block_start`

	var bookBlocks []Booking
	err := db.Select(&bookBlocks, q, f.ID, start, end)
	logger.Println("Selected blocks: ", bookBlocks)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Ensure locale is set
	for i, b := range bookBlocks {
		bookBlocks[i].Start = time.Date(b.Start.Year(), b.Start.Month(), b.Start.Day(),
			b.Start.Hour(), b.Start.Minute(), 0, 0, time.Local)
		bookBlocks[i].End = time.Date(b.End.Year(), b.End.Month(), b.End.Day(),
			b.End.Hour(), b.End.Minute(), 0, 0, time.Local)

	}
	return bookBlocks, nil
}

// Get donations between start and end two-tailed inclusive (where the family may be either the donor, or donee
func (f *Family) getShortDonations(start, end time.Time) (gifts ShortDonations, err error) {
	q := `SELECT * FROM donation 
			WHERE (donation.donor_id = $1 OR donation.donee_id = $1)
			AND (
					(donation.date_sent <= $3 AND donation.date_sent >= $2)
				OR 
					(donation.date_sent <= $2 AND donation.date_sent >= $3)
				)`

	err = db.Select(&gifts, q, f.ID, start, end)
	if err != nil {
		if err == sql.ErrNoRows {
			return gifts, nil
		} else {
			return nil, err
		}
	}
	return gifts, nil
}

func (donor *Family) GetID() int {
	return donor.ID
}

/*
 * Send a donation -- including saving it to db
 */
func (donor *Family) GiveCharity(donee *Family, amount float64) (*Donation, error) {
	d := new(Donation)
	d.Sender = donor
	d.Recipient = donee
	d.Amount = amount
	err := d.save()
	return d, err
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
		return nil, &ClientSafeError{Msg:"invalid token"}
	}

	q := `SELECT family_id, family_name, children 
			FROM users NATURAL JOIN family 
			WHERE users.username = $1
				AND family.family_id = users.family_id`

	fdata := new(Family)
	err := db.QueryRow(q, username).Scan(&fdata.ID, &fdata.Name, &fdata.Children)
	if err != nil {
		logger.Println(err)
		return nil, &ClientSafeError{Msg:"could not retrieve family information"}
	}

	var uids []int
	q = `SELECT user_id FROM users WHERE users.family_id = $1`
	err = db.Select(&uids, q, fdata.ID)
	if err != nil {
		return fdata, err
	}
	// fill parents slice
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

func GetFamilyByID(ID int) (f *Family, err error) {
	f = new(Family)
	q := `SELECT * FROM family WHERE family_id = $1`
	err = db.QueryRow(q, ID).Scan(&f.ID, &f.Name, &f.Children)
	if err != nil {
		logger.Println(err)
	}
	return f, err
}

/*
 * User sans password field
 */
type User struct {
	UserID    int
	Role      int
	Username  string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	FamilyID  int
	Bonus     float64
	BonusNote string
}

/*
 * Initializes reciever struct based on the given UID, a user from the db.
 */
func (u *User) init(UID int) error {
	q := `SELECT 	user_id, user_role, username, first_name, last_name, 
					email, phone_number, family_id, COALESCE(bonus_hours, 0), COALESCE(bonus_note, '')
			FROM users 
			WHERE user_id = $1`
	err := db.QueryRow(q, UID).Scan(&u.UserID, &u.Role, &u.Username, &u.FirstName, &u.LastName,
		&u.Email, &u.Phone, &u.FamilyID, &u.Bonus, &u.BonusNote)
	if err != nil {
		return err
	}
	return nil
}

// Get bookings for a user by ID
func getUserBookings(start time.Time, end time.Time, UID int) ([]Booking, error) {
	/* Get all bookings in range start-now  (start > block_start & end > blocK_end) */
	q := `SELECT booking_id, block_id, family_id, user_id, block_start, block_end, room_id, modifier
			FROM booking NATURAL JOIN time_block WHERE (
					time_block.block_id = booking.block_id
					AND booking.user_id = $1
					AND time_block.block_start >= $2 AND time_block.block_start < $3
					AND time_block.block_end > $2 AND  time_block.block_end <= $3
			) ORDER BY block_start`

	var bookBlocks []Booking
	err := db.Select(&bookBlocks, q, UID, start, end)
	logger.Println("Selected blocks: ", bookBlocks)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return bookBlocks, nil
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

func getUIDFromName(userName string) (uid int, err error) {
	q := `SELECT user_id FROM users WHERE username = $1`
	err = db.QueryRow(q, userName).Scan(&uid)
	if err != nil {
		return -1, err
	}
	return
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

func getTaughtRooms(r *http.Request) (roomNames []string, err error) {
	//get username
	sesh, _ := store.Get(r, "loginSession")
	uname, ok := sesh.Values["username"].(string)
	if !ok {
		return nil, &ClientSafeError{Msg: "Invalid User Token, please re-login", Code: http.StatusNetworkAuthenticationRequired}
	}
	//get teaches rooms
	roomNames, err = userNameTeachesRooms(uname)
	return
}

/*
 * Get the room a teacher teaches by username
 */
func userNameTeachesRooms(userName string) (roomNames []string, err error) {
	q := `SELECT room.room_name FROM room NATURAL JOIN users
				WHERE room.teacher_id = users.user_id
						AND users.username = $1`

	err = db.Select(&roomNames, q, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			err = &ClientSafeError{Msg: "Teacher not assigned a room", Code: http.StatusBadRequest}
		}
		logger.Println(err)
	}
	return
}
