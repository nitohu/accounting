package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/nitohu/err"
)

// Transaction model
// TODO: Implement Forecasted and Booked in database
type Transaction struct {
	// Database fields
	ID              int64
	Name            string
	Description     string
	Active          bool
	TransactionDate time.Time
	CreateDate      time.Time
	LastUpdate      time.Time
	Amount          float64
	FromAccount     int64
	ToAccount       int64
	TransactionType string
	CategoryID      int64

	// Computed fields
	FromAccountName    string
	ToAccountName      string
	TransactionDateStr string
	Category           Category
}

// EmptyTransaction ..
func EmptyTransaction() Transaction {
	t := Transaction{
		ID:              0,
		Name:            "",
		Description:     "",
		Active:          false,
		TransactionDate: time.Now().Local(),
		CreateDate:      time.Now().Local(),
		LastUpdate:      time.Now().Local(),
		Amount:          0.0,
		FromAccount:     0,
		ToAccount:       0,
		TransactionType: "",
		CategoryID:      0,
	}

	return t
}

func bookIntoAccount(cr *sql.DB, id int64, t *Transaction, invert bool) err.Error {
	acc, e := FindAccountByID(cr, id)

	if !e.Empty() {
		e.AddTraceback("bookIntoAccount()", "Error while finding account")
		return e
	}

	e = acc.Book(cr, t, invert)

	if !e.Empty() {
		e.AddTraceback("bookIntoAccount()", "Error while booking into account")
		return e
	}

	return err.Error{}
}

// Create 's a transaction with the current values of the object
func (t *Transaction) Create(cr *sql.DB) err.Error {
	// Requirements for creating a transaction
	if t.ID != 0 {
		var err err.Error
		err.Init("Transaction.Create()", "This object already has an id")
		return err
	} else if t.Amount == 0.0 {
		var err err.Error
		err.Init("Transaction.Create()", "The Amount of this transaction is 0")
		return err
	}

	var id int64
	var categID interface{}

	// Initializing variables
	query := "INSERT INTO transactions ( name, active, transaction_date, last_update, create_date, amount,"
	query += " account_id, to_account, transaction_type, description, category_id"
	query += ") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;"

	t.CreateDate = time.Now().Local()
	t.LastUpdate = time.Now().Local()

	categID = t.CategoryID

	if t.CategoryID == 0 {
		categID = nil
	}

	var e error

	// 4 cases of executing the query
	if t.ToAccount == 0 && t.FromAccount > 0 {
		e = cr.QueryRow(query,
			t.Name,
			t.Active,
			t.TransactionDate,
			t.LastUpdate,
			t.CreateDate,
			t.Amount,
			t.FromAccount,
			nil,
			t.TransactionType,
			t.Description,
			categID,
		).Scan(&id)
	} else if t.ToAccount > 0 && t.FromAccount > 0 {
		e = cr.QueryRow(query,
			t.Name,
			t.Active,
			t.TransactionDate,
			t.LastUpdate,
			t.CreateDate,
			t.Amount,
			t.FromAccount,
			t.ToAccount,
			t.TransactionType,
			t.Description,
			categID,
		).Scan(&id)
		// TODO: This should throw an error
	} else if t.ToAccount == 0 && t.FromAccount == 0 {
		e = cr.QueryRow(query,
			t.Name,
			t.Active,
			t.TransactionDate,
			t.LastUpdate,
			t.CreateDate,
			t.Amount,
			nil,
			nil,
			t.TransactionType,
			t.Description,
			categID,
		).Scan(&id)
	} else if t.ToAccount > 0 && t.FromAccount == 0 {
		e = cr.QueryRow(query,
			t.Name,
			t.Active,
			t.TransactionDate,
			t.LastUpdate,
			t.CreateDate,
			t.Amount,
			nil,
			t.ToAccount,
			t.TransactionType,
			t.Description,
			categID,
		).Scan(&id)
	}

	if e != nil {
		var err err.Error
		err.Init("Transaction.Create()", e.Error())
		return err
	}

	// Writing id to object
	t.ID = id

	// Book transaction into FromAccount if it's given
	if t.FromAccount > 0 {
		err := bookIntoAccount(cr, t.FromAccount, t, true)

		if !err.Empty() {
			err.AddTraceback("Transaction.Create()", "Error while booking into FromAccount")
			return err
		}
	}

	// Book the transaction into ToAccount if it's given
	if t.ToAccount > 0 {
		err := bookIntoAccount(cr, t.ToAccount, t, false)

		if !err.Empty() {
			err.AddTraceback("Transaction.Create()", "Error while booking into ToAccount")
			return err
		}
	}

	return err.Error{}
}

