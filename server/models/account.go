package models

import (
	"database/sql"
	"time"
)

type Account struct {
	ID              int
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

func (a *Account) FindByID(cr *sql.DB, accountId int) error {
	query := "SELECT * FROM accounts WHERE id=$1"

	err := cr.QueryRow(query, accountId).Scan(
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

func FindAccountById(cr *sql.DB, accountId int) Account {
	a := EmptyAccount()

	err := a.FindByID(cr, accountId)

	if err != nil {
		panic(err)
	}

	return a
}
