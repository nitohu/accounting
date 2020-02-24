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

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/api")[1]
	var body []byte

	w.Header().Set("Content-Type", "application/json; charset: utf-8")

	// Parse a body if there is one
	if r.ContentLength > 0 {
		b := r.Body
		body = make([]byte, r.ContentLength)
		if _, err := io.ReadFull(b, body); err != nil {
			fmt.Println("[ERROR] API.ServeHTTP():", err)
		}
	}

	api.obj = nil
	api.id = 0

	switch {
	case path == "/categories" || path == "/categories/":
		// Check if the db attribute is appended to the url
		// If yes, return only that one record
		api.getIDFromURL(r)
		if api.id > 0 {
			api.getCategoryByID(w, r)
			return
		}
		api.getCategories(w, r)
	case path == "/categories/create" || path == "/categories/create/":
		c := models.EmptyCategory()
		if err := json.Unmarshal(body, &c); err != nil {
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		c.ID = 0
		api.obj = c
		api.updateCategory(w, r)
	case path == "/categories/update":
		c := models.EmptyCategory()
		if err := json.Unmarshal(body, &c); err != nil {
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		api.obj = c
		// Validate if the ID is existing
		if err := c.FindByID(db, c.ID); err != nil {
			fmt.Fprintf(w, "{'error': '%s'}", err)
			return
		}
		api.id = c.ID
		api.updateCategory(w, r)
	}
}

func (api *API) getIDFromURL(r *http.Request) {
	api.id = 0
	idStr := r.FormValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("[ERROR] API.getIDFromURL():", err)
			return
		}
		api.id = int64(id)
	}
}

/*
	##############################
	#                            #
	#         Categories         #
	#                            #
	##############################
*/

// Getters

// getCategories gets all Categories
// Also handles if
func (api API) getCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/categories/: Method must be GET.'}")
		return
	}

	c, err := models.GetAllCategories(db)
	if err != nil {
		log.Printf("[ERROR] API.getCategories(): Error getting Categories.\n%s\n", err)
		return
	}
	d, err := json.Marshal(c)
	if err != nil {
		log.Printf("[ERROR] API.getCategories(): Error parsing JSON.\n%s\n", err)
		return
	}
	data := string(d)
	fmt.Fprint(w, data)
}

// getCategoryByID gets a category by it's ID
func (api API) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "{'error': '/categories/: Method must be GET.'}")
		return
	}
	if api.id <= 0 {
		fmt.Fprintln(w, "{'error': 'Please provide a valid ID.'}")
		return
	}

	c := models.EmptyCategory()

	if err := c.FindByID(db, int64(api.id)); err != nil {
		log.Println("API.getCategoryByID():", err)
		fmt.Fprintln(w, "{'error': 'There was an unexpected error while searching the record by ID.'}")
		return
	}

	d, err := json.Marshal(c)
	if err != nil {
		log.Println("API.getCategoryByID():", err)
		fmt.Fprintln(w, "{'error': 'There was an unexpected error while converting the data to JSON.'}")
		return
	}
	data := string(d)
	fmt.Fprintln(w, data)
}

// Setter

// getCategoryByID write to a category
// if this category does not exist or has no ID, create one
// Method MUST be POST
// Returns the data written to the database
func (api API) updateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, "{'error': '/categories/create: Method must be POST.'}")
		return
	}

	c := models.EmptyCategory()

	// Get data if ID is present
	if api.id > 0 {
		if err := c.FindByID(db, api.id); err != nil {
			log.Println("[WARN] API.updateCategory():", err)
			fmt.Fprint(w, "{'error': 'Error getting your ID.'}")
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

	fmt.Println("Name:", reqData.Name)
	fmt.Println("Hex:", reqData.Hex)

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

	// Create return data
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("[WARN] API.updateCategory()", err)
		fmt.Fprint(w, "{'error': 'Error converting object to JSON.'}")
		return
	}

	// Write data to the client
	fmt.Fprint(w, string(data))
}
