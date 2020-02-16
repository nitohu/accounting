package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"./models"
)

// Root
func handleRoot(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logWarn("handleRoot", "%s", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if ctx["Transactions"], err = models.GetLatestTransactions(db, 10); err != nil {
		logWarn("handleRoot", "Error while getting transactions:\n%s", err)
	}
	if ctx["Accounts"], err = models.GetLimitAccounts(db, 4); err != nil {
		logWarn("handleRoot", "Error while gettings Accounts:\n%s", err)
	}

	err = tmpl.ExecuteTemplate(w, "index.html", ctx)

	if err != nil {
		logError("handleLogin", "%s", err)
	}
}

// Settings
func handleSettings(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)
	ctx["Title"] = "Settings"

	if err != nil {
		logWarn("handleSettings", "%s", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
	}

	settings := ctx["Settings"].(models.Settings)

	fmt.Println(settings.StartDate)
	fmt.Println(settings.StartDateForm)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodPost {
		err = tmpl.ExecuteTemplate(w, "settings.html", ctx)

		if err != nil {
			logError("handleSettings", "Error while executing the template:\n%s", err)
		}
		return
	}

	// Receive data from from
	settings.Name = r.FormValue("name")
	settings.Email = r.FormValue("email")
	// pw := sha256.Sum256([]byte(r.FormValue("password")))

	settings.Currency = r.FormValue("currency")
	sdate := r.FormValue("start_date")
	interval := r.FormValue("calc_interval")
	settings.CalcUoM = r.FormValue("calc_uom")

	// Converting the hashed password to a string
	// password := fmt.Sprintf("%X", pw)
	startDate, err := time.Parse(dateFormLayout, sdate)

	settings.CalcInterval, _ = strconv.ParseInt(interval, 10, 64)

	if err != nil {
		logWarn("handleSettings", "Error while parsing the start_date:\n%s", err)
		startDate = time.Now()
	}

	settings.StartDate = startDate

	// err = settings.Save(db, password)
	err = settings.Save(db)

	ctx["Settings"] = settings

	if err != nil {
		logError("handleSettings", "Error while saving the settings:\n%s", err)
	}

	ctx["SaveSuccess"] = true

	err = tmpl.ExecuteTemplate(w, "settings.html", ctx)

	if err != nil {
		logError("handleLogin()", "Error while rendering the template:\n%s", err)
	}
}

// 404 Page
func pageNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		ctx, err := createContextFromSession(db, session)

		if err != nil {
			logError("handlePageNotFound", "%s", err)
			http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = tmpl.ExecuteTemplate(w, "404.html", ctx)

		if err != nil {
			logError("handlePageNotFound", "%s", err)
		}
	})
}

/*
	##############################
	#                            #
	#            User            #
	#                            #
	##############################
*/

func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		logInfo(r.URL.Path, "User is already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	pw := r.FormValue("password")

	// Authenticate
	success, err := Authenticate(pw)

	if err != nil {
		logWarn("handleLogin", "Error while authenticating:\n%s", err)
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	if success == false {
		msg := make(map[string]string)
		msg["Error"] = "The authentication was not successful."
		tmpl.ExecuteTemplate(w, "login.html", msg)
		return
	}

	// Authentication was successful
	// Generate and save session key
	sessionKey := GenerateSessionKey()

	query := "UPDATE settings SET session_key=$1;"

	_, err = db.Exec(query, sessionKey)

	if err != nil {
		logError("handleLogin", "Error while saving session key:%s", err)
		msg := make(map[string]string)
		msg["Error"] = "There was an error while saving your session key. Please try again."
		tmpl.ExecuteTemplate(w, "login.html", msg)
		return
	}

	// Successfully logged in
	session.Values["authenticated"] = true
	session.Values["key"] = sessionKey

	err = session.Save(r, w)

	if err != nil {
		logWarn("handleLogin", "%s", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	session.Values["authenticated"] = false
	session.Values["key"] = ""

	session.Save(r, w)

	http.Redirect(w, r, "/login/", http.StatusSeeOther)
}
