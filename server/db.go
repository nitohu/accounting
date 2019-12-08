package main

import (
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
