package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nitohu/err"
)

// APIHandler for the Accounting app
type APIHandler struct {
	id  int64
	obj interface{}
	key API
}

const (
	errorID    = "{'error': 'Please provide a valid ID.'}"
	errorGetID = "{'error': 'There was an unexpected error while getting the record by ID.'}"
)

func (api APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/api")[1]
	var body []byte

	w.Header().Set("Content-Type", "application/json; charset: utf-8")

	// Logging
	log.Printf("[INFO] %s: %s\n", r.URL.Path, r.Method)

	// Authorize
	if !api.authorize(w, r) {
		return
	}

	// Parse a body if there is one
	if r.ContentLength > 0 {
		b := r.Body
		body = make([]byte, r.ContentLength)
		if _, e := io.ReadFull(b, body); e != nil {
			var err err.Error
			err.Init("APIHandler.ServeHTTP()", e.Error())
			log.Println("[WARN]", err)
		}
	}

	api.obj = nil
	api.id = 0

	api.multiplexer(w, r, path, body)
}

func (api *APIHandler) authorize(w http.ResponseWriter, r *http.Request) bool {
	key := r.Header.Get("Authorization")
	if key != "" {
		k := strings.Split(key, " ")
		key := k[1]
		if k[0] == "Bearer" {
			var id int64
			var dbKey, fullKey string
			var local bool
			query := "SELECT id,api_key,local_key FROM api WHERE api_prefix=$1 AND active=true"
			k = strings.Split(key, ".")
			prefix := k[0]

			if e := db.QueryRow(query, prefix).Scan(&id, &dbKey, &local); e != nil {
				log.Println("[ERROR]", e)
			}

			if local {
				fullKey = key
			} else {
				fullKey = fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
			}

			if fullKey == dbKey {
				// Client is authenticated
				var a API
				if err := a.FindByPrefix(db, prefix); !err.Empty() {
					w.WriteHeader(400)
					err.AddTraceback("api.authorize", "Error while fetching API record.")
					log.Println("[ERROR]", err)
					fmt.Fprintln(w, "{'error': 'There was an unexpected error while fetching the API record.'}")
					return false
				}
				api.key = a
				query = "UPDATE api SET last_use=$1 WHERE id=$2"
				if _, err := db.Exec(query, time.Now(), id); err != nil {
					w.WriteHeader(400)
					log.Println("[ERROR]", err)
					fmt.Fprintln(w, "{'error': 'There was an internal server error.'}")
					return false
				}
				return true
			}
			w.WriteHeader(401)
			fmt.Fprintln(w, "{'error': 'The provided API Key is invalid.'}")
		} else {
			w.Header().Set("WWW-Authenticate", "Bearer realm=\"Valid API Key must be provided for access to the API\"")
			w.WriteHeader(401)
			fmt.Fprintln(w, "{'error': 'Please provide an API key in the request headers.'}")
		}
	} else {
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"Valid API Key must be provided for access to the API\"")
		w.WriteHeader(401)
		fmt.Fprintln(w, "{'error': 'Please provide an API key in the request headers.'}")
	}
	return false
}

func (api *APIHandler) checkAccessRight(w http.ResponseWriter, accessRight string) bool {
	if !StrContains(api.key.AccessRights, accessRight) {
		w.WriteHeader(403)
		fmt.Fprintf(w, "{'error': 'The provided access key does not have the mandatory rights to perform this action.'}")
		return false
	}
	return true
}

