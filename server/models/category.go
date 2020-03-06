package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
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
func (c *Category) Create(cr *sql.DB) error {
	if c.ID > 0 {
		err := "This category already has an ID. Maybe try saving it?"
		return errors.New(err)
	} else if c.Name == "" {
		err := "This category has no name. Please set a name before creating a category."
		return errors.New(err)
	}

	query := "INSERT INTO categories (name, create_date, last_update, active, hex) "
	query += "VALUES ($1, $2, $3, $4, $5) RETURNING id;"

	if c.Hex == "" {
		c.Hex = "#ffffff"
	}

	err := cr.QueryRow(query,
		c.Name,
		time.Now(),
		time.Now(),
		true,
		c.Hex,
	).Scan(&c.ID)

	if err != nil {
		fmt.Println("Category.Create(): Traceback: Error while executing query")
		return err
	}

	c.computeFields(cr)

	return nil
}

// Save the current category to the database
func (c *Category) Save(cr *sql.DB) error {
	if c.ID <= 0 {
		err := "This category has no ID. Maybe create it first?"
		return errors.New(err)
	} else if c.Name == "" {
		err := "This category has no name. Please set a name before saving it."
		return errors.New(err)
	}

	query := "UPDATE categories SET name=$1, hex=$2, last_update=$3, active=$4 WHERE id=$5"

	if c.Hex == "" {
		c.Hex = "#ffffff"
	}

	_, err := cr.Exec(query,
		c.Name,
		c.Hex,
		time.Now(),
		c.Active,
		c.ID,
	)

	if err != nil {
		fmt.Println("Category.Save(): Traceback: Error while saving the category to the database.")
		return err
	}

	c.computeFields(cr)

	return nil
}

// Delete the current category from the database
func (c *Category) Delete(cr *sql.DB) error {
	if c.ID <= 0 {
		return errors.New("Category.Delete(): ID must be bigger than 0")
	}
	query := "DELETE FROM categories WHERE id=$1"

	if _, err := cr.Exec(query, c.ID); err != nil {
		log.Println("[ERROR] Category.Delete():", err)
	}

	c.ID = 0
	c.Name = ""
	c.Hex = ""
	c.Active = false

	return nil
}

func (c *Category) computeFields(cr *sql.DB) error {
	transQuery := "SELECT id FROM transactions where category_id=$1;"
	res, err := cr.Query(transQuery, c.ID)
	if err != nil {
		fmt.Println("[ERROR] Category.Save(): Error executing query f")
	}

	c.TransactionIDs = nil

	for res.Next() {
		var id int64
		if err = res.Scan(&id); err != nil {
			fmt.Println("[WARN] Category.computeFields(): Could not scan ID, skipping row")
			continue
		}
		c.TransactionIDs = append(c.TransactionIDs, id)
	}

	c.TransactionCount = len(c.TransactionIDs)

	return nil
}

// FindByID finds a category in the database
func (c *Category) FindByID(cr *sql.DB, id int64) error {
	if id <= 0 {
		err := "ID must be a positive number: " + string(id)
		return errors.New(err)
	}

	query := "SELECT id,name,create_date,last_update,active,hex FROM categories WHERE id=$1;"

	err := cr.QueryRow(query, id).Scan(
		&c.ID,
		&c.Name,
		&c.CreateDate,
		&c.LastUpdate,
		&c.Active,
		&c.Hex,
	)

	if err != nil {
		fmt.Println("Category.FindByID(): Traceback: Error executing SELECT query.")
		return err
	}

	c.computeFields(cr)

	return nil
}

// FindCategoryByID is similar to FindByID but returns the category
func FindCategoryByID(cr *sql.DB, categoryID int64) (Category, error) {
	t := EmptyCategory()

	err := t.FindByID(cr, categoryID)

	if err != nil {
		return t, err
	}

	return t, nil
}


// GetAllCategories returns all categories
func GetAllCategories(cr *sql.DB) ([]Category, error) {
	var categories []Category
	query := "SELECT id FROM categories;"

	res, err := cr.Query(query)

	if err != nil {
		fmt.Println("GetAllCategories(): Traceback: Error while executing the query")
		return nil, err
	}

	for res.Next() {
		var id int64

		if err = res.Scan(&id); err != nil {
			fmt.Println("[INFO] Skipping Record")
			fmt.Printf("[WARN] GetAllCategories(): Error while scanning the id:\n%s\n", err)
		} else {
			cat := EmptyCategory()

			if err = cat.FindByID(cr, id); err != nil {
				fmt.Println("[INFO] Skipping Record")
				fmt.Printf("[WARN] GetAllCategories: %s\n", err)
			} else {
				cat.computeFields(cr)
				categories = append(categories, cat)
			}
		}
	}

	return categories, nil
}
