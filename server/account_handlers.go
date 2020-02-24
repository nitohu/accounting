package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"./models"

	"github.com/gorilla/mux"
)

/*
	##############################
	#                            #
	#          Accounts          #
	#                            #
	##############################
*/

func handleAccountOverview(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountOverview", "%s", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx["Title"] = "Accounts"
	if ctx["Accounts"], err = models.GetAllAccounts(db); err != nil {
		logWarn("handleAccountOverView", "Error while getting accounts:\n%s", err)
	}

	err = tmpl.ExecuteTemplate(w, "accounts.html", ctx)

	if err != nil {
		logError("handleAccountOverview", "%s", err)
	}
}

func handleAccountForm(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountEditing", "%s", err)
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
		fmt.Println(idStr)
		if accountID, err = strconv.Atoi(idStr[0]); err != nil {
			log.Printf("[WARN] handleAccountForm(): %s\n", err)
		}

		// Error getting account by id, account probably doesn't exit
		// so return to the form
		if err := account.FindByID(db, int64(accountID)); err != nil {
			log.Println("[WARN] handleAccountForm():", err)
			account.ID = 0
		}
	}

	ctx["Title"] = "Edit " + account.Name
	ctx["Header"] = "Edit " + account.Name
	ctx["Btn"] = "Save Account"
	ctx["Account"] = account

	// Method is GET
	// Return the form
	if r.Method != http.MethodPost {
		err := tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		if err != nil {
			logError("handleAccountEditing", "%s", err)
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
		logError("handleAccountEditing", "Wrong value for bank type: %s", account.BankType)
		ctx["Error"] = "There was an error in the form."
		tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		return
	}

	// Save or create the account
	fmt.Println(account.ID)
	if account.ID <= 0 {
		if err := account.Create(db); err != nil {
			log.Println("handleAccountEditing():", err)
		}
	} else {
		if err = account.Save(db); err != nil {
			log.Println("handleAccountEditing():", err)
		}
	}

	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}

func handleAccountDeletion(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	_, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountDeletion", "%s", err)
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		logError("handleAccountDeletion", "%s", err)
		return
	}

	account, err := models.FindAccountByID(db, int64(id))

	if err != nil {
		logError("handleAccountDeletion", "%s", err)
		return
	}

	account.Delete(db)

}