// Save 's the current values of the object to the database
func (t *Transaction) Save(cr *sql.DB) err.Error {
	if t.ID == 0 {
		var err err.Error
		err.Init("Transaction.Save()", "This transaction as no ID, maybe create it first?")
		return err
	} else if t.Amount == 0.0 {
		var err err.Error
		err.Init("Transaction.Save()", "The Amount of the transaction with the id "+fmt.Sprintf("%d", t.ID)+" is 0")
		return err
	}

	// Get old data
	var oldAmount float64
	var accountID, toAccountID interface{}
	var TransactionDateStr time.Time

	query := "SELECT amount, transaction_date, account_id, to_account FROM transactions WHERE id=$1"

	row := cr.QueryRow(query, t.ID)

	e := row.Scan(&oldAmount, &TransactionDateStr, &accountID, &toAccountID)

	if e != nil {
		var err err.Error
		err.Init("Transaction.Save()", e.Error())
		return err
	}

	// Write values to database
	query = "UPDATE transactions SET name=$2, active=$3, transaction_date=$4, last_update=$5, amount=$6, account_id=$7,"
	query += "to_account=$8, transaction_type=$9, description=$10, category_id=$11 WHERE id=$1"

	var categID interface{}

	categID = t.CategoryID

	if t.CategoryID == 0 {
		categID = nil
	}

	if t.ToAccount == 0 && t.FromAccount > 0 {
		_, e = cr.Exec(query,
			t.ID,
			t.Name,
			t.Active,
			t.TransactionDate,
			t.LastUpdate,
			t.Amount,
			t.FromAccount,
			nil,
			t.TransactionType,
			t.Description,
			categID,
		)
	} else if t.ToAccount > 0 && t.FromAccount > 0 {
		_, e = cr.Exec(query,
			t.ID,
			t.Name,
			t.Active,
			t.TransactionDate,
			t.LastUpdate,
			t.Amount,
			t.FromAccount,
			t.ToAccount,
			t.TransactionType,
			t.Description,
			categID,
		)
	} else if t.ToAccount == 0 && t.FromAccount == 0 {
		// TODO: This case shouldn't be allowed
		_, e = cr.Exec(query,
			t.ID,
			t.Name,
			t.Active,
			t.TransactionDate,
			t.LastUpdate,
			t.Amount,
			nil,
			nil,
			t.TransactionType,
			t.Description,
			categID,
		)
	} else if t.ToAccount > 0 && t.FromAccount == 0 {
		_, e = cr.Exec(query,
			t.ID,
			t.Name,
			t.Active,
			t.TransactionDate,
			t.LastUpdate,
			t.Amount,
			nil,
			t.ToAccount,
			t.TransactionType,
			t.Description,
			categID,
		)
	}

	if e != nil {
		var err err.Error
		err.Init("Transaction.Save()", e.Error())
		return err
	}

	/*
		Change balance on accounts
	*/

	temp := EmptyTransaction()

	temp.Name = t.Name
	temp.Amount = oldAmount

	originChanged := false
	destinationChanged := false

	// The origin account changed
	if accountID != t.FromAccount {
		// Remove booking from old account
		// Amount needs to be inverted
		// Make sure the old amount gets removed from the old account
		if accountID != nil {
			err := bookIntoAccount(cr, accountID.(int64), &temp, false)

			if !err.Empty() {
				err.AddTraceback("Transaction.Save()", "Error while redoing booking from old origin account: "+accountID.(string))
				return err
			}
		}

		// Book the amount into the new account
		err := bookIntoAccount(cr, t.FromAccount, t, true)
		originChanged = true

		if !err.Empty() {
			err.AddTraceback("Transaction.Save()", "Error while booking to new origin account: "+fmt.Sprintf("%d", t.FromAccount))
			return err
		}
	}

	// The destination account changed
	if toAccountID != t.ToAccount {
		// Remove booking from old destination account
		// Make sure the old amount is used
		destinationChanged = true
		if toAccountID != nil {
			err := bookIntoAccount(cr, toAccountID.(int64), &temp, true)

			if !err.Empty() {
				err.AddTraceback("Transaction.Save()", "Error while removing transaction from the old receiving account: "+toAccountID.(string))
				return err
			}
		}

		// Book into new account
		err := bookIntoAccount(cr, t.ToAccount, t, false)

		if !err.Empty() {
			err.AddTraceback("Transaction.Save()", "Error while booking transaction into new receiving account: "+fmt.Sprintf("%d", t.ToAccount))
			return err
		}
	}

	// The amount of the transaction changed
	if t.Amount != oldAmount {
		diff := t.Amount - oldAmount

		temp.Amount = diff
		// The origin did not change
		// So book the difference into the origin account
		// If the origin has changed we've already booked into the accounts
		if !originChanged {
			// Calculate difference and book the difference into the origin account
			diff := t.Amount - oldAmount

			temp.Amount = diff

			err := bookIntoAccount(cr, t.FromAccount, &temp, true)

			if !err.Empty() {
				err.AddTraceback("Transaction.Save()", "Error while booking difference into origin account: "+fmt.Sprintf("%d", t.FromAccount))
				return err
			}
		}

		// The destination did not change
		if !destinationChanged {

			err := bookIntoAccount(cr, t.ToAccount, &temp, false)

			if !err.Empty() {
				err.AddTraceback("Transaction.Save()", "Error while booking difference into destination account: "+fmt.Sprintf("%d", t.ToAccount))
				return err
			}
		}
	}

	return err.Error{}
}

