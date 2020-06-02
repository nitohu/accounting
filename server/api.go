package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nitohu/accounting/server/models"

	"github.com/nitohu/err"
)

// API Handler for the Accounting app
type API struct {
	id  int64
	obj interface{}
}

const (
	errorID    = "{'error': 'Please provide a valid ID.'}"
	errorGetID = "{'error': 'There was an unexpected error while getting the record by ID.'}"
)

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/api")[1]
	var body []byte

	w.Header().Set("Content-Type", "application/json; charset: utf-8")

	// Logging
	log.Printf("[INFO] %s: %s\n", r.URL.Path, r.Method)

	// Parse a body if there is one
	if r.ContentLength > 0 {
		b := r.Body
		body = make([]byte, r.ContentLength)
		if _, e := io.ReadFull(b, body); e != nil {
			var err err.Error
			err.Init("API.ServeHTTP()", e.Error())
			log.Println("[WARN]", err)
		}
	}

	api.obj = nil
	api.id = 0

	api.multiplexer(w, r, path, body)
}

func (api *API) multiplexer(w http.ResponseWriter, r *http.Request, path string, body []byte) {
	switch path {
	//
	// Categories
	//
	case "/categories":
		c := models.EmptyCategory()
		if len(body) > 0 {
			if err := json.Unmarshal(body, &c); err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "{'error': '%s'}", err)
				return
			}
		}
		api.id = c.ID
		if api.id > 0 {
			api.getCategoryByID(w, r)
			return
		}
		api.getCategories(w, r)
	case "/categories/create":
		c := models.EmptyCategory()
		if err := json.Unmarshal(body, &c); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		c.ID = 0
		api.obj = c
		api.updateCategory(w, r)
	case "/categories/update":
		c := models.EmptyCategory()
		if err := json.Unmarshal(body, &c); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		api.obj = c
		// Validate if the ID is existing
		if err := c.FindByID(db, c.ID); !err.Empty() {
			w.WriteHeader(500)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		api.id = c.ID
		api.updateCategory(w, r)
	case "/categories/delete":
		c := models.EmptyCategory()
		if err := json.Unmarshal(body, &c); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		api.id = c.ID
		api.deleteCategory(w, r)

	//
	// Accounts
	//
	case "/accounts":
		api.id = 0
		a := models.EmptyAccount()
		if len(body) > 0 {
			if err := json.Unmarshal(body, &a); err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "{'error': '%s'}", err)
				return
			}
		}
		api.id = a.ID
		if api.id > 0 {
			api.getAccountByID(w, r)
			return
		}
		api.getAccounts(w, r)
	case "/accounts/create":
		a := models.EmptyAccount()
		if err := json.Unmarshal(body, &a); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		a.ID = 0
		api.obj = a
		api.updateAccount(w, r)
	case "/accounts/update":
		a := models.EmptyAccount()
		if err := json.Unmarshal(body, &a); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		api.obj = a
		// Validate if the ID is existing
		if a.ID > 0 {
			if err := a.FindByID(db, a.ID); !err.Empty() {
				w.WriteHeader(500)
				fmt.Fprintf(w, "{'error': 'ID is not existing in the database'}")
				return
			}
			api.id = a.ID
		}
		api.updateAccount(w, r)
	case "/accounts/delete":
		a := models.EmptyAccount()
		if err := json.Unmarshal(body, &a); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		api.id = a.ID
		api.deleteAccount(w, r)
	//
	// Transactions
	//
	case "/transactions":
		api.id = 0
		t := models.EmptyTransaction()
		if len(body) > 0 {
			if e := json.Unmarshal(body, &t); e != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "{'error': '%s'}", e)
				return
			}
		}
		api.id = t.ID
		if api.id > 0 {
			api.getTransactionByID(w, r)
			return
		}
		api.getTransactions(w, r)
	case "/transactions/update":
		api.id = 0
		t := models.EmptyTransaction()
		if len(body) > 0 {
			if e := json.Unmarshal(body, &t); e != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "{'error': '%s'}", e)
				return
			}
		}
		// Validate if the ID is existing
		api.id = t.ID
		api.obj = t
		api.updateTransaction(w, r)
	case "/transactions/delete":
		api.id = 0
		t := models.Transaction{}
		if len(body) > 0 {
			if e := json.Unmarshal(body, &t); e != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "{'error': '%s'}", e)
				return
			}
		}
		api.id = t.ID
		api.deleteTransaction(w, r)
	case "/statistics":
		api.id = 0
		s := models.Statistic{}
		if len(body) > 0 {
			if e := json.Unmarshal(body, &s); e != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "{'error': '%s'}", e)
				return
			}
		}
		api.id = s.ID
		if api.id > 0 {
			api.getStatisticByID(w, r)
			return
		}
		api.getStatistics(w, r)
	default:
		w.WriteHeader(404)
		fmt.Fprint(w, "{'error': '404 Not Found', 'status': 404}")
	}
}

