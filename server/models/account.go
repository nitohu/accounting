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
func (a *Account) Create(cr *sql.DB) error {

	if a.ID != 0 {
		return errors.New("This object already has an id")
	}

	var id int64

	query := "INSERT INTO accounts ( name, active, balance, balance_forecast, iban,"
	query += " bank_code, account_nr, bank_name, bank_type, create_date, last_update"
	query += ") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;"

	a.CreateDate = time.Now().Local()
	a.LastUpdate = time.Now().Local()

	err := cr.QueryRow(query,
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

	if err != nil {
		fmt.Println("Origin: Account.Create")
		return err
	}

	a.ID = id

	return nil
}

// Save 's the current values of the object to the database
func (a *Account) Save(cr *sql.DB) error {
	if a.ID == 0 {
		return errors.New("This account as now id, maybe create it first?")
	}

	query := "UPDATE accounts SET name=$2, active=$3, balance=$4, balance_forecast=$5, iban=$6,"
	query += " bank_code=$7, account_nr=$8, bank_name=$9, bank_type=$10, last_update=$11 WHERE id=$1"

	res, err := cr.Exec(query,
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
	if err != nil {
		return err
	}
	if _, err = res.RowsAffected(); err != nil {
		return err

	}
	if err = a.ComputeFields(cr); err != nil {
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

// ComputeFields computes the fields for this model
// Gets automatically called in Account.Save() and Account.FindByID()
func (a *Account) ComputeFields(cr *sql.DB) error {
	if a.ID <= 0 {
		return errors.New("This account has no ID")
	}

	query := "SELECT COUNT(*) FROM transactions WHERE account_id=$1;"

	err := cr.QueryRow(query, a.ID).Scan(
		&a.TransactionCount,
	)
	if err != nil {
		return err
	}

	return nil
}

// Book books a transaction in the account
// Also saves the new balance to the database
func (a *Account) Book(cr *sql.DB, t *Transaction, invert bool) error {

	amount := t.Amount

	if invert == true {
		amount = amount * -1
	}

	fmt.Printf("Book function of account; %s (%d)\n", a.Name, a.ID)
	fmt.Printf("Transaction; %s (%d) %f\n", t.Name, t.ID, amount)

	a.BalanceForecast += amount
	a.Balance += amount

	// currentTime := time.Now().Local()

	// TODO: Unnecessary for now, will be more important for later
	// features (forecasting and later booking)
	// if currentTime.After(t.TransactionDate) {
	// 	a.Balance += amount
	// }

	err := a.Save(cr)

	if err != nil {
		return err
	}

	return nil
}

// FindByID finds an account with it's id
func (a *Account) FindByID(cr *sql.DB, accountID int64) error {
	query := "SELECT * FROM accounts WHERE id=$1"

	err := cr.QueryRow(query, accountID).Scan(
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

	if err != nil {
		fmt.Println("Traceback: Account.FindById():", a.ID)
		return err
	}

	if err = a.ComputeFields(cr); err != nil {
		return err
	}

	return nil
}

// FindAccountByID is similar to FindByID but returns the account
func FindAccountByID(cr *sql.DB, accountID int64) (Account, error) {
	a := EmptyAccount()

	err := a.FindByID(cr, accountID)

	if err != nil {
		return a, err
	}

	return a, nil
}

// GetAllAccounts does that what you expect
func GetAllAccounts(cr *sql.DB) ([]Account, error) {
	var accounts []Account
	query := "SELECT id FROM accounts"

	idRows, err := cr.Query(query)

	if err != nil {
		return accounts, err
	}

	for idRows.Next() {
		var id int64
		err := idRows.Scan(&id)

		if err != nil {
			fmt.Printf("[WARN] %s GetAllAccounts():\n[INFO] Skipping Record\n%s", time.Now().Local(), err)
		} else {
			a := EmptyAccount()
			err = a.FindByID(cr, id)

			if err != nil {
				fmt.Printf("[WARN] %s GetAllAccounts():\n[INFO] Skipping Record\n%s", time.Now().Local(), err)
			} else {
				accounts = append(accounts, a)
			}
		}

	}

	return accounts, nil
}

// GetLimitAccounts returns a limited number of accounts
func GetLimitAccounts(cr *sql.DB, number int) ([]Account, error) {
	var result []Account

	if number <= 0 {
		err := "The number of the accounts must be bigger than 0."
		return result, errors.New(err)
	}

	query := "SELECT id FROM accounts LIMIT $1"

	res, err := cr.Query(query, number)

	if err != nil {
		fmt.Printf("[INFO] %s Account.GetLimitAccount(): Traceback: Error executing query.\n", time.Now())
		return nil, err
	}

	for res.Next() {
		var id int64
		if err = res.Scan(&id); err != nil {
			return nil, err
		}
		acc := EmptyAccount()
		if err = acc.FindByID(cr, id); err != nil {
			return nil, err
		}
		result = append(result, acc)
	}

	return result, nil
}
