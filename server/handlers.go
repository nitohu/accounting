package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"./models"

	"github.com/gorilla/mux"
)

// Root
func handleRoot(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

	err = tmpl.ExecuteTemplate(w, "index.html", ctx)

	if err != nil {
		logError("handleLogin", "%s", err)
	}
}

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
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

	ctx["Title"] = "All Accounts"

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
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

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
	account.UserID = ctx["User"].(models.User).ID

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
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

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
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

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

/*
	##############################
	#                            #
	#        Transactions        #
	#                            #
	##############################
*/

func handleTransactionOverview(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleTransactionOverview", "%s", err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	ctx["Title"] = "Transactions"

	if r.Method != http.MethodPost {
		err := tmpl.ExecuteTemplate(w, "transactions.html", ctx)
		if err != nil {
			logError("handleTransactionOverview", "%s", err)
		}
		return
	}
}

func handleTransactionForm(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleTransactionForm", "%s", err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	t := models.EmptyTransaction()

	transactionID := vars["id"]
	id, _ := strconv.ParseInt(transactionID, 10, 64)
	t.FindByID(db, id)

	ctx["Title"] = "Edit Transaction"
	ctx["Transaction"] = t

	if r.Method != http.MethodPost {
		err = tmpl.ExecuteTemplate(w, "transaction_form.html", ctx)
		if err != nil {
			logError("handleTransactionForm", "%s", err)
		}
		return
	}

	t.Name = r.FormValue("name")
	// TODO: Handle error
	t.Amount, _ = strconv.ParseFloat(r.FormValue("amount"), 64)
	t.FromAccount, _ = strconv.ParseInt(r.FormValue("fromAccount"), 0, 64)
	t.LastUpdate = time.Now().Local()
	t.TransactionDDate = time.Now().Local()
	t.UserID = ctx["User"].(models.User).ID
	t.ToAccount = 0

	toAccount, err := strconv.ParseInt(r.FormValue("toAccount"), 0, 64)

	if err != nil {
		logError("handleTransactionForm", "%s", err)
	} else if toAccount != 0 {
		t.ToAccount = toAccount
	}

	if t.ID == 0 {
		t.Active = true
		err = t.Create(db)
	} else {
		err = t.Save(db)
	}

	if err != nil {
		logError("handleTransactionForm", "%s", err)
	}

	http.Redirect(w, r, "/transactions/", http.StatusSeeOther)
}

// TODO: Replace w/ API
func handleTransactionDeletion(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	_, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleTransactionDeletion", "%s", err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		logError("handleTransactionDeletion", "%s", err)
		http.Error(w, "There was an error parsing the id", http.StatusInternalServerError)
		return
	}

	t, err := models.FindTransactionByID(db, int64(id))

	if err != nil {
		logError("handleTransactionDeletion", "%s", err)
		http.Error(w, "Model not found", http.StatusInternalServerError)
		return
	}

	err = t.Delete(db)

	if err != nil {
		logError("handleTransactionDeletion", "%s", err)
		http.Error(w, "Error while deleting", http.StatusInternalServerError)
		return
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
	session, _ := store.Get(r, "session")

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		logInfo(r.URL.Path, "User is already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	u := models.EmptyUser()

	mail := r.FormValue("email")
	pw := r.FormValue("password")

	logInfo("handleLogin()", "Login try %s", mail)

	err := u.Login(db, mail, pw)

	if err != nil {
		fmt.Printf("[ERROR] %s %s", time.Now().Local(), err)
		tmpl.ExecuteTemplate(w, "login.html", map[string]string{"err": "Wrong credentials."})
		return
	}

	logInfo("handleLogin()", "Login try of %s was successful", mail)

	logInfo("handleLogin()", "UID: %d\n", u.ID)

	// Successfully logged in
	session.Values["authenticated"] = true
	session.Values["email"] = u.Email
	session.Values["uid"] = u.ID

	err = session.Save(r, w)

	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	session.Values["authenticated"] = false
	session.Values["user"] = models.EmptyUser()

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
