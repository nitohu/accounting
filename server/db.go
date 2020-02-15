package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
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

// Authenticate the master password
func Authenticate(password string) (bool, error) {
	var pw string
	query := "SELECT password FROM settings LIMIT 1"

	if err := db.QueryRow(query).Scan(&pw); err != nil {
		logWarn("Authenticate()", "Traceback: failed to get master password from database.")
		return false, err
	}

	passw := sha256.Sum256([]byte(password))
	password = fmt.Sprintf("%X", passw)

	if pw == password {
		return true, nil
	}

	return false, nil
}