func (api API) sendResult(w http.ResponseWriter, data interface{}) {
	d, err := json.Marshal(data)
	if err != nil {
		log.Println("[ERROR] API.sendResult():", err)
		fmt.Fprint(w, "{'error': 'Server error while parsing data into JSON.'}")
		return
	}

	fmt.Fprint(w, string(d))
}

/*
	##############################
	#                            #
	#         Categories         #
	#                            #
	##############################
*/

// getCategories gets all Categories
// Also handles if
func (api API) getCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/categories/: Method must be GET.'}")
		return
	}

	c, err := models.GetAllCategories(db)
	if !err.Empty() {
		err.AddTraceback("API.getCategories()", "Error getting all categories.")
		log.Println("[ERROR]", err)
		w.WriteHeader(500)
		fmt.Fprint(w, "{'error': 'Server error while fetching accounts.'}")
		return
	}

	api.sendResult(w, c)
}

// getCategoryByID gets a category by it's ID
func (api API) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/categories/: Method must be POST for getting categories by ID.'}")
		return
	}
	if api.id <= 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, errorID)
		return
	}

	c := models.EmptyCategory()

	if err := c.FindByID(db, api.id); !err.Empty() {
		err.AddTraceback("API.getCategoryByID()", "Error getting category:"+fmt.Sprintf("%d", api.id))
		log.Println("[ERROR]", err)
		w.WriteHeader(400)
		fmt.Fprintln(w, errorGetID)
		return
	}

	api.sendResult(w, c)
}

// getCategoryByID write to a category
// if this category does not exist or has no ID, create one
// Method MUST be POST
// Returns the data written to the database
func (api API) updateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/categories/create: Method must be POST.'}")
		return
	}

	c := models.EmptyCategory()

	// Get data if ID is present
	if api.id > 0 {
		if err := c.FindByID(db, api.id); !err.Empty() {
			err.AddTraceback("API.updateCategory()", "Error getting category by ID: "+fmt.Sprintf("%d", api.id))
			log.Println("[WARN]", err)
			w.WriteHeader(400)
			fmt.Fprint(w, errorGetID)
			return
		}
	} else {
		c.CreateDate = time.Now()
	}

	// Create necessary variables for parsing the data
	reqData := api.obj.(models.Category)
	name := c.Name
	hex := c.Hex

	// Error catching when the client wants to create a category but provides no name
	if c.ID == 0 && reqData.Name == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "{'error': 'Please provide a name for creating a category.'}")
		return
	}

	// Update data only if it is given
	if reqData.Name != "" {
		name = reqData.Name
	}
	if reqData.Hex != "" {
		hex = reqData.Hex
	}

	// Update object
	c.Name = name
	c.Hex = hex
	c.LastUpdate = time.Now()

	// Save the object to the database
	var e err.Error
	if c.ID == 0 {
		e = c.Create(db)
	} else {
		e = c.Save(db)
	}
	if !e.Empty() {
		var err err.Error
		err.AddTraceback("API.updateCategory()", "Error creating/saving the category.")
		log.Println("[WARN]", err)
		w.WriteHeader(500)
		fmt.Fprint(w, "{'error': 'Error creating/saving the category.'}")
		return
	}

	api.sendResult(w, c)
}

// deletes a category with an ID
func (api API) deleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/categories/delete: Method must be DELETE.'}")
		return
	}
	if api.id <= 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, errorID)
		return
	}

	c := models.Category{
		ID: api.id,
	}

	if e := c.Delete(db); !e.Empty() {
		e.AddTraceback("API.deleteCategory", "Error deleting category with ID: "+fmt.Sprintf("%d", api.id))
		log.Println("[WARN]", e)
		w.WriteHeader(400)
		fmt.Fprintf(w, "{'error': 'There was an error deleting the record from the database.'}")
		return
	}

	fmt.Fprintf(w, "{'success': 'The record with the id %d was successfully deleted.'}", api.id)
}

/*
	##############################
	#                            #
	#          Accounts          #
	#                            #
	##############################
*/

