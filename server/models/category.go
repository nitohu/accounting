package models

import (
	"database/sql"
	"errors"
	"fmt"
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
	TransactionIDs []int64
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
	if c.ID >= 0 {
		err := "This category already has an ID. Maybe try saving it?"
		return errors.New(err)
	} else if c.Name == "" {
		err := "This category has no name. Please set a name before creating a category."
		return errors.New(err)
	}

	query := "INSERT INTO categories (name, create_date, last_update, active, hex) "
	query += "VALUES ($1, $2, $2, $4, $5);"

	if c.Hex == "" {
		c.Hex = "#ffffff"
	}

	res, err := cr.Exec(query,
		c.Name,
		time.Now(),
		time.Now(),
		true,
		c.Hex,
	)

	if err != nil {
		fmt.Println("Category.Create(): Traceback: Error while executing query")
		return err
	}

	if c.ID, err = res.LastInsertId(); err != nil {
		fmt.Println("Category.Create(): Traceback: Error while getting last insert ID")
		return err
	}

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

	return nil
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
				categories = append(categories, cat)
			}
		}
	}

	return categories, nil
}
