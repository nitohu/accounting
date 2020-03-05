package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"./models"
)

// API Handler for the Accounting app
type API struct {
	id  int64
	obj interface{}
}

// APIAccount holds the account type and additional fields for parsing
// from JSON
type APIAccount struct {
	models.Account

	Balance string
	Active  string
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
		if _, err := io.ReadFull(b, body); err != nil {
			fmt.Println("[ERROR] API.ServeHTTP():", err)
		}

		// body = api.replaceBodyJSON(body)
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
		if err := c.FindByID(db, c.ID); err != nil {
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
		a := APIAccount{
			Account: models.EmptyAccount(),
		}
		if err := json.Unmarshal(body, &a); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		a.ID = 0
		api.obj = a
		api.updateAccount(w, r)
	default:
		w.WriteHeader(404)
		fmt.Fprint(w, "{'error': '404 Not Found', 'status': 404}")
	}
}

// func (api API) replaceBodyJSON(body []byte) []byte {
// 	bodyStr := string(body)

// 	bodyStr = strings.Replace(bodyStr, "Balance", "apiBalance", -1)
// 	bodyStr = strings.Replace(bodyStr, "Active", "apiActive", -1)

// 	return []byte(bodyStr)
// }

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
		fmt.Fprint(w, "{'error': '/api/categories/: Method must be GET.'}")
		return
	}

	c, err := models.GetAllCategories(db)
	if err != nil {
		fmt.Fprint(w, "{'error': 'Server error while fetching accounts.'}")
		log.Printf("[ERROR] API.getCategories(): Error getting Categories.\n%s\n", err)
		return
	}

	api.sendResult(w, c)
}

// getCategoryByID gets a category by it's ID
func (api API) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/categories/: Method must be GET.'}")
		return
	}
	if api.id <= 0 {
		fmt.Fprintln(w, errorID)
		return
	}

	c := models.EmptyCategory()

	if err := c.FindByID(db, int64(api.id)); err != nil {
		log.Println("[ERROR] API.getCategoryByID():", err)
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
		fmt.Fprint(w, "{'error': '/api/categories/create: Method must be POST.'}")
		return
	}

	c := models.EmptyCategory()

	// Get data if ID is present
	if api.id > 0 {
		if err := c.FindByID(db, api.id); err != nil {
			log.Println("[WARN] API.updateCategory():", err)
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
	var err error
	if c.ID == 0 {
		err = c.Create(db)
	} else {
		err = c.Save(db)
	}
	if err != nil {
		log.Println("[WARN] API.updateCategory()", err)
		fmt.Fprint(w, "{'error': 'Error creating/saving the category.'}")
		return
	}

	api.sendResult(w, c)
}

// deletes a category with an ID
func (api API) deleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		fmt.Fprint(w, "{'error': '/api/categories/delete: Method must be DELETE.'}")
		return
	}
	if api.id <= 0 {
		fmt.Fprintln(w, errorID)
		return
	}

	c := models.Category{
		ID: api.id,
	}

	if err := c.Delete(db); err != nil {
		log.Println("[WARN] API.deleteCategory():", err)
		fmt.Fprintf(w, "{'error': 'There was an error deleting the record from the database.'}")
		return
	}

	fmt.Fprintf(w, "{'success': 'You've successfully deleted the record with the id %d'}", api.id)
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
		fmt.Fprint(w, "{'error': '/api/accounts/: Method must be GET.'}")
		return
	}

	acc, err := models.GetAllAccounts(db)
	if err != nil {
		log.Println("[ERROR] API.getAccounts():", err)
		fmt.Fprint(w, "{'error': 'Server error while fetching accounts.'}")
		return
	}

	api.sendResult(w, acc)
}

// Returns a specific account with the given id
func (api API) getAccountByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/api/accounts/: Method must be GET.'}")
		return
	}
	if api.id <= 0 {
		fmt.Fprint(w, errorID)
		return
	}

	acc := models.EmptyAccount()
	if err := acc.FindByID(db, api.id); err != nil {
		fmt.Fprint(w, errorGetID)
		return
	}

	api.sendResult(w, acc)
}

func (api API) updateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, "{'error': '/api/categories/create: Method must be POST.'}")
		return
	}

	acc := APIAccount{
		Account: models.EmptyAccount(),
	}

	if api.id > 0 {
		if err := acc.Account.FindByID(db, api.id); err != nil {
			log.Println("[WARN] API.updateCategory():", err)
			fmt.Fprint(w, errorGetID)
			return
		}
	} else {
		acc.Account.CreateDate = time.Now()
		acc.Account.Active = true
	}

	reqData := api.obj.(APIAccount)

	fmt.Println(reqData)

	if acc.ID == 0 && reqData.Name == "" && reqData.Balance == "" {
		fmt.Fprint(w, "{'error': 'Please provide a name and balance for creating a category.'}")
		return
	}

	active := acc.Account.Active
	name := acc.Account.Name
	balance := acc.Account.Balance
	iban := acc.Account.Iban
	bankCode := acc.Account.BankCode
	accountNr := acc.Account.AccountNr
	bankName := acc.Account.BankName
	bankType := acc.Account.BankType

	if reqData.Name != "" {
		name = reqData.Name
	}
	if reqData.Balance != "" {
		b, err := strconv.ParseFloat(reqData.Balance, 64)
		if err != nil {
			log.Printf("[WARN] API.updateAccount(): Balance will keep the old value:\n%s\n", err)
		} else {
			balance = b
		}
	}
	if reqData.Active != "" {
		a, err := strconv.ParseBool(reqData.Active)
		if err != nil {
			log.Printf("[WARN] API.updateAccount(): Active will keep the old value:\n%s\n", err)
		} else {
			active = a
		}
	}
	if reqData.Iban != "" {
		iban = reqData.Iban
	}
	if reqData.BankCode != "" {
		bankCode = reqData.BankCode
	}
	if reqData.AccountNr != "" {
		accountNr = reqData.AccountNr
	}
	if reqData.BankName != "" {
		bankName = reqData.BankName
	}
	if reqData.BankType != "" {
		bankType = reqData.BankType
	}
	// if reqData.Balance

	acc.Account.Active = active
	acc.Account.Name = name
	acc.Account.Balance = balance
	acc.Account.Iban = iban
	acc.Account.BankCode = bankCode
	acc.Account.AccountNr = accountNr
	acc.Account.BankName = bankName
	acc.Account.BankType = bankType
	acc.Account.LastUpdate = time.Now()

	fmt.Println(acc.Account)

	var err error
	if acc.ID <= 0 {
		err = acc.Account.Create(db)
	} else {
		err = acc.Account.Save(db)
	}
	if err != nil {
		log.Println("[ERROR] API.updateAccount():", err)
		return
	}

	api.sendResult(w, acc.Account)
}