func (api *APIHandler) multiplexer(w http.ResponseWriter, r *http.Request, path string, body []byte) {
	switch path {
	//
	// Categories
	//
	case "/categories":
		if !api.checkAccessRight(w, "category.read") {
			return
		}
		c := EmptyCategory()
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
		if !api.checkAccessRight(w, "category.write") {
			return
		}
		c := EmptyCategory()
		if err := json.Unmarshal(body, &c); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		c.ID = 0
		api.obj = c
		api.updateCategory(w, r)
	case "/categories/update":
		if !api.checkAccessRight(w, "category.write") {
			return
		}
		c := EmptyCategory()
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
		if !api.checkAccessRight(w, "category.delete") {
			return
		}
		c := EmptyCategory()
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
		if !api.checkAccessRight(w, "account.read") {
			return
		}
		api.id = 0
		a := EmptyAccount()
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
		if !api.checkAccessRight(w, "account.write") {
			return
		}
		a := EmptyAccount()
		if err := json.Unmarshal(body, &a); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		a.ID = 0
		api.obj = a
		api.updateAccount(w, r)
	case "/accounts/update":
		if !api.checkAccessRight(w, "account.write") {
			return
		}
		a := EmptyAccount()
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
		if !api.checkAccessRight(w, "account.delete") {
			return
		}
		a := EmptyAccount()
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
		if !api.checkAccessRight(w, "transaction.read") {
			return
		}
		api.id = 0
		t := EmptyTransaction()
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
		if !api.checkAccessRight(w, "transaction.write") {
			return
		}
		api.id = 0
		t := EmptyTransaction()
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
		if !api.checkAccessRight(w, "transaction.delete") {
			return
		}
		api.id = 0
		t := Transaction{}
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
		if !api.checkAccessRight(w, "statistics.read") {
			return
		}
		api.id = 0
		s := Statistic{}
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

func (api APIHandler) sendResult(w http.ResponseWriter, data interface{}) {
	d, err := json.Marshal(data)
	if err != nil {
		log.Println("[ERROR] APIHandler.sendResult():", err)
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
func (api APIHandler) getCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/categories/: Method must be GET.'}")
		return
	}

	c, err := GetAllCategories(db)
	if !err.Empty() {
		err.AddTraceback("APIHandler.getCategories()", "Error getting all categories.")
		log.Println("[ERROR]", err)
		w.WriteHeader(500)
		fmt.Fprint(w, "{'error': 'Server error while fetching accounts.'}")
		return
	}

	api.sendResult(w, c)
}

// getCategoryByID gets a category by it's ID
func (api APIHandler) getCategoryByID(w http.ResponseWriter, r *http.Request) {
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

	c := EmptyCategory()

	if err := c.FindByID(db, api.id); !err.Empty() {
		err.AddTraceback("APIHandler.getCategoryByID()", "Error getting category:"+fmt.Sprintf("%d", api.id))
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
func (api APIHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/categories/create: Method must be POST.'}")
		return
	}

	c := EmptyCategory()

	// Get data if ID is present
	if api.id > 0 {
		if err := c.FindByID(db, api.id); !err.Empty() {
			err.AddTraceback("APIHandler.updateCategory()", "Error getting category by ID: "+fmt.Sprintf("%d", api.id))
			log.Println("[WARN]", err)
			w.WriteHeader(400)
			fmt.Fprint(w, errorGetID)
			return
		}
	} else {
		c.CreateDate = time.Now()
	}

	// Create necessary variables for parsing the data
	reqData := api.obj.(Category)
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
		err.AddTraceback("APIHandler.updateCategory()", "Error creating/saving the category.")
		log.Println("[WARN]", err)
		w.WriteHeader(500)
		fmt.Fprint(w, "{'error': 'Error creating/saving the category.'}")
		return
	}

	api.sendResult(w, c)
}

// deletes a category with an ID
func (api APIHandler) deleteCategory(w http.ResponseWriter, r *http.Request) {
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

	c := Category{
		ID: api.id,
	}

	if e := c.Delete(db); !e.Empty() {
		e.AddTraceback("APIHandler.deleteCategory", "Error deleting category with ID: "+fmt.Sprintf("%d", api.id))
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
func (api APIHandler) getAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/accounts/: Method must be GET.'}")
		return
	}

	acc, e := GetAllAccounts(db)
	if !e.Empty() {
		e.AddTraceback("APIHandler.getAccount()", "Error while getting accounts")
		log.Println("[ERROR]", e)
		w.WriteHeader(400)
		fmt.Fprint(w, "{'error': 'Server error while fetching accounts.'}")
		return
	}

	api.sendResult(w, acc)
}

// Returns a specific account with the given id
func (api APIHandler) getAccountByID(w http.ResponseWriter, r *http.Request) {
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

	acc := EmptyAccount()
	if err := acc.FindByID(db, api.id); !err.Empty() {
		err.AddTraceback("APIHandler.getAccountByID", "Error while getting account: "+fmt.Sprintf("%d", api.id))
		log.Println("[ERROR]", err)

		w.WriteHeader(400)
		fmt.Fprint(w, errorGetID)
		return
	}

	api.sendResult(w, acc)
}

func (api APIHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		fmt.Fprint(w, "{'error': '/api/categories/create: Method must be POST.'}")
		return
	}

	acc := EmptyAccount()

	if api.id > 0 {
		if err := acc.FindByID(db, api.id); !err.Empty() {
			err.AddTraceback("APIHandler.updateAccount()", "Error while getting account:"+fmt.Sprintf("%d", api.id))
			log.Println("[ERROR]", err)
			w.WriteHeader(400)
			fmt.Fprint(w, errorGetID)
			return
		}
	} else {
		acc.CreateDate = time.Now()
		acc.Active = true
	}

	reqData := api.obj.(Account)

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
		e.AddTraceback("APIHandler.updateAccount()", "Error while writing account to the database.")
		log.Println("[ERROR]", e)
		return
	}

	api.sendResult(w, acc)
}

