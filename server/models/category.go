package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/nitohu/err"
)

// Category ...
type Category struct {
	ID         int64
	Name       string
	Hex        string
	CreateDate time.Time
	LastUpdate time.Time
	Active     bool

	// Computed fields
	TransactionIDs   []int64
	TransactionCount int
}

// EmptyCategory returns an empty category
func EmptyCategory() Category {
	cat := Category{
		ID:         0,
		Name:       "",
		Hex:        "",
		CreateDate: time.Now(),
		LastUpdate: time.Now(),
		Active:     false,
	}

	return cat
}

// Create a new category in the database
func (c *Category) Create(cr *sql.DB) err.Error {
	if c.ID > 0 {
		var err err.Error
		err.Init("Category.Create()", "This category already has an ID. Maybe try saving it?")
		return err
	} else if c.Name == "" {
		var err err.Error
		err.Init("Category.Create()", "This category has no name. Please set a name before creating a category.")
		return err
	}

	query := "INSERT INTO categories (name, create_date, last_update, active, hex) "
	query += "VALUES ($1, $2, $3, $4, $5) RETURNING id;"

	if c.Hex == "" {
		c.Hex = "#ffffff"
	}

	e := cr.QueryRow(query,
		c.Name,
		time.Now(),
		time.Now(),
		true,
		c.Hex,
	).Scan(&c.ID)

	if e != nil {
		var err err.Error
		err.Init("Category.Create()", e.Error())
		return err
	}

	c.computeFields(cr)

	return err.Error{}
}

// Save the current category to the database
func (c *Category) Save(cr *sql.DB) err.Error {
	if c.ID <= 0 {
		var err err.Error
		err.Init("Category.Save()", "This category has no ID. Maybe create it first?")
		return err
	} else if c.Name == "" {
		var err err.Error
		err.Init("Category.Save()", "This category has no name. Please set a name before saving it.")
		return err
	}

	query := "UPDATE categories SET name=$1, hex=$2, last_update=$3, active=$4 WHERE id=$5"

	if c.Hex == "" {
		c.Hex = "#ffffff"
	}

	_, e := cr.Exec(query,
		c.Name,
		c.Hex,
		time.Now(),
		c.Active,
		c.ID,
	)

	if e != nil {
		var err err.Error
		err.Init("Category.Save()", e.Error())
		return err
	}

	c.computeFields(cr)

	return err.Error{}
}

// Delete the current category from the database
func (c *Category) Delete(cr *sql.DB) err.Error {
	if c.ID <= 0 {
		var err err.Error
		err.Init("Category.Delete()", "ID must be bigger than 0")
		return err
	}
	query := "DELETE FROM categories WHERE id=$1"

	if _, e := cr.Exec(query, c.ID); e != nil {
		var err err.Error
		err.Init("Category.Delete()", e.Error())
		return err
	}

	c.ID = 0
	c.Name = ""
	c.Hex = ""
	c.Active = false

	return err.Error{}
}

func (c *Category) computeFields(cr *sql.DB) {
	transQuery := "SELECT id FROM transactions where category_id=$1;"
	res, e := cr.Query(transQuery, c.ID)
	if e != nil {
		log.Printf("[ERROR] Category.computeFields(): Error getting transaction IDs\n%s\n", e)
	}

	c.TransactionIDs = nil

	for res.Next() {
		var id int64
		if e = res.Scan(&id); e != nil {
			log.Println("[WARN] Category.computeFields(): Could not scan ID, skipping row")
			continue
		}
		c.TransactionIDs = append(c.TransactionIDs, id)
	}

	c.TransactionCount = len(c.TransactionIDs)
}

// FindByID finds a category in the database
func (c *Category) FindByID(cr *sql.DB, id int64) err.Error {
	if id <= 0 {
		var err err.Error
		err.Init("Category.FindByID()", "ID must be a positive number: "+string(id))
		return err
	}

	query := "SELECT id,name,create_date,last_update,active,hex FROM categories WHERE id=$1;"

	e := cr.QueryRow(query, id).Scan(
		&c.ID,
		&c.Name,
		&c.CreateDate,
		&c.LastUpdate,
		&c.Active,
		&c.Hex,
	)

	if e != nil {
		var err err.Error
		err.Init("Category.FindByID()", e.Error())
		return err
	}

	c.computeFields(cr)

	return err.Error{}
}

// FindCategoryByID is similar to FindByID but returns the category
func FindCategoryByID(cr *sql.DB, categoryID int64) (Category, err.Error) {
	t := EmptyCategory()

	e := t.FindByID(cr, categoryID)

	if !e.Empty() {
		var err err.Error
		err.Init("FindCategoryByID()", e.Error())
		return t, err
	}

	return t, err.Error{}
}

// GetAllCategories returns all categories
func GetAllCategories(cr *sql.DB) ([]Category, err.Error) {
	var categories []Category
	query := "SELECT id FROM categories;"

	res, e := cr.Query(query)

	if e != nil {
		var err err.Error
		err.Init("GetAllCategories()", "Error while executing the query")
		return nil, err
	}

	for res.Next() {
		var id int64

		if e = res.Scan(&id); e != nil {
			log.Println("[INFO] GetAllCategories(): Skipping Record")
			log.Printf("[WARN] GetAllCategories(): Error while scanning the id:\n%s\n", e)
		} else {
			cat := EmptyCategory()

			if e = cat.FindByID(cr, id); e != nil {
				log.Println("[INFO] GetAllCategories(): Skipping Record")
				log.Printf("[WARN] GetAllCategories(): %s\n", e)
			} else {
				cat.computeFields(cr)
				categories = append(categories, cat)
			}
		}
	}

	return categories, err.Error{}
}
