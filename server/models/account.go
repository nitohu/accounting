package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/nitohu/err"
)

// Account object
type Account struct {
	ID              int64
	Name            string
	Active          bool
	Balance         float64
	BalanceForecast float64
	Iban            string
	BankCode        string
	AccountNr       string
	BankName        string
	BankType        string
	CreateDate      time.Time
	LastUpdate      time.Time

	// Computed Fields
	TransactionCount int64
}

// EmptyAccount ...
func EmptyAccount() Account {
	a := Account{
		ID:               0,
		Name:             "",
		Active:           true,
		Balance:          0.0,
		BalanceForecast:  0.0,
		Iban:             "",
		BankCode:         "",
		AccountNr:        "",
		BankName:         "",
		BankType:         "",
		CreateDate:       time.Now().Local(),
		LastUpdate:       time.Now().Local(),
		TransactionCount: 0,
	}

	return a
}

// Create 's an account with the current values of the object
func (a *Account) Create(cr *sql.DB) err.Error {

	if a.ID != 0 {
		var err err.Error
		err.Init("Account.Create()", "This object already has an id")
		return err
	}

	var id int64

	query := "INSERT INTO accounts ( name, active, balance, balance_forecast, iban,"
	query += " bank_code, account_nr, bank_name, bank_type, create_date, last_update"
	query += ") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;"

	a.CreateDate = time.Now().Local()
	a.LastUpdate = time.Now().Local()

	e := cr.QueryRow(query,
		a.Name,
		a.Active,
		a.Balance,
		a.BalanceForecast,
		a.Iban,
		a.BankCode,
		a.AccountNr,
		a.BankName,
		a.BankType,
		a.CreateDate,
		a.LastUpdate,
	).Scan(&id)

	if e != nil {
		var err err.Error
		err.Init("Account.Create()", e.Error())
		return err
	}

	a.ID = id

	return err.Error{}
}

// Save 's the current values of the object to the database
func (a *Account) Save(cr *sql.DB) err.Error {
	if a.ID == 0 {
		var err err.Error
		err.Init("Account.Save()", "This account as no ID, maybe create it first?")
		return err
	}

	query := "UPDATE accounts SET name=$2, active=$3, balance=$4, balance_forecast=$5, iban=$6,"
	query += " bank_code=$7, account_nr=$8, bank_name=$9, bank_type=$10, last_update=$11 WHERE id=$1"

	res, e := cr.Exec(query,
		a.ID,
		a.Name,
		a.Active,
		a.Balance,
		a.BalanceForecast,
		a.Iban,
		a.BankCode,
		a.AccountNr,
		a.BankName,
		a.BankType,
		time.Now().Local(),
	)
	if e != nil {
		var err err.Error
		err.Init("Account.Save()", e.Error())
		return err
	}
	if _, e = res.RowsAffected(); e != nil {
		var err err.Error
		err.Init("Account.Save()", e.Error())
		return err

	}
	a.computeFields(cr)

	return err.Error{}
}

// Delete 's the account
func (a *Account) Delete(cr *sql.DB) err.Error {
	if a.ID == 0 {
		var err err.Error
		err.Init("Account.Delete()", "The account you want to delete does not have an id")
		return err
	}

	query := "DELETE FROM accounts WHERE id=$1"

	_, e := cr.Exec(query, a.ID)

	if e != nil {
		var err err.Error
		err.Init("Account.Delete()", e.Error())
		return err
	}

	return err.Error{}
}

// ComputeFields computes the fields for this model
// Gets automatically called in Account.Save() and Account.FindByID()
func (a *Account) computeFields(cr *sql.DB) {
	query := "SELECT COUNT(*) FROM transactions WHERE account_id=$1;"

	e := cr.QueryRow(query, a.ID).Scan(
		&a.TransactionCount,
	)
	if e != nil {
		var err err.Error
		err.Init("Account.ComputeFields()", "Error getting the number of transactions for account "+fmt.Sprintf("%d", a.ID))
		log.Println("[WARN]", err)
	}
}

