package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

/*
Query:

SELECT
    t.id,
    t.name,
    t.active,
    t.transaction_date,
    t.create_date,
    t.last_update,
    t.amount,
    t.transaction_type,
    u.name,
    a.name
FROM transactions AS t
JOIN users AS u ON t.user_id=u.id
JOIN accounts AS a ON t.account_id=a.id
WHERE t.user_id= *id*;
*/

// Transaction model
type Transaction struct {
	ID              int64
	Name            string
	Active          bool
	TransactionDate time.Time
	CreateDate      time.Time
	LastUpdate      time.Time
	Amount          float64
	FromAccount     int64
	ToAccount       int64
	TransactionType string
	UserID          int
}

// EmptyTransaction ..
func EmptyTransaction() Transaction {
	t := Transaction{
		ID:              0,
		Name:            "",
		Active:          false,
		TransactionDate: time.Now().Local(),
		CreateDate:      time.Now().Local(),
		LastUpdate:      time.Now().Local(),
		Amount:          0.0,
		FromAccount:     0,
		ToAccount:       0,
		TransactionType: "",
		UserID:          0,
	}

	return t
}

// Create 's a transaction with the current values of the object
func (t *Transaction) Create(cr *sql.DB) error {

	if t.ID != 0 {
		return errors.New("This object already has a user id")
	} else if t.UserID == 0 {
		return errors.New("No user is linked to the transaction")
	} else if t.Amount == 0.0 {
		return errors.New("The Amount of this transaction is 0")
	}

	query := "INSERT INTO transactions ( name, active, transaction_date, last_update, create_date, amount,"
	query += " account_id, to_account, transaction_type, user_id"
	query += ") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"

	t.CreateDate = time.Now().Local()
	t.LastUpdate = time.Now().Local()

	res, err := cr.Exec(query,
		t.Name,
		t.Active,
		t.TransactionDate,
		t.LastUpdate,
		t.CreateDate,
		t.Amount,
		t.FromAccount,
		t.ToAccount,
		t.TransactionType,
		t.UserID,
	)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()

	if rowCount, err := res.RowsAffected(); err != nil || rowCount < 1 {
		return errors.New("No rows affected. ID: " + fmt.Sprintf("%s", id))
	}

	t.ID = id

	return nil
}

// Save 's the current values of the object to the database
func (t *Transaction) Save(cr *sql.DB) error {

	if t.ID == 0 {
		return errors.New("This account as now id, maybe create it first?")
	} else if t.Amount == 0.0 {
		return errors.New("The Amount of the transaction with the id " + fmt.Sprintf("%s", t.ID) + " is 0")
	}

	query := "UPDATE accounts SET name=$2, active=$3, transaction_date=$4, last_update=$5, amount=$6, account_id=$7,"
	query += " to_account=$8, transaction_type=$9 WHERE id=$1"

	t.TransactionDate = time.Now().Local()

	res, err := cr.Exec(query,
		t.ID,
		t.Name,
		t.Active,
		t.TransactionDate,
		t.LastUpdate,
		t.Amount,
		t.FromAccount,
		t.ToAccount,
		t.TransactionType,
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

// Delete 's the transtaction
func (t *Transaction) Delete(cr *sql.DB) error {
	if t.ID == 0 {
		return errors.New("The account you want to delete does not have an id")
	}

	query := "DELETE FROM accounts WHERE id=$1"

	_, err := cr.Exec(query, t.ID)

	if err != nil {
		return err
	}

	return nil
}

// FindByID finds a transaction with it's id
func (t *Transaction) FindByID(cr *sql.DB, transactionID int) error {
	query := "SELECT * FROM accounts WHERE id=$1"

	err := cr.QueryRow(query, transactionID).Scan(
		&t.ID,
		&t.Name,
		&t.Active,
		&t.TransactionDate,
		&t.LastUpdate,
		&t.CreateDate,
		&t.Amount,
		&t.FromAccount,
		&t.ToAccount,
		&t.TransactionType,
		&t.UserID,
	)

	if err != nil {
		return err
	}

	return nil
}

// FindTransactionByID is similar to FindByID but returns the transaction
func FindTransactionByID(cr *sql.DB, transactionID int) Transaction {
	t := EmptyTransaction()

	err := t.FindByID(cr, transactionID)

	if err != nil {
		panic(err)
	}

	return t
}
