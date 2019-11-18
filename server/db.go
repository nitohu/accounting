package main

import (
	"database/sql"
	"fmt"

	"./models"
)

var db *sql.DB

func dbInit(host, user, password, dbname string, port int) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Printf("An error occurred while connecting to the database.\n"+
			"%s", err)
	}

	fmt.Printf("[INFO] Successfully connected to postgres!\n")
	return db
}

func dbCreateAccount(cr *sql.DB, account models.Account) int64 {
	query := "INSERT INTO accounts ( name, active, balance, balance_forecast, iban, account_holder,"
	query += " bank_code, account_nr, bank_name, bank_type, create_date, last_update, user_id"
	query += " ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);"

	res, err := cr.Exec(query,
		account.Name,
		account.Active,
		account.Balance,
		account.BalanceForecast,
		account.Iban,
		account.Holder,
		account.BankCode,
		account.AccountNr,
		account.BankName,
		account.BankType,
		account.CreateDate,
		account.LastUpdate,
		account.UserID,
	)

	if err != nil {
		msg := fmt.Sprintf("dbCreateAccount(): %s", err)
		panic(msg)
	}

	id, _ := res.LastInsertId()

	if rowCount, err := res.RowsAffected(); err != nil || rowCount < 1 {
		logError("dbCreateAccount", "No rows affected. ID: %s", id)
	}

	return id
}
