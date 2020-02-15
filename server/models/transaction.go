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

func bookIntoAccount(cr *sql.DB, id int64, t *Transaction, invert bool) error {
	acc, err := FindAccountByID(cr, id)

	if err != nil {
		fmt.Println("Origin: bookIntoAccount()")
		fmt.Println("Hint: Finding Account")
		return err
	}

	err = acc.Book(cr, t, invert)

	if err != nil {
		fmt.Println("Origin: bookIntoAccount()")
		fmt.Println("Hint: Booking into Account")
		return err
	}

	return nil
}

// Create 's a transaction with the current values of the object
func (t *Transaction) Create(cr *sql.DB) error {
	// Requirements for creating a transaction
	if t.ID != 0 {
		return errors.New("This object already has an id")
	} else if t.Amount == 0.0 {
		return errors.New("The Amount of this transaction is 0")
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

	var err error

	// 4 cases of executing the query
	if t.ToAccount == 0 && t.FromAccount > 0 {
		err = cr.QueryRow(query,
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
		err = cr.QueryRow(query,
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
		err = cr.QueryRow(query,
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
		err = cr.QueryRow(query,
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

	if err != nil {
		fmt.Println("Traceback: models.Transaction.Create(): Inserting transaction into database")
		return err
	}

	// Writing id to object
	t.ID = id

	// Book transaction into FromAccount if it's given
	if t.FromAccount > 0 {
		err := bookIntoAccount(cr, t.FromAccount, t, true)

		if err != nil {
			fmt.Println("Traceback: models.Transaction.Create(): Book into FromAccount")
			return err
		}
	}

	// Book the transaction into ToAccount if it's given
	if t.ToAccount > 0 {
		err := bookIntoAccount(cr, t.ToAccount, t, false)

		if err != nil {
			fmt.Println("Traceback: models.Transaction.Create(): Book into ToAccount")
			return err
		}
	}

	return nil
}

// Save 's the current values of the object to the database
func (t *Transaction) Save(cr *sql.DB) error {
	if t.ID == 0 {
		fmt.Println("Traceback: models.Transaction.Create(): #1")
		return errors.New("This transaction as no ID, maybe create it first?")
	} else if t.Amount == 0.0 {
		fmt.Println("Traceback: models.Transaction.Create(): #2")
		return errors.New("The Amount of the transaction with the id " + fmt.Sprintf("%d", t.ID) + " is 0")
	}

	// Get old data
	var oldAmount float64
	var accountID, toAccountID interface{}
	var TransactionDateStr time.Time

	query := "SELECT amount, transaction_date, account_id, to_account FROM transactions WHERE id=$1"

	row := cr.QueryRow(query, t.ID)

	err := row.Scan(&oldAmount, &TransactionDateStr, &accountID, &toAccountID)

	if err != nil {
		fmt.Println("Traceback: models.Transaction.Create(): Fetching old data from database")
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

	t.TransactionDate = time.Now().Local()

	if t.ToAccount == 0 && t.FromAccount > 0 {
		_, err = cr.Exec(query,
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
		_, err = cr.Exec(query,
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
		_, err = cr.Exec(query,
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
		_, err = cr.Exec(query,
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

	if err != nil {
		fmt.Println("Traceback: models.Transaction.Create(): Write the new data to the database.")
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

			if err != nil {
				fmt.Println("Traceback: models.Transaction.Create(): Redo booking from old origin account.")
				return err
			}
		}

		// Book the amount into the new account
		err = bookIntoAccount(cr, t.FromAccount, t, true)
		originChanged = true

		if err != nil {
			fmt.Println("Traceback: models.Transaction.Create(): Booking to new origin account.")
			return err
		}
	}

	// The destination account changed
	if toAccountID != t.ToAccount {
		// Remove booking from old destination account
		// Make sure the old amount is used
		destinationChanged = true
		if toAccountID != nil {
			err := bookIntoAccount(cr, toAccountID.(int64), &temp, false)

			if err != nil {
				fmt.Println("Traceback: models.Transaction.Create(): Removing transaction from the old receiving account.")
				return err
			}
		}

		// Book into new account
		err = bookIntoAccount(cr, t.ToAccount, t, true)

		if err != nil {
			fmt.Println("Traceback: models.Transaction.Create(): Book transaction into new receiving account.")
			return err
		}
	}

	// The amount of the transaction changed
	if t.Amount != oldAmount {
		// The origin did not change
		// So book the difference into the origin account
		// If the origin has changed we've already booked into the accounts
		if !originChanged {
			// Calculate difference and book the difference into the origin account
			diff := t.Amount - oldAmount

			temp.Amount = diff

			err := bookIntoAccount(cr, t.FromAccount, &temp, true)

			if err != nil {
				fmt.Println("Traceback: models.Transaction.Create(): Booking difference into origin account.")
				return err
			}
		}

		// The destination did not change
		if !destinationChanged {
			diff := t.Amount - oldAmount

			temp.Amount = diff

			err := bookIntoAccount(cr, t.ToAccount, &temp, false)

			if err != nil {
				fmt.Println("Traceback: models.Transaction.Create(): Booking difference into destination account.")
				return err
			}
		}
	}

	return nil
}

// Delete 's the transtaction
func (t *Transaction) Delete(cr *sql.DB) error {
	if t.ID == 0 {
		return errors.New("The transaction you want to delete does not have an id")
	}

	if t.FromAccount > 0 {
		fmt.Println("delete transaction fromAccount")
		err := bookIntoAccount(cr, t.FromAccount, t, false)

		if err != nil {
			fmt.Println("Origin: Transaction.Delete")
			fmt.Println("Hint: booking into t.FromAccount")
			return err
		}
	}
	if t.ToAccount > 0 {
		fmt.Println("delete transaction toAccount")
		err := bookIntoAccount(cr, t.ToAccount, t, true)

		if err != nil {
			fmt.Println("Origin: Transaction.Delete")
			fmt.Println("Hint: booking into t.FromAccount")
			return err
		}
	}

	query := "DELETE FROM transactions WHERE id=$1"

	_, err := cr.Exec(query, t.ID)

	if err != nil {
		return err
	}

	return nil
}

// ComputeFields computes the fields which are not directly received
// from the database
func (t *Transaction) ComputeFields(cr *sql.DB) error {
	// Compute: FromAccountName
	if t.FromAccount != 0 {
		fromAccount, err := FindAccountByID(cr, t.FromAccount)

		if err != nil {
			return err
		}

		t.FromAccountName = fromAccount.Name
	} else {
		t.FromAccountName = "External Account"
	}

	// Compute: ToAccountName
	if t.ToAccount != 0 {
		toAccount, err := FindAccountByID(cr, t.ToAccount)

		if err != nil {
			return err
		}

		t.ToAccountName = toAccount.Name
	} else {
		t.ToAccountName = "External Account"
	}

	// Compute: TransactionDateStr
	t.TransactionDateStr = t.TransactionDate.Format("Monday 02 January 2006 - 15:04")

	return nil
}

// FindByID finds a transaction with it's id
func (t *Transaction) FindByID(cr *sql.DB, transactionID int64) error {
	query := "SELECT id, name, active, transaction_date, last_update, create_date, "
	query += "amount, account_id, to_account, transaction_type, description, category_id "
	query += "FROM transactions WHERE id=$1 "
	query += "ORDER BY transaction_date"

	var fromAccountID, toAccountID, categID interface{}

	err := cr.QueryRow(query, transactionID).Scan(
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

	if fromAccountID != nil {
		t.FromAccount = fromAccountID.(int64)
	}

	if toAccountID != nil {
		t.ToAccount = toAccountID.(int64)
	}

	if categID != nil {
		t.CategoryID = categID.(int64)
	}

	if err != nil {
		return err
	}

	t.ComputeFields(cr)

	if err != nil {
		return err
	}

	return nil
}

// FindTransactionByID is similar to FindByID but returns the transaction
func FindTransactionByID(cr *sql.DB, transactionID int64) (Transaction, error) {
	t := EmptyTransaction()

	err := t.FindByID(cr, transactionID)

	if err != nil {
		return t, err
	}

	return t, nil
}

// GetAllTransactions does that what you expect
func GetAllTransactions(cr *sql.DB) ([]Transaction, error) {
	var transactions []Transaction
	query := "SELECT id FROM transactions"

	idRows, err := cr.Query(query)

	if err != nil {
		return transactions, err
	}

	for idRows.Next() {
		var id int64

		if err = idRows.Scan(&id); err != nil {
			fmt.Printf("[INFO] Skipping Record\n[WARN] %s GetAllTransactions():\n%s", time.Now().Local(), err)
		} else {
			t := EmptyTransaction()

			if err = t.FindByID(cr, id); err != nil {
				fmt.Printf("[INFO] Skipping Record\n[WARN] %s GetAllTransactions():\n%s", time.Now().Local(), err)
			} else {
				transactions = append(transactions, t)
			}
		}

	}

	return transactions, nil
}

// GetLatestransactions returns a limited number of the latest transactions
// latest transactions are sorted by their transaction_date
func GetLatestTransactions(cr *sql.DB, amount int) ([]Transaction, error) {
	var transactions []Transaction
	query := "SELECT id FROM transactions ORDER BY transaction_date DESC LIMIT $1"

	rows, err := cr.Query(query, amount)

	if err != nil {
		fmt.Printf("[WARN] %s GetLatestTransactions(): Traceback: Error while running the query:\n%s", time.Now().Local(), err)
		return transactions, err
	}

	for rows.Next() {
		var id int64

		if err = rows.Scan(&id); err != nil {
			fmt.Printf("[INFO] Skipping Record\n[WARN] %s GetLatestTransactions():\n%s", time.Now().Local(), err)
		} else {
			t := EmptyTransaction()

			if err = t.FindByID(cr, id); err != nil {

				fmt.Printf("[INFO] Skipping Record\n[WARN] %s GetLatestTransactions(): Error while finding transaction by id:\n%s",
					time.Now().Local(), err)

			} else {
				transactions = append(transactions, t)
			}
		}
	}

	return transactions, nil
}