// Delete 's the transtaction
func (t *Transaction) Delete(cr *sql.DB) err.Error {
	if t.ID == 0 {
		var err err.Error
		err.Init("Transaction.Delete()", "The transaction you want to delete does not have an id")
		return err
	}

	if t.FromAccount > 0 {
		err := bookIntoAccount(cr, t.FromAccount, t, false)

		if !err.Empty() {
			err.AddTraceback("Transaction.Delete()", "Redo booking from origin account: "+fmt.Sprintf("%d", t.FromAccount))
			return err
		}
	}
	if t.ToAccount > 0 {
		err := bookIntoAccount(cr, t.ToAccount, t, true)

		if !err.Empty() {
			// TODO: If FromAccount was bigger than 0, it's already booked at this time
			// Make sure the booking is reverted again.
			err.AddTraceback("Transaction.Delete()", "Redo booking from recipient account: "+fmt.Sprintf("%d", t.FromAccount))
			return err
		}
	}

	query := "DELETE FROM transactions WHERE id=$1"

	_, e := cr.Exec(query, t.ID)

	if e != nil {
		var err err.Error
		err.Init("Transaction.Delete()", e.Error())
		return err
	}

	return err.Error{}
}

// ComputeFields computes the fields which are not directly received
// from the database
func (t *Transaction) computeFields(cr *sql.DB) {
	// Compute: FromAccountName
	if t.FromAccount != 0 {
		fromAccount, err := FindAccountByID(cr, t.FromAccount)

		if !err.Empty() {
			err.AddTraceback("Transaction.computeFields()", "Error while finding origin account by ID.")
			log.Println("[WARN]", err)
		}

		t.FromAccountName = fromAccount.Name
	} else {
		t.FromAccountName = "External Account"
	}

	// Compute: ToAccountName
	if t.ToAccount != 0 {
		toAccount, err := FindAccountByID(cr, t.ToAccount)

		if !err.Empty() {
			err.AddTraceback("Transaction.computeFields()", "Error while finding recipient account by ID:"+fmt.Sprintf("%d", t.CategoryID))
			log.Println("[WARN]", err)
		}

		t.ToAccountName = toAccount.Name
	} else {
		t.ToAccountName = "External Account"
	}

	// Compute: TransactionDateStr
	t.TransactionDateStr = t.TransactionDate.Format("02.01.2006 - 15:04")

	// Compute: Category
	if t.CategoryID > 0 {
		var err err.Error
		if t.Category, err = FindCategoryByID(cr, t.CategoryID); !err.Empty() {
			err.AddTraceback("Transaction.computeFields()", "Error while finding category by ID: "+fmt.Sprintf("%d", t.CategoryID))
			log.Println("[WARN]", err)
		}
	}
}