// Gives back all accounts
func (api API) getAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/accounts/: Method must be GET.'}")
		return
	}

	acc, e := models.GetAllAccounts(db)
	if !e.Empty() {
		e.AddTraceback("API.getAccount()", "Error while getting accounts")
		log.Println("[ERROR]", e)
		w.WriteHeader(400)
		fmt.Fprint(w, "{'error': 'Server error while fetching accounts.'}")
		return
	}

	api.sendResult(w, acc)
}

// Returns a specific account with the given id
func (api API) getAccountByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/accounts/: Method must be GET.'}")
		return
	}
	if api.id <= 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, errorID)
		return
	}

	acc := models.EmptyAccount()
	if err := acc.FindByID(db, api.id); !err.Empty() {
		err.AddTraceback("API.getAccountByID", "Error while getting account: "+fmt.Sprintf("%d", api.id))
		log.Println("[ERROR]", err)

		w.WriteHeader(400)
		fmt.Fprint(w, errorGetID)
		return
	}

	api.sendResult(w, acc)
}

func (api API) updateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/categories/create: Method must be POST.'}")
		return
	}

	acc := models.EmptyAccount()

	if api.id > 0 {
		if err := acc.FindByID(db, api.id); !err.Empty() {
			err.AddTraceback("API.updateAccount()", "Error while getting account:"+fmt.Sprintf("%d", api.id))
			log.Println("[ERROR]", err)
			w.WriteHeader(400)
			fmt.Fprint(w, errorGetID)
			return
		}
	} else {
		acc.CreateDate = time.Now()
		acc.Active = true
	}

	reqData := api.obj.(models.Account)

	// Validate user input
	if acc.ID == 0 && reqData.Name == "" && reqData.Balance == 0.0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "{'error': 'Please provide a name and balance for creating a category.'}")
		return
	}

	// Set fields which are not empty
	if reqData.Name != "" {
		acc.Name = reqData.Name
	}
	if reqData.Balance != 0.0 {
		acc.Balance = reqData.Balance
	}
	if reqData.Iban != "" {
		acc.Iban = reqData.Iban
	}
	if reqData.BankCode != "" {
		acc.BankCode = reqData.BankCode
	}
	if reqData.AccountNr != "" {
		acc.AccountNr = reqData.AccountNr
	}
	if reqData.BankName != "" {
		acc.BankName = reqData.BankName
	}
	if reqData.BankType != "" {
		acc.BankType = reqData.BankType
	}

	acc.LastUpdate = time.Now()

	var e err.Error
	if acc.ID <= 0 {
		e = acc.Create(db)
	} else {
		e = acc.Save(db)
	}
	if !e.Empty() {
		e.AddTraceback("API.updateAccount()", "Error while writing account to the database.")
		log.Println("[ERROR]", e)
		return
	}

	api.sendResult(w, acc)
}

// deletes an account with an ID
func (api API) deleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/accounts/delete: Method must be DELETE.'}")
		return
	}
	if api.id <= 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, errorID)
		return
	}

	a, err := models.FindAccountByID(db, api.id)
	if !err.Empty() {
		err.AddTraceback("API.deleteAccount()", "Error getting account: "+fmt.Sprintf("%d", api.id))
		log.Println("[WARN]", err)
		w.WriteHeader(400)
		fmt.Fprintln(w, "{'error': 'There was an error finding the record in the database.'}")
		return
	}

	if a.TransactionCount > 0 {
		w.WriteHeader(403)
		fmt.Fprintln(w, "{'error': 'You cannot delete this record because it has transactions referenced to it'}")
		return
	}

	if err := a.Delete(db); !err.Empty() {
		err.AddTraceback("API.deleteAccount()", "Error while deleting account.")
		log.Println("[ERROR]", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "{'error': 'There was an error deleting the record from the database.'}")
		return
	}

	log.Printf("[INFO] api.deleteAccount(): Account with ID %d was successfully deleted.\n", api.id)
	fmt.Fprintf(w, "{'success': 'The record with the id %d was successfully deleted.'}", api.id)
}

/*
	##############################
	#                            #
	#        Transactions        #
	#                            #
	##############################
*/

// Returns all transactions
func (api API) getTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/transactions/: Method must be GET.'}")
		return
	}

	acc, e := models.GetLatestTransactions(db, -1)
	if !e.Empty() {
		e.AddTraceback("API.getTransactions()", "Error while getting transactions")
		log.Println("[ERROR]", e)
		fmt.Fprint(w, "{'error': 'Server error while fetching accounts.'}")
		return
	}

	api.sendResult(w, acc)
}

