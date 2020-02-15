package models

import (
	"database/sql"
	"time"
)

// Settings object
type Settings struct {
	Name         string
	Email        string
	StartDate    time.Time
	CalcInterval int64
	CalcUoM      string
	Currency     string

	lastUpdate time.Time

	StartDateForm string
}

// InitializeSettings creates an empty settings object and initializes it
func InitializeSettings(cr *sql.DB) (Settings, error) {
	s := Settings{
		Name:          "",
		Email:         "",
		StartDate:     time.Now(),
		CalcInterval:  0,
		CalcUoM:       "",
		Currency:      "",
		lastUpdate:    time.Now(),
		StartDateForm: "",
	}

	err := s.Init(cr)

	if err != nil {
		return s, err
	}

	return s, nil
}

// Init the settings
func (s *Settings) Init(cr *sql.DB) error {
	query := "SELECT name,email,last_update,start_date,calc_interval,calc_uom,currency FROM settings;"

	err := cr.QueryRow(query).Scan(
		&s.Name,
		&s.Email,
		&s.lastUpdate,
		&s.StartDate,
		&s.CalcInterval,
		&s.CalcUoM,
		&s.Currency,
	)

	s.computeFields(cr)

	return err
}

// Save the current Settings object to the database
// func (s *Settings) Save(cr *sql.DB, password string) error {
func (s *Settings) Save(cr *sql.DB) error {
	query := "UPDATE settings SET name=$1,email=$2,last_update=$3,start_date=$4,calc_interval=$5,calc_uom=$6,currency=$7;"

	_, err := cr.Exec(query,
		s.Name,
		s.Email,
		time.Now(),
		s.StartDate,
		s.CalcInterval,
		s.CalcUoM,
		s.Currency,
	)

	s.computeFields(cr)

	return err
}

// GetLastUpdate returns the value of the last_update field
func (s *Settings) GetLastUpdate() time.Time {
	return s.lastUpdate
}

func (s *Settings) computeFields(cr *sql.DB) {
	s.StartDateForm = s.StartDate.Format("Monday 02 January 2006")
}
