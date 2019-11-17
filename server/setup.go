package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func initDb(host, user, password, dbname string, port int) *sql.DB {
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

func initTemplates(files ...string) (*template.Template, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}

	var paths []string

	if !strings.HasSuffix(workingDir, "/server") {
		workingDir += "/server"
	}

	for i := range files {
		file := files[i]
		path := workingDir + "/templates/" + file

		paths = append(paths, path)
	}

	tmpl, err := template.ParseFiles(paths...)

	return tmpl, err
}
