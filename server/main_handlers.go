package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"./models"
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
	// pw := sha256.Sum256([]byte(r.FormValue("password")))

	settings.Currency = r.FormValue("currency")
	sdate := r.FormValue("salary_date")
	interval := r.FormValue("calc_interval")
	settings.CalcUoM = r.FormValue("calc_uom")
	settings.SetAPIKey(r.FormValue("api_key"))

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

	if !err.Empty() {
		err.AddTraceback("handleSettings", "Error while saving the settings to the database.")
		log.Println("[ERROR]", err)
	}

	e = tmpl.ExecuteTemplate(w, "settings.html", ctx)

	if e != nil {
		err.Init("handleSettings()", e.Error())
		log.Println("[ERROR]", err)
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
