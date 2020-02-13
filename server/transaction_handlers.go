package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"./models"

	"github.com/gorilla/mux"
)

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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
	fmt.Println(r.FormValue("datetime"))

	// Format the time received from the form
	tDate := r.FormValue("datetime")
	transactionDate, err := time.Parse(dtFormLayout, tDate)

	if err != nil {
		transactionDate = time.Now().Local()
		logWarn("handleTransactionForm", "Error while parsing the transaction date:\n%s", err)
	}

	t.Name = r.FormValue("name")
	// TODO: Handle error
	t.Amount, _ = strconv.ParseFloat(r.FormValue("amount"), 64)
	t.FromAccount, _ = strconv.ParseInt(r.FormValue("fromAccount"), 0, 64)
	t.LastUpdate = time.Now().Local()
	t.TransactionDate = transactionDate
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