// Returns a specific transaction with the given id
func (api API) getTransactionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/transactions/: Method must be GET.'}")
		return
	}
	if api.id <= 0 {
		fmt.Fprint(w, errorID)
		return
	}

	t := models.EmptyTransaction()
	if err := t.FindByID(db, api.id); !err.Empty() {
		err.AddTraceback("API.getTransactionByID", "Error while getting transaction: "+fmt.Sprintf("%d", api.id))
		log.Println("[WARN]", err)
		fmt.Fprint(w, errorGetID)
		return
	}

	api.sendResult(w, t)
}

func (api API) updateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(400)
		fmt.Fprint(w, "{'error': '/api/transactions/update: Method must be POST.'}")
		return
	}

	t := models.EmptyTransaction()
	if api.id > 0 {
		if e := t.FindByID(db, api.id); !e.Empty() {
			e.AddTraceback("api.updateTransaction()", "Error while searching transaction per ID.")
			log.Println("[ERROR]", e)
			w.WriteHeader(400)
			fmt.Fprintf(w, errorGetID)
			return
		}
	} else {
		t.CreateDate = time.Now()
		t.LastUpdate = time.Now()
	}

	req := api.obj.(models.Transaction)

	if req.Name == "" || req.Amount <= 0 || (req.FromAccount <= 0 && req.ToAccount <= 0) {
		w.WriteHeader(400)
		fmt.Fprint(w, "{'error', 'At least one required field was empty (Name, Amount) or both accounts were 0'}")
		return
	}

	t.Name = req.Name
	t.Amount = req.Amount
	t.Active = req.Active

	var emptyTime time.Time

	if req.Description != "" {
		t.Description = req.Description
	}
	if req.TransactionDate != emptyTime {
		t.TransactionDate = req.TransactionDate
	}
	if req.TransactionType != "" {
		t.TransactionType = req.TransactionType
	}
	if req.CategoryID >= 0 {
		t.CategoryID = req.CategoryID
	}
	if req.FromAccount >= 0 {
		t.FromAccount = req.FromAccount
	}
	if req.ToAccount >= 0 {
		t.ToAccount = req.ToAccount
	}

	t.LastUpdate = time.Now()

	var e err.Error
	if api.id > 0 {
		e = t.Save(db)
	} else {
		e = t.Create(db)
	}
	if !e.Empty() {
		w.WriteHeader(400)
		fmt.Println("{'error': 'An error occured while saving/creating the transaction.'}")
		e.AddTraceback("api.updateTransaction()", "Error while creating/saving the transaction.")
		return
	}

	api.sendResult(w, t)
}

func (api API) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error', '/api/transactions/delete: Method must be DELETE.'}")
		return
	} else if api.id <= 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, errorID)
		return
	}

	// Validate that the ID is existing in the database
	t := models.Transaction{}

	if e := t.FindByID(db, api.id); !e.Empty() {
		w.WriteHeader(400)
		fmt.Fprint(w, errorGetID)
		e.AddTraceback("api.deleteTransaction()", "Error while finding record by ID: "+fmt.Sprintf("%d", api.id))
		log.Println("[WARN]", e)
		return
	}

	if e := t.Delete(db); !e.Empty() {
		w.WriteHeader(400)
		fmt.Fprint(w, "{'error', 'There was an unexpected error deleting the transaction from the database'}")
		e.AddTraceback("api.deleteTransaction()", "Error deleting the transaction "+fmt.Sprintf("%d", api.id))
		log.Println("[ERROR]", e)
		return
	}

	log.Printf("[INFO] api.deleteTransaction(): Transaction with ID %d was successfully deleted.\n", api.id)
	fmt.Fprintf(w, "{'success': 'The record with the id %d was successfully deleted.'}", api.id)
}

/*
	##############################
	#                            #
	#         Statistics         #
	#                            #
	##############################
*/

// Returns all Statistics
func (api API) getStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/statistics: Method must be GET.")
		return
	}

	stats, e := models.GetAllStatistics(db)
	if !e.Empty() {
		e.AddTraceback("api.getStatistics()", "Error while getting statistics.")
		log.Println("[ERROR]", e)
		fmt.Fprint(w, "{'error': 'Server error while getting the statistics.'}")
		return
	}

	api.sendResult(w, stats.GetArray())
}

// Returns specific Statistic
func (api API) getStatisticByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/statistics: Method must be GET.")
		return
	}
	if api.id <= 0 {
		fmt.Fprint(w, errorID)
		return
	}

	s := models.Statistic{}
	if e := s.FindByID(db, api.id); !e.Empty() {
		e.AddTraceback("api.getStatisticByID()", "Error getting Statistic by ID.")
		log.Println("[WARN]", e)
		fmt.Fprint(w, errorGetID)
		return
	}

	api.sendResult(w, s)
}
