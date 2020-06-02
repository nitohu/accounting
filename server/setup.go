package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/nitohu/err"
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

func getCmdLineArgs(args []string) (map[string]string, err.Error) {
	if len(args) <= 1 {
		return nil, err.Error{}
	}

	res := make(map[string]string)

	for i := 1; i < len(args); i += 2 {
		var kw, val string
		kw = args[i]
		if i+1 < len(args) {
			val = args[i+1]
		}
		if kw == "-h" || kw == "--help" {
			fmt.Println("Accounting Command line arguments")
			fmt.Println("\n\t-h, --help\t\tLoad the current page")
			fmt.Println("\t-c, --config <path>\tLoad the config file with the given <path>")
			fmt.Println("\t-p, --port <port>\tSet the port this app is running at (standard 80)")
			fmt.Println("\t-H, --dbhost <addr>\tSpecifies the host address of the postgres database")
			fmt.Println("\t-u, --dbuser <name>\tSpecifies the database user")
			fmt.Println("\t-pw, --dbpassword <pw>\tSpecifies the database password")
			fmt.Println("\t-P, --dbport <port>\tSpecifies the database port")
			os.Exit(0)
			return nil, err.Error{}
		} else if kw == "-c" || kw == "--config" {
			if e := readConfFile(&res, val); !e.Empty() {
				return nil, e
			}
		} else if kw == "-H" || kw == "--dbhost" {
			res["dbhost"] = val
		} else if kw == "-u" || kw == "--dbuser" {
			res["dbuser"] = val
		} else if kw == "-pw" || kw == "--dbpassword" {
			res["dbpassword"] = val
		} else if kw == "-P" || kw == "--dbport" {
			res["dbport"] = val
		} else if kw == "-p" || kw == "--port" {
			res["port"] = val
		} else if kw == "-d" || kw == "--database" {
			res["dbdatabase"] = val
		}
	}

	if err := validateCmdlineData(res); !err.Empty() {
		return nil, err
	}

	return res, err.Error{}
}

func readConfFile(res *map[string]string, filePath string) err.Error {
	data, error := ioutil.ReadFile(filePath)
	if error != nil {
		var e err.Error
		e.Init("getCmdLineArgs()", error.Error())
		return e
	}

	var kw, val string
	kwPassed := false
	for x := 0; x < len(data); x++ {
		c := data[x]
		if c == byte('=') {
			kwPassed = true
		} else if c == byte('\n') {
			if kw != "" && val != "" {
				(*res)[kw] = val
			}

			kwPassed = false
			kw = ""
			val = ""
		} else if c == byte(';') {
			var comment string
			for c != byte('\n') {
				x++
				c = data[x]
				comment += string(c)
			}
		} else if kwPassed {
			val += string(c)
		} else {
			kw += string(c)
		}
	}
	return err.Error{}
}

func validateCmdlineData(data map[string]string) err.Error {
	var err err.Error

	if val, ok := data["dbhost"]; !ok || val == "" {
		err.Init("validateCmdlineData()", "Please provide a host for the database.")
	}
	if val, ok := data["dbuser"]; !ok || val == "" {
		err.Init("validateCmdlineData()", "Please provide a user for the database.")
	}
	if val, ok := data["dbpassword"]; !ok || val == "" {
		err.Init("validateCmdlineData()", "Please provide a password for the database.")
	}
	if val, ok := data["dbhost"]; !ok || val == "" {
		err.Init("validateCmdlineData()", "Please provide a host for the database.")
	}
	if val, ok := data["dbdatabase"]; !ok || val == "" {
		err.Init("validateCmdlineData()", "Please provide a database name.")
	}

	return err
}
