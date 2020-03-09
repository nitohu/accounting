package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nitohu/err"
)

/*
Statistic object
Visualisation Types:
 * number
 * minibar -> mini bar chart, e.g.: Dashboard, latest transactions
*/
type Statistic struct {
	ID            int64
	Active        bool
	Name          string
	ComputeQuery  string
	LastUpdate    time.Time
	Description   string
	Keys          string
	Value         string
	ExecutionDate time.Time
	Visualisation string
}

// EmptyStatistic returns an empty statistic
func EmptyStatistic() Statistic {
	stat := Statistic{
		ID:            0,
		Active:        false,
		Name:          "",
		ComputeQuery:  "",
		LastUpdate:    time.Now(),
		ExecutionDate: time.Now(),
		Description:   "",
		Keys:          "",
		Value:         "",
		Visualisation: "",
	}

	return stat
}

// Create the current object in the database
func (s *Statistic) Create(cr *sql.DB) err.Error {
	if s.ID > 0 {
		var err err.Error
		err.Init("Statistic.Create()", "The Statistic "+s.Name+" already has an ID. Maybe try saving it?")
		return err
	} else if s.Name == "" {
		var err err.Error
		err.Init("Statistic.Create()", "The Statistic does not have a Name.")
		return err
	} else if s.ComputeQuery == "" {
		var err err.Error
		err.Init("Statistic.Create()", "The Statistic does not have a query.")
		return err
	}

	query := "INSERT INTO statistics (active,name,compute_query,last_update,create_date,description,visualisation,keys,"
	query += "value,execution_date) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);"

	res, e := cr.Exec(query,
		s.Active,
		s.Name,
		s.ComputeQuery,
		time.Now(),
		time.Now(),
		s.Description,
		s.Visualisation,
		s.Keys,
		s.Value,
		s.ExecutionDate,
	)

	if e != nil {
		var err err.Error
		err.Init("Statistic.Create()", e.Error())
		return err
	}

	if s.ID, e = res.LastInsertId(); e != nil {
		var err err.Error
		err.Init("Statistic.Create()", e.Error())
		return err
	}

	if err := s.Compute(cr); !err.Empty() {
		err.AddTraceback("Statistic.Create()", "Error computing the value of "+s.Name)
		return err
	}

	return err.Error{}
}

// Save the current object to the database
func (s *Statistic) Save(cr *sql.DB) err.Error {
	if s.ID <= 0 {
		var err err.Error
		err.Init("Statistic.Save()", "The Statistic "+s.Name+" does not have an ID. Maybe try creating it first?")
		return err
	} else if s.Name == "" {
		var err err.Error
		err.Init("Statistic.Save()", "The Statistic does not have a Name.")
		return err
	} else if s.ComputeQuery == "" {
		var err err.Error
		err.Init("Statistic.Save()", "The Statistic does not have a query.")
		return err
	}

	query := "UPDATE statistics SET name=$1,active=$2,compute_query=$3,last_update=$4,description=$5,visualisation=$6,keys=$7,"
	query += "value=$8,execution_date=$9 WHERE id=$10;"

	s.LastUpdate = time.Now()

	res, e := cr.Exec(query,
		s.Name,
		s.Active,
		s.ComputeQuery,
		s.LastUpdate,
		s.Description,
		s.Visualisation,
		s.Keys,
		s.Value,
		s.ExecutionDate,
		s.ID,
	)

	if e != nil {
		var err err.Error
		err.Init("Statistic.Save()", e.Error())
		return err
	}

	rows, e := res.RowsAffected()
	if e != nil {
		var err err.Error
		err.Init("Statistic.Save()", "Error while fetching the affected rows.")
		return err
	}

	if rows == 0 {
		var err err.Error
		err.Init("Statistic.Save()", "No rows where affected while saving.")
		return err
	}

	if err := s.Compute(cr); !err.Empty() {
		err.AddTraceback("Statistic.Save()", "Error computing the value of "+s.Name)
		return err
	}

	return err.Error{}
}

// Compute the value with the ComputeQuery
func (s *Statistic) Compute(cr *sql.DB) err.Error {
	// Make sure the salary_date (still start_date, will be changed soon) is always in the future
	settings, e := InitializeSettings(cr)
	if !e.Empty() {
		e.AddTraceback("Statistic.Compute()", "There was an error initializing the settings.")
		return e
	}
	if e := settings.ShiftSalaryDate(cr); !e.Empty() {
		e.AddTraceback("Statistic.Compute()", "Error while shifting the salary date.")
		return e
	}

	if error := cr.QueryRow(s.ComputeQuery).Scan(&s.Value); error != nil {
		var err err.Error
		err.Init("Statistic.Compute()", error.Error())
		return err
	}
	return err.Error{}
}

// FindByID finds a statistic by it's ID and sets it's value to the current object
func (s *Statistic) FindByID(cr *sql.DB, id int64) err.Error {
	if id <= 0 {
		var err err.Error
		err.Init("Statistic.FindByID()", "ID must be greater than 0")
		return err
	}

	query := "SELECT id,active,name,compute_query,last_update,execution_date,description,visualisation,keys,"
	query += "value FROM statistics WHERE id=$1"

	e := cr.QueryRow(query, id).Scan(
		&s.ID,
		&s.Active,
		&s.Name,
		&s.ComputeQuery,
		&s.LastUpdate,
		&s.ExecutionDate,
		&s.Description,
		&s.Visualisation,
		&s.Keys,
		&s.Value,
	)

	if e != nil {
		var err err.Error
		err.Init("Statistic.FindByID()", "Error fetching data from database")
		return err
	}

	if err := s.Compute(cr); !err.Empty() {
		err.AddTraceback("Statistic.FindByID()", "Error computing the value of "+s.Name)
		return err
	}

	return err.Error{}
}

// GetAllStatistics returns all statistics from the database
func GetAllStatistics(cr *sql.DB) ([]Statistic, err.Error) {
	query := "SELECT id FROM statistics"
	var stats []Statistic

	rows, e := cr.Query(query)
	if e != nil {
		var err err.Error
		err.Init("GetAllStatistics():", e.Error())
		return nil, err
	}

	for rows.Next() {
		var id int64
		s := EmptyStatistic()
		if e = rows.Scan(&id); e != nil {
			var err err.Error
			err.Init("GetAllStatistics()", e.Error())
			return stats, err
		}
		if err := s.FindByID(cr, id); !err.Empty() {
			err.AddTraceback("GetAllStatistics()", "Error while getting statistic: "+fmt.Sprintf("%d", id))
			return stats, err
		}
		stats = append(stats, s)
	}

	return stats, err.Error{}
}
