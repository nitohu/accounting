package models

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// User ...
type User struct {
	Active       bool
	ID           int
	Name         string
	Email        string
	TotalBalance float64
	CreateDate   time.Time
	LastUpdate   time.Time
}

// EmptyUser creates an empty user
func EmptyUser() User {
	t := time.Now().Local()
	res := User{false, 0, "", "", 0.0, t, t}

	return res
}

// CreateUser creates a user from a map
func CreateUser(user map[string]interface{}) User {
	var u User

	u.FromMap(user)

	return u
}

// Login verifies and logs a user in
func (u *User) Login(cr *sql.DB, email, password string) error {

	var user User
	var dbPass string

	query := "SELECT id, email, password FROM users WHERE email=$1"

	err := cr.QueryRow(query, email).Scan(&user.ID, &user.Email, &dbPass)

	if err != nil {
		return errors.New("[ERROR] No user with this email (" + email + ").\n")
	}

	fmt.Println("USERID: " + strconv.Itoa(user.ID))

	p := sha256.Sum256([]byte(password))
	pw := fmt.Sprintf("%X", p)

	if strings.Compare(pw, dbPass) != 0 {
		return errors.New("[ERROR] Wrong password.\n")
	}

	u.FindByID(cr, user.ID)

	fmt.Println("USERID: " + strconv.Itoa(u.ID))

	return (nil)
}

func (u *User) FindByID(cr *sql.DB, uid int) error {
	query := "SELECT active, id, name, email, total_balance, create_date, last_update FROM users WHERE id=$1"

	err := cr.QueryRow(query, uid).Scan(
		&u.Active,
		&u.ID,
		&u.Name,
		&u.Email,
		&u.TotalBalance,
		&u.CreateDate,
		&u.LastUpdate,
	)

	if err != nil {
		return err
	}

	return nil
}

// FindUserByID ...
func FindUserByID(cr *sql.DB, uid int) User {
	user := EmptyUser()

	err := user.FindByID(cr, uid)

	if err != nil {
		panic(err)
	}

	return user
}

// ToMap enables to save the User model into a session
func (u User) ToMap() map[string]interface{} {
	user := make(map[string]interface{})

	user["active"] = u.Active
	user["id"] = u.ID
	user["name"] = u.Name
	user["email"] = u.Email
	user["balance"] = u.TotalBalance
	user["createDate"] = u.CreateDate
	user["lastUpdate"] = u.LastUpdate

	return user
}

// FromMap takes a map and saves it into the user model
func (u User) FromMap(user map[string]interface{}) {

	u.Active = user["active"].(bool)
	u.ID = user["id"].(int)
	u.Name = user["name"].(string)
	u.Email = user["email"].(string)
	u.TotalBalance = user["balance"].(float64)
	u.CreateDate = user["createDate"].(time.Time)
	u.LastUpdate = user["LastUpdate"].(time.Time)

}