// Book books a transaction in the account
// Also saves the new balance to the database
func (a *Account) Book(cr *sql.DB, t *Transaction, invert bool) err.Error {
	amount := t.Amount
	if invert == true {
		amount = amount * -1
	}

	a.BalanceForecast += amount
	a.Balance += amount

	// currentTime := time.Now().Local()

	// TODO: Unnecessary for now, will be more important for later
	// features (forecasting and later booking)
	// if currentTime.After(t.TransactionDate) {
	// 	a.Balance += amount
	// }

	e := a.Save(cr)

	if e.Empty() == false {
		var err err.Error
		err.AddTraceback("Account.Book()", "Error while saving the Account "+a.Name)
		return err
	}

	return err.Error{}
}

// FindByID finds an account with it's id
func (a *Account) FindByID(cr *sql.DB, accountID int64) err.Error {
	query := "SELECT * FROM accounts WHERE id=$1"

	e := cr.QueryRow(query, accountID).Scan(
		&a.ID,
		&a.Name,
		&a.Active,
		&a.Balance,
		&a.BalanceForecast,
		&a.Iban,
		&a.BankCode,
		&a.AccountNr,
		&a.BankName,
		&a.BankType,
		&a.CreateDate,
		&a.LastUpdate,
	)

	if e != nil {
		var err err.Error
		err.Init("Account.FindById():", e.Error())
		return err
	}

	a.computeFields(cr)

	return err.Error{}
}

// FindAccountByID is similar to FindByID but returns the account
func FindAccountByID(cr *sql.DB, accountID int64) (Account, err.Error) {
	a := EmptyAccount()

	e := a.FindByID(cr, accountID)

	if !e.Empty() {
		msg := fmt.Sprintf("Error while getting account with ID %d from the database.", a.ID)
		e.AddTraceback("FindAccountByID()", msg)
		return a, e
	}

	return a, err.Error{}
}

// GetAllAccounts does that what you expect
func GetAllAccounts(cr *sql.DB) ([]Account, err.Error) {
	var accounts []Account
	query := "SELECT id FROM accounts"

	idRows, e := cr.Query(query)

	if e != nil {
		var err err.Error
		err.Init("GetAllAccounts()", e.Error())
		return accounts, err
	}

	for idRows.Next() {
		var id int64
		e = idRows.Scan(&id)

		if e != nil {
			log.Println("[INFO] GetAllAccounts(): Skipping record.")
			log.Println("[WARN] GetAllAccounts(): ", e)
		} else {
			a := EmptyAccount()
			err := a.FindByID(cr, id)

			if !err.Empty() {
				err.AddTraceback("GetAllAccounts()", "Error while finding account: "+fmt.Sprintf("%d", id))
				log.Println("[WARN]", err)
			} else {
				accounts = append(accounts, a)
			}
		}

	}

	return accounts, err.Error{}
}

// GetLimitAccounts returns a limited number of accounts
func GetLimitAccounts(cr *sql.DB, number int) ([]Account, err.Error) {
	var result []Account

	if number <= 0 {
		var err err.Error
		err.Init("GetLimitAccounts()", "The number of the accounts must be bigger than 0.")
		return result, err
	}

	query := "SELECT id FROM accounts LIMIT $1"

	res, e := cr.Query(query, number)

	if e != nil {
		var err err.Error
		err.Init("Account.GetLimitAccount()", e.Error())
		return nil, err
	}

	for res.Next() {
		var id int64
		if e = res.Scan(&id); e != nil {
			var err err.Error
			err.Init("Account.GetLimitAccount()", e.Error())
			return nil, err
		}
		acc := EmptyAccount()
		if e = acc.FindByID(cr, id); e != nil {
			var err err.Error
			err.Init("Account.GetLimitAccount()", e.Error())
			return nil, err
		}
		result = append(result, acc)
	}

	return result, err.Error{}
}
