package main

import (
	"fmt"
	"log"
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
	if r.URL.Path != "/transactions/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, e := createContextFromSession(db, session)
	if !e.Empty() {
		e.AddTraceback("handleTransactionOverview()", "Error while creating the context.")
		fmt.Println("[ERROR]", e)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx["Title"] = "Transactions"
	if ctx["Transactions"], e = models.GetAllTransactions(db); !e.Empty() {
		e.AddTraceback("handleTransactionOverview()", "Error while getting all transactions.")
		fmt.Println("[WARN]", e)
	}

	if err := tmpl.ExecuteTemplate(w, "transactions.html", ctx); err != nil {
		e.Init("handleTransactionOverview", err.Error())
		fmt.Println("[ERROR]", e)
	}
	return
}

func handleTransactionForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/transactions/form/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if !err.Empty() {
		err.AddTraceback("handleTransactionForm()", "Error while creating the context.")
		fmt.Println("[ERROR]", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx["Title"] = "Create Transaction"
	ctx["Btn"] = "Create Transaction"

	vars := r.URL.Query()

	t := models.EmptyTransaction()

	// Get the current transaction
	if transactionID, ok := vars["id"]; ok {
		id, e := strconv.Atoi(transactionID[0])
		if e != nil {
			err.Init("handleTransactionForm()", e.Error())
			log.Println("[ERROR]", err)
			return
		}

		if err = t.FindByID(db, int64(id)); !err.Empty() {
			err.AddTraceback("handleTransactionForm()", "Error finding transaction: "+fmt.Sprintf("%d", id))
			log.Println("[WARN]", err)
			t.ID = 0
		} else {
			ctx["Title"] = "Edit " + t.Name
			ctx["Btn"] = "Save Transaction"
		}
	}

	ctx["Transaction"] = t

	// Get the accounts
	if ctx["Accounts"], err = models.GetAllAccounts(db); !err.Empty() {
		err.AddTraceback("handleTransactionForm()", "Error while getting the accounts.")
		log.Println("[WARN]", err)
	}

	// Get Categories
	if ctx["Categories"], err = models.GetAllCategories(db); !err.Empty() {
		err.AddTraceback("handleTransactionForm()", "Error while getting the categories.")
		log.Println("[WARN]", err)
	}

	if r.Method != http.MethodPost {
		e := tmpl.ExecuteTemplate(w, "transaction_form.html", ctx)
		if e != nil {
			err.Init("handleTransactionForm()", e.Error())
			log.Println("[ERROR]", err)
		}
		return
	}

	// Format the time received from the form
	tDate := r.FormValue("datetime")
	transactionDate, e := time.Parse(dtLayout, tDate)

	if e != nil {
		err.Init("handleTransactionForm()", e.Error())
		log.Println("[INFO] handleTransactionForm(): Using current time as transaction date.")
		log.Println("[WARN]", err)

		transactionDate = time.Now().Local()
	}

	t.Name = r.FormValue("name")

	if t.Amount, e = strconv.ParseFloat(r.FormValue("amount"), 64); e != nil {
		err.Init("handleTransactionsForm()", e.Error())
		log.Println("[WARN]", err)
		t.Amount = 0
	}
	if t.FromAccount, e = strconv.ParseInt(r.FormValue("fromAccount"), 0, 64); e != nil {
		err.Init("handleTransactionsForm()", e.Error())
		log.Println("[WARN]", err)
		t.FromAccount = 0
	}
	if t.ToAccount, e = strconv.ParseInt(r.FormValue("toAccount"), 0, 64); e != nil {
		err.Init("handleTransactionsForm()", e.Error())
		log.Println("[WARN]", err)
		t.ToAccount = 0
	}
	if t.CategoryID, e = strconv.ParseInt(r.FormValue("category"), 0, 64); e != nil {
		err.Init("handleTransactionsForm()", e.Error())
		log.Println("[WARN]", err)
		t.CategoryID = 0
	}
	t.LastUpdate = time.Now().Local()
	t.TransactionDate = transactionDate
	t.Description = r.FormValue("description")

	if t.ID == 0 {
		t.Active = true
		err = t.Create(db)
	} else {
		err = t.Save(db)
	}

	if !err.Empty() {
		err.AddTraceback("handleTransactionForm()", "Error while writing the transaction to the database.")
	}

	http.Redirect(w, r, "/transactions/", http.StatusSeeOther)
}

// TODO: Replace w/ API
func handleTransactionDeletion(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	_, err := createContextFromSession(db, session)

	if !err.Empty() {
		err.AddTraceback("handleTransactionDeletion()", "Error while creating the context.")
		log.Println("[ERROR]", err)
		http.Redirect(w, r, "/logout/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)

	id, e := strconv.Atoi(vars["id"])

	if e != nil {
		err.Init("handleTransactionDeletion()", e.Error())
		log.Println("[ERROR]", err)
		http.Error(w, "There was an error parsing the id", http.StatusInternalServerError)
		return
	}

	t, err := models.FindTransactionByID(db, int64(id))

	if !err.Empty() {
		err.AddTraceback("handleTransactionDeletion()", "Error while finding transaction by ID.")
		log.Println("[ERROR]", err)
		http.Error(w, "Model not found", http.StatusInternalServerError)
		return
	}

	err = t.Delete(db)

	if !err.Empty() {
		err.AddTraceback("handleTransactionDeletion()", "Error while deleting the transaction from the database.")
		log.Println("[ERROR]", err)
		http.Error(w, "Error while deleting", http.StatusInternalServerError)
	}
}
