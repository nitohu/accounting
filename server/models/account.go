package models

import "time"

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
	UserID          User
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
	}

	return a
}
