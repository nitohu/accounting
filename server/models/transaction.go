package models

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
	ID              int
	Name            string
	Active          bool
	TransactionDate string
	CreateDate      string
	LastUpdate      string
	Amount          float64
	FromAccount     string
	ToAccount       string
	TransactionType string
	User            string
}

// EmptyTransaction ..
func EmptyTransaction() Transaction {
	t := Transaction{
		ID:              0,
		Name:            "",
		Active:          false,
		TransactionDate: "",
		CreateDate:      "",
		LastUpdate:      "",
		Amount:          0.0,
		FromAccount:     "",
		ToAccount:       "",
		TransactionType: "",
		User:            "",
	}

	return t
}
