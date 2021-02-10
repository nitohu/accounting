package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nitohu/accounting/server/models"
	"github.com/nitohu/err"
)

// Root
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if !err.Empty() {
		err.AddTraceback("handleRoot", "Error while creating the context.")
		log.Println("[ERROR]", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx["Title"] = "Dashboard"

	if ctx["Transactions"], err = models.GetLatestTransactions(db, 10); !err.Empty() {
		err.AddTraceback("handleRoot", "Error while getting the latest transactions.")
		log.Println("[WARN]", err)
	}
	if ctx["Accounts"], err = models.GetLimitAccounts(db, 4); !err.Empty() {
		err.AddTraceback("handleRoot", "Error while getting accounts.")
		log.Println("[WARN]", err)
	}
	if ctx["Statistics"], err = models.GetAllStatistics(db); !err.Empty() {
		err.AddTraceback("handleRoot", "Error while getting statistics.")
		log.Println("[WARN]", err)
	}

	e := tmpl.ExecuteTemplate(w, "index.html", ctx)

	if e != nil {
		err.Init("handleRoot", e.Error())
		log.Println("[ERROR]", err)
	}
}

// Settings
func handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/settings/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)
	ctx["Title"] = "Settings"

	if !err.Empty() {
		err.AddTraceback("handleSettings", "Error while creating the context.")
		log.Println("[WARN]", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		return
	}

	settings := ctx["Settings"].(models.Settings)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodPost {
		e := tmpl.ExecuteTemplate(w, "settings.html", ctx)

		if e != nil {
			err.Init("handleSettings", e.Error())
			log.Println("[ERROR]", err)
		}
		return
	}

	// Receive data from from
	settings.Name = r.FormValue("name")
	settings.Email = r.FormValue("email")
	// Old and new password
	opw := r.FormValue("old_password")
	npw := r.FormValue("password")
	if opw != "" && npw != "" {
		if err := settings.UpdateMasterPassword(db, opw, npw); !err.Empty() {
			log.Println("[ERROR]", err.Error())
			ctx["Error"] = "Error: " + err.Error()
			if e := tmpl.ExecuteTemplate(w, "settings.html", ctx); e != nil {
				err.Init("handleSettings", e.Error())
				log.Println("[ERROR]" + err.Error())
			}
			return
		}
	}

	settings.Currency = r.FormValue("currency")
	sdate := r.FormValue("salary_date")
	interval := r.FormValue("calc_interval")
	settings.CalcUoM = r.FormValue("calc_uom")
	// settings.SetAPIKey(r.FormValue("api_key"))
	// api_key := r.FormValue("api_key")

	// Converting the hashed password to a string
	// password := fmt.Sprintf("%X", pw)
	startDate, e := time.Parse(dateSettingsLayout, sdate)

	settings.CalcInterval, _ = strconv.ParseInt(interval, 10, 64)

	if e != nil {
		err.Init("handleSettings", e.Error())
		log.Println("[INFO] Using current time as start date.")
		log.Println("[WARN]", err)
		startDate = time.Now()
	}

	settings.SalaryDate = startDate

	// err = settings.Save(db, password)
	err = settings.Save(db)

	ctx["Settings"] = settings
	ctx["Success"] = "Settings saved."

	if !err.Empty() {
		err.AddTraceback("handleSettings", "Error while saving the settings to the database.")
		ctx["Error"] = "There was an error while saving the settings to the database. Please check the logs."
		ctx["Success"] = ""
		log.Println("[ERROR]", err)
	}

	e = tmpl.ExecuteTemplate(w, "settings.html", ctx)

	if e != nil {
		err.Init("handleSettings()", e.Error())
		log.Println("[ERROR]", err)
	}
}

