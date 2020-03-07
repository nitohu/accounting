package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/nitohu/err"

	"./models"
)

/*
	##############################
	#                            #
	#          Accounts          #
	#                            #
	##############################
*/

func handleAccountOverview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/accounts/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, e := createContextFromSession(db, session)

	if !e.Empty() {
		e.AddTraceback("handleAccountOverview", "Error creating context from session.")
		log.Println("[ERROR]", e)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx["Title"] = "Accounts"
	if ctx["Accounts"], e = models.GetAllAccounts(db); !e.Empty() {
		e.AddTraceback("handleAccountOverview", "Error while getting the accounts.")
		fmt.Println("[ERROR]", e)
	}

	err := tmpl.ExecuteTemplate(w, "accounts.html", ctx)

	if err != nil {
		e.Init("handleAccountOverview", err.Error())
		fmt.Println("[ERROR]", e)
	}
}

func handleAccountForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/accounts/form/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, e := createContextFromSession(db, session)

	if !e.Empty() {
		e.AddTraceback("handleAccountForm", "Error creating context from session.")
		log.Println("[ERROR]", e)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var accountID int

	account := models.EmptyAccount()
	vars := r.URL.Query()

	ctx["Title"] = "Create Account"
	ctx["Header"] = "Create Account"
	ctx["Btn"] = "Create Account"
	ctx["Account"] = account

	if idStr, ok := vars["id"]; ok {
		var e error
		if accountID, e = strconv.Atoi(idStr[0]); e != nil {
			var err err.Error
			err.Init("handleAccountForm()", e.Error())
			log.Println("[WARN]", err)
		}

		// Error getting account by id, account probably doesn't exit
		// so return to the form
		if e := account.FindByID(db, int64(accountID)); !e.Empty() {
			e.AddTraceback("handleAccountForm()", e.Error())
			log.Println("[WARN]", e)
			account.ID = 0
		} else {
			ctx["Title"] = "Edit Account"
			ctx["Header"] = "Edit " + account.Name
			ctx["Btn"] = "Save Account"
		}
	}

	ctx["Account"] = account

	// Method is GET
	// Return the form
	if r.Method != http.MethodPost {
		if e := tmpl.ExecuteTemplate(w, "account_form.html", ctx); e != nil {
			var err err.Error
			err.Init("handleAccountForm()", e.Error())
			log.Println("[ERROR]", err)
		}
		return
	}

	// Method is POST
	// Process the form
	account.Name = r.FormValue("name")
	account.Balance, _ = strconv.ParseFloat(r.FormValue("balance"), 64)
	account.BalanceForecast = account.Balance
	account.Iban = r.FormValue("iban")
	account.BankCode = r.FormValue("bankCode")
	account.AccountNr = r.FormValue("accountNumber")
	account.BankType = r.FormValue("accountType")

	if account.BankType == "bank" {
		account.BankName = r.FormValue("bankName")
	} else {
		account.BankName = r.FormValue("providerName")
	}

	if account.BankType != "bank" && account.BankType != "online" {
		var err err.Error
		err.Init("handleAccountForm", "Unkown bank type: "+account.BankType)
		log.Println("[WARN]", err)

		ctx["Error"] = "Unkown bank type: " + account.BankType
		tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		return
	}

	// Save or create the account
	if account.ID <= 0 {
		if err := account.Create(db); !err.Empty() {
			err.AddTraceback("handleAccountForm()", "Error while creating the account.")
			log.Println("[ERROR]", err)
		}
	} else {
		if err := account.Save(db); !err.Empty() {
			err.AddTraceback("handleAccountForm", "Error while saving the account.")
			log.Println("[ERROR]", err)
		}
	}

	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}
