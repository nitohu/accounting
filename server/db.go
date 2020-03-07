package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"

	"github.com/nitohu/err"
)

var db *sql.DB

func dbInit(host, user, password, dbname string, port int) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println("[ERROR] dbInit():", err)
		return nil
	}

	fmt.Printf("[INFO] Successfully connected to postgres!\n")
	return db
}

// Authenticate the master password
func Authenticate(password string) (bool, err.Error) {
	var pw string
	query := "SELECT password FROM settings LIMIT 1"

	if e := db.QueryRow(query).Scan(&pw); e != nil {
		var err err.Error
		err.Init("Authenticate()", e.Error())
		return false, err
	}

	passw := sha256.Sum256([]byte(password))
	password = fmt.Sprintf("%X", passw)

	if pw == password {
		return true, err.Error{}
	}

	return false, err.Error{}
}