// API Settings Overview
func handleAPISettingsOverview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/settings/api/" {
		handleNotFound(w, r)
		return
	}

	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)
	if !err.Empty() {
		err.AddTraceback("handleAPISettingsOverview", "Error while creating the context.")
		log.Println("[WARN]", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	apiKeys, err := models.GetAllAPIKeys(db)
	if !err.Empty() {
		err.AddTraceback("handleAPISettingsOverview", "Error while fetching API Keys.")
		log.Println("[ERROR]", err)
	}
	ctx["Title"] = "API Settings"
	ctx["APIKeys"] = apiKeys

	if r.Method != http.MethodPost {
		if e := tmpl.ExecuteTemplate(w, "settings_api.html", ctx); e != nil {
			err.Init("handleAPISettingsOverview", e.Error())
			log.Println("[ERROR]", err)
		}
		return
	}
}

// API Settings Form
func handleAPISettings(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/settings/api/form/" {
		handleNotFound(w, r)
		return
	}

	session, _ := store.Get(r, "session")
	ctx, e := createContextFromSession(db, session)
	if !e.Empty() {
		e.AddTraceback("handleAPISettings", "Error while creating the context.")
		log.Println("[WARN]", e)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx["Title"] = "Create API Key"
	ctx["Btn"] = "Create"

	vars := r.URL.Query()
	key := models.API{}
	if keyID, ok := vars["id"]; ok {
		k, err := strconv.Atoi(keyID[0])
		if err != nil {
			log.Println("[WARN] handleAPISettings(): ", e)
		}
		ctx["Title"] = "Edit API Key"
		ctx["Btn"] = "Save"
		if e := key.FindByID(db, int64(k)); !e.Empty() {
			e.AddTraceback("handleAPISettings", "Error while gettings API Key.")
			log.Println("[WARN]", e)
		}
	}
	ctx["API"] = key

	// Format access rights for rendering
	var arf [][]string
	ar := models.GetAllAccessRights()
	for x := 0; x < 3; x++ {
		var row []string
		for y := 0; y < (len(ar) / 3); y++ {
			i := y + x*(len(ar)/3)
			row = append(row, ar[i])
		}
		arf = append(arf, row)
	}
	ctx["AccessRights"] = arf
	ctx["FormattedAR"] = key.FormatAccessRights()

	if r.Method != http.MethodPost {
		e := tmpl.ExecuteTemplate(w, "settings_api_form.html", ctx)
		if e != nil {
			var err err.Error
			err.Init("handleAPISettings()", e.Error())
			log.Println("[ERROR]", err)
		}
		return
	}

	active := r.FormValue("active")
	local := r.FormValue("local")
	key.Name = r.FormValue("name")

	// Validate access rights, if they are invalid return an error to the user
	key.AccessRights = strings.Split(r.FormValue("access_rights"), ";")
	if !models.ValidateAccessRights(key.AccessRights) {
		ctx["Error"] = "One of the access rights you've typed in is not in the list of access rights."

		er := tmpl.ExecuteTemplate(w, "settings_api_form.html", ctx)
		if er != nil {
			var err err.Error
			err.Init("handleAPISettings()", e.Error())
			log.Println("[ERROR]", err)
		}
		return
	}

	// Convert FormValue strings into bool
	if active == "on" {
		key.Active = true
	} else {
		key.Active = false
	}
	if local == "on" {
		key.LocalKey = true
	} else {
		key.LocalKey = false
	}

	log.Printf("Name: %s Active: %t Local: %t\n", key.Name, key.Active, key.LocalKey)

	// Save the key and either render the next page
	if key.ID > 0 {
		if e = key.Save(db); !e.Empty() {
			fmt.Println("[WARN]", e.Error())
			ctx["Error"] = "There was an error while saving this item to the database, please check the logs."

			if er := tmpl.ExecuteTemplate(w, "settings_api_form.html", ctx); er != nil {
				var err err.Error
				err.Init("handleAPISettings()", e.Error())
				log.Println("[ERROR]", err)
				return
			}
		}
		http.Redirect(w, r, "/settings/api/", http.StatusSeeOther)
	} else {
		ctx["RawKey"] = key.GenerateAPIKey()

		if e = key.Create(db); !e.Empty() {
			ctx["Error"] = "There was an error while saving this item to the database, please check the logs."
			fmt.Println("[WARN]", e.Error())
		} else {
			ctx["Success"] = "Success"
		}

		if er := tmpl.ExecuteTemplate(w, "settings_api_form.html", ctx); er != nil {
			var err err.Error
			err.Init("handleAPISettings()", e.Error())
			log.Println("[ERROR]", err)
		}
	}
}

// 404 Page
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if !err.Empty() {
		err.AddTraceback("handleNotFound()", "Error creating context.")
		log.Println(err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(404)

	ctx["Title"] = "404 - Not Found"

	e := tmpl.ExecuteTemplate(w, "404.html", ctx)

	if e != nil {
		err.Init("handleNotFound", e.Error())
		log.Println(err)
	}
}

/*
	##############################
	#                            #
	#            User            #
	#                            #
	##############################
*/

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		log.Println("[INFO] User is already logged in. Redirecting....")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	pw := r.FormValue("password")

	// Authenticate
	success, err := Authenticate(pw)
	if !err.Empty() {
		err.AddTraceback("handleLogin()", "Error while authenticating the user.")
		log.Println("[ERROR]", err)
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	if success == false {
		log.Println("[INFO] handleLogin(): Failed login attempt.")
		msg := make(map[string]string)
		msg["Error"] = "The authentication was not successful."
		tmpl.ExecuteTemplate(w, "login.html", msg)
		return
	}

	// Authentication was successful
	// Generate and save session key
	sessionKey := GenerateSessionKey()

	query := "UPDATE settings SET session_key=$1;"

	_, e := db.Exec(query, sessionKey)

	if e != nil {
		err.Init("handleLogin()", e.Error())
		log.Println("[ERROR]", err)

		msg := make(map[string]string)
		msg["Error"] = "There was an error while saving your session key. Please try again."
		tmpl.ExecuteTemplate(w, "login.html", msg)
		return
	}

	// Successfully logged in
	session.Values["authenticated"] = true
	session.Values["key"] = sessionKey

	e = session.Save(r, w)

	if e != nil {
		err.Init("handleLogin", e.Error())
		fmt.Println("[ERROR]", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	session.Values["authenticated"] = false
	session.Values["key"] = ""

	session.Save(r, w)

	http.Redirect(w, r, "/login/", http.StatusSeeOther)
}
