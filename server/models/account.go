package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Account object
type Account struct {
	ID              int64
	Name            string
	Active          bool
	Balance         float64
	BalanceForecast float64
	Iban            string
	Holder          string
	BankCode        string
	AccountNr       string
	BankName        string
	BankType        string
	CreateDate      time.Time
	LastUpdate      time.Time
	UserID          int
}

// EmptyAccount ...
func EmptyAccount() Account {
	a := Account{
		ID:              0,
		Name:            "",
		Active:          false,
		Balance:         0.0,
		BalanceForecast: 0.0,
		Iban:            "",
		Holder:          "",
		BankCode:        "",
		AccountNr:       "",
		BankName:        "",
		BankType:        "",
		CreateDate:      time.Now().Local(),
		LastUpdate:      time.Now().Local(),
		UserID:          0,
	}

	return a
}

// Create 's an account with the current values of the object
func (a *Account) Create(cr *sql.DB) error {

	if a.ID != 0 {
		return errors.New("This object already has a user id")
	} else if a.UserID == 0 {
		return errors.New("No user is linked to the account")
	}

	query := "INSERT INTO accounts ( name, active, balance, balance_forecast, iban, account_holder,"
	query += " bank_code, account_nr, bank_name, bank_type, create_date, last_update, user_id"
	query += " ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);"

	a.CreateDate = time.Now().Local()
	a.LastUpdate = time.Now().Local()

	res, err := cr.Exec(query,
		a.Name,
		a.Active,
		a.Balance,
		a.BalanceForecast,
		a.Iban,
		a.Holder,
		a.BankCode,
		a.AccountNr,
		a.BankName,
		a.BankType,
		a.CreateDate,
		a.LastUpdate,
		a.UserID,
	)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()

	if rowCount, err := res.RowsAffected(); err != nil || rowCount < 1 {
		return errors.New("No rows affected. ID: " + fmt.Sprintf("%d", id))
	}

	a.ID = id

	return nil
}

// Save 's the current values of the object to the database
func (a *Account) Save(cr *sql.DB) error {

	if a.ID == 0 {
		return errors.New("This account as now id, maybe create it first?")
	}

	query := "UPDATE accounts SET name=$2, active=$3, balance=$4, balance_forecast=$5, iban=$6, account_holder=$7,"
	query += " bank_code=$8, account_nr=$9, bank_name=$10, bank_type=$11, last_update=$12 WHERE id=$1"

	res, err := cr.Exec(query,
		a.ID,
		a.Name,
		a.Active,
		a.Balance,
		a.BalanceForecast,
		a.Iban,
		a.Holder,
		a.BankCode,
		a.AccountNr,
		a.BankName,
		a.BankType,
		time.Now().Local(),
	)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()

	if err != nil {
		return err
	}

	return nil
}

// Delete 's the account
func (a *Account) Delete(cr *sql.DB) error {
	if a.ID == 0 {
		return errors.New("The account you want to delete does not have an id")
	}

	query := "DELETE FROM accounts WHERE id=$1"

	_, err := cr.Exec(query, a.ID)

	if err != nil {
		return err
	}

	return nil
}

// FindByID finds an account with it's id
func (a *Account) FindByID(cr *sql.DB, accountID int) error {
	query := "SELECT * FROM accounts WHERE id=$1"

	err := cr.QueryRow(query, accountID).Scan(
		&a.ID,
		&a.Name,
		&a.Active,
		&a.Balance,
		&a.BalanceForecast,
		&a.Iban,
		&a.Holder,
		&a.BankCode,
		&a.AccountNr,
		&a.BankName,
		&a.BankType,
		&a.CreateDate,
		&a.LastUpdate,
		&a.UserID,
	)

	if err != nil {
		return err
	}

	return nil
}

// FindAccountByID is similar to FindByID but returns the account
func FindAccountByID(cr *sql.DB, accountID int) Account {
	a := EmptyAccount()

	err := a.FindByID(cr, accountID)

	if err != nil {
		panic(err)
	}

	return a
}