// deletes an account with an ID
func (api APIHandler) deleteAccount(w http.ResponseWriter, r *http.Request) {
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

	a, err := FindAccountByID(db, api.id)
	if !err.Empty() {
		err.AddTraceback("APIHandler.deleteAccount()", "Error getting account: "+fmt.Sprintf("%d", api.id))
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
		err.AddTraceback("APIHandler.deleteAccount()", "Error while deleting account.")
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
func (api APIHandler) getTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/transactions/: Method must be GET.'}")
		return
	}

	acc, e := GetLatestTransactions(db, -1)
	if !e.Empty() {
		e.AddTraceback("APIHandler.getTransactions()", "Error while getting transactions")
		log.Println("[ERROR]", e)
		fmt.Fprint(w, "{'error': 'Server error while fetching accounts.'}")
		return
	}

	api.sendResult(w, acc)
}

// Returns a specific transaction with the given id
func (api APIHandler) getTransactionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/transactions/: Method must be GET.'}")
		return
	}
	if api.id <= 0 {
		fmt.Fprint(w, errorID)
		return
	}

	t := EmptyTransaction()
	if err := t.FindByID(db, api.id); !err.Empty() {
		err.AddTraceback("APIHandler.getTransactionByID", "Error while getting transaction: "+fmt.Sprintf("%d", api.id))
		log.Println("[WARN]", err)
		fmt.Fprint(w, errorGetID)
		return
	}

	api.sendResult(w, t)
}

func (api APIHandler) updateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(400)
		fmt.Fprint(w, "{'error': '/api/transactions/update: Method must be POST.'}")
		return
	}

	t := EmptyTransaction()
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

	req := api.obj.(Transaction)

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

func (api APIHandler) deleteTransaction(w http.ResponseWriter, r *http.Request) {
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
	t := Transaction{}

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
func (api APIHandler) getStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/statistics: Method must be GET.")
		return
	}

	stats, e := GetAllStatistics(db)
	if !e.Empty() {
		e.AddTraceback("api.getStatistics()", "Error while getting statistics.")
		log.Println("[ERROR]", e)
		fmt.Fprint(w, "{'error': 'Server error while getting the statistics.'}")
		return
	}

	api.sendResult(w, stats.GetArray())
}

// Returns specific Statistic
func (api APIHandler) getStatisticByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/statistics: Method must be GET.")
		return
	}
	if api.id <= 0 {
		fmt.Fprint(w, errorID)
		return
	}

	s := Statistic{}
	if e := s.FindByID(db, api.id); !e.Empty() {
		e.AddTraceback("api.getStatisticByID()", "Error getting Statistic by ID.")
		log.Println("[WARN]", e)
		fmt.Fprint(w, errorGetID)
		return
	}

	api.sendResult(w, s)
}
