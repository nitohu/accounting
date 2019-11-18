package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

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
