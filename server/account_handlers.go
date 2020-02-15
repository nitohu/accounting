package main

import (
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

	ctx["Title"] = "All Accounts"
	if ctx["Accounts"], err = models.GetAllAccounts(db); err != nil {
		logWarn("handleAccountOverView", "Error while getting accounts:\n%s", err)
	}

	err = tmpl.ExecuteTemplate(w, "accounts.html", ctx)

	if err != nil {
		logError("handleAccountOverview", "%s", err)
	}
}

func handleAccountCreation(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountCreation", "%s", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx["Title"] = "Create Account"
	ctx["Header"] = "Create a new account"
	ctx["Btn"] = "Create Account"
	ctx["Account"] = models.EmptyAccount()

	if r.Method != http.MethodPost {
		err = tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		if err != nil {
			logError("handleAccountCreation", "%s", err)
		}
		return
	}

	account := models.EmptyAccount()

	account.Name = r.FormValue("name")
	account.Balance, _ = strconv.ParseFloat(r.FormValue("balance"), 64)
	account.BalanceForecast = account.Balance
	account.Iban = r.FormValue("iban")
	account.BankCode = r.FormValue("bankCode")
	account.AccountNr = r.FormValue("accountNumber")
	account.BankName = r.FormValue("bankName")
	account.BankType = r.FormValue("accountType")

	err = account.Create(db)

	if err != nil {
		logError("handleAccountCreation", "%s", err)
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func handleAccountEditing(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountEditing", "%s", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])

	account := models.EmptyAccount()
	account.FindByID(db, int64(accountID))

	ctx["Title"] = "Edit Account: " + account.Name
	ctx["Header"] = "Edit " + account.Name
	ctx["Btn"] = "Save Account"
	ctx["Account"] = account

	if r.Method != http.MethodPost {
		err := tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		if err != nil {
			logError("handleAccountEditing", "%s", err)
		}
		return
	}

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
		err = tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		return
	}

	err = account.Save(db)

	if err != nil {
		logError("handleAccountEditing", "%s", err)
		http.Redirect(w, r, "/accounts/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/accounts/", http.StatusSeeOther)
}

func handleAccountDeletion(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	_, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountDeletion", "%s", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
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
