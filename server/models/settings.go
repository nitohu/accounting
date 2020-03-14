package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/nitohu/err"
)

// Settings object
type Settings struct {
	Name         string
	Email        string
	SalaryDate    time.Time
	CalcInterval int64
	CalcUoM      string
	Currency     string

	lastUpdate time.Time

	SalaryDateForm string
}

// InitializeSettings creates an empty settings object and initializes it
func InitializeSettings(cr *sql.DB) (Settings, err.Error) {
	s := Settings{
		Name:          "",
		Email:         "",
		SalaryDate:     time.Now(),
		CalcInterval:  0,
		CalcUoM:       "",
		Currency:      "",
		lastUpdate:    time.Now(),
		SalaryDateForm: "",
	}

	e := s.Init(cr)

	if !e.Empty() {
		e.AddTraceback("e.InitializeSettings()", "Error initializing the settings")
		return s, e
	}

	return s, err.Error{}
}

// Init the settings
func (s *Settings) Init(cr *sql.DB) err.Error {
	query := "SELECT name,email,last_update,salary_date,calc_interval,calc_uom,currency FROM settings;"

	e := cr.QueryRow(query).Scan(
		&s.Name,
		&s.Email,
		&s.lastUpdate,
		&s.SalaryDate,
		&s.CalcInterval,
		&s.CalcUoM,
		&s.Currency,
	)
	if e != nil {
		var err err.Error
		err.Init("Settings.Init()", e.Error())
		return err
	}

	s.computeFields(cr)

	return err.Error{}
}

// Save the current Settings object to the database
// func (s *Settings) Save(cr *sql.DB, password string) error {
func (s *Settings) Save(cr *sql.DB) err.Error {
	query := "UPDATE settings SET name=$1,email=$2,last_update=$3,salary_date=$4,calc_interval=$5,calc_uom=$6,currency=$7;"

	_, e := cr.Exec(query,
		s.Name,
		s.Email,
		time.Now(),
		s.SalaryDate,
		s.CalcInterval,
		s.CalcUoM,
		s.Currency,
	)
	if e != nil {
		var err err.Error
		err.Init("Settings.Save()", e.Error())
		return err
	}

	s.computeFields(cr)

	return err.Error{}
}

// GetLastUpdate returns the value of the last_update field
func (s *Settings) GetLastUpdate() time.Time {
	return s.lastUpdate
}

// ShiftSalaryDate ..
func (s *Settings) ShiftSalaryDate(cr *sql.DB) err.Error {
	n := time.Now()
	shifted := false
	for n.After(s.SalaryDate) {
		shifted = true
		s.SalaryDate = s.SalaryDate.AddDate(0, 1, 0)
	}
	if shifted {
		if e := s.Save(cr); !e.Empty() {
			e.AddTraceback("Settings.shiftSalaryDate()", "Error while saving the settings to the database.")
			return e
		}
		log.Println("[INFO] Settings.ShiftSalaryDate(): Salary date shifted.")
	}
	return err.Error{}
}

func (s *Settings) computeFields(cr *sql.DB) {
	s.SalaryDateForm = s.SalaryDate.Format("Monday 02 January 2006")
}