// FindByID finds a transaction with it's id
func (t *Transaction) FindByID(cr *sql.DB, transactionID int64) err.Error {
	query := "SELECT id, name, active, transaction_date, last_update, create_date, "
	query += "amount, account_id, to_account, transaction_type, description, category_id "
	query += "FROM transactions WHERE id=$1 "
	query += "ORDER BY transaction_date"

	var fromAccountID, toAccountID, categID interface{}

	e := cr.QueryRow(query, transactionID).Scan(
		&t.ID,
		&t.Name,
		&t.Active,
		&t.TransactionDate,
		&t.LastUpdate,
		&t.CreateDate,
		&t.Amount,
		&fromAccountID,
		&toAccountID,
		&t.TransactionType,
		&t.Description,
		&categID,
	)
	if e != nil {
		var err err.Error
		err.Init("Transaction.FindByID()", e.Error())
		return err
	}

	if fromAccountID != nil {
		t.FromAccount = fromAccountID.(int64)
	}

	if toAccountID != nil {
		t.ToAccount = toAccountID.(int64)
	}

	if categID != nil {
		t.CategoryID = categID.(int64)
	}

	t.computeFields(cr)

	return err.Error{}
}

// FindTransactionByID is similar to FindByID but returns the transaction
func FindTransactionByID(cr *sql.DB, transactionID int64) (Transaction, err.Error) {
	t := EmptyTransaction()

	e := t.FindByID(cr, transactionID)

	if !e.Empty() {
		e.AddTraceback("FindTransactionByID()", "Error while finding transaction by ID: "+fmt.Sprintf("%d", t.ID))
		return t, e
	}

	return t, err.Error{}
}

// GetAllTransactions does that what you expect
func GetAllTransactions(cr *sql.DB) ([]Transaction, err.Error) {
	var transactions []Transaction
	query := "SELECT id FROM transactions"

	idRows, e := cr.Query(query)

	if e != nil {
		var err err.Error
		err.Init("GetAllTransactions()", e.Error())
		return transactions, err
	}

	for idRows.Next() {
		var id int64

		if e = idRows.Scan(&id); e != nil {
			log.Println("[INFO] GetAllTransactions(): Skipping record")
			log.Printf("[WARN] GetAllTransactions: %s\n", e)
		} else {
			t := EmptyTransaction()

			if err := t.FindByID(cr, id); !err.Empty() {
				log.Printf("[INFO] GetAllTransactions(): Skipping record with ID %d\n", t.ID)
				log.Printf("[WARN] GetAllTransactions(): %s\n", err)
			} else {
				transactions = append(transactions, t)
			}
		}

	}

	return transactions, err.Error{}
}

// GetLatestTransactions returns a limited number of the latest transactions
// latest transactions are sorted by their transaction_date
func GetLatestTransactions(cr *sql.DB, amount int) ([]Transaction, err.Error) {
	var transactions []Transaction
	query := "SELECT id FROM transactions ORDER BY transaction_date DESC"

	if amount > 0 {
		query += " LIMIT $1"
	}

	var rows *sql.Rows
	var e error
	if amount > 0 {
		rows, e = cr.Query(query, amount)
	} else {
		rows, e = cr.Query(query)
	}

	if e != nil {
		var err err.Error
		err.Init("GetLatestTransactions()", e.Error())
		return transactions, err
	}

	for rows.Next() {
		var id int64

		if e = rows.Scan(&id); e != nil {
			log.Println("[INFO] GetLatestTransactions(): Skipping record")
			log.Printf("[WARN] GetLatestTransactions: %s\n", e)
		} else {
			t := EmptyTransaction()

			if err := t.FindByID(cr, id); !err.Empty() {
				log.Printf("[INFO] GetLatestTransactions(): Skipping record with ID: %d\n", t.ID)
				log.Println("[WARN]", err)
			} else {
				transactions = append(transactions, t)
			}
		}
	}

	return transactions, err.Error{}
}
