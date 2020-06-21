package models

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nitohu/err"
)

// API is the database interface to the api table
type API struct {
	ID           int
	Active       bool
	Name         string
	CreateDate   time.Time
	LastUpdate   time.Time
	apiKey       string
	APIPrefix    string
	AccessRights []string
	LocalKey     bool
}

//GetAllAccessRights returns a list of all existing api access rights
func GetAllAccessRights() []string {
	return []string{
		"transaction.read",
		"transaction.write",
		"transaction.delete",
		"account.read",
		"account.write",
		"account.delete",
		"statistic.read",
		"statistic.write",
		"statistic.delete",
		"settings.read",
		"settings.write",
		"settings.delete",
		"category.read",
		"category.write",
		"category.delete",
	}
}

func formatAccessRights(rights []string) string {
	res := ""

	for i := 0; i < len(rights); i++ {
		r := rights[i]
		res += r
		if i < (len(rights) - 1) {
			res += ";"
		}
	}

	return res
}

func generateRandomKey(length int) string {
	key := ""
	seed := time.Now().Nanosecond()
	s := rand.NewSource(int64(seed))
	r := rand.New(s)

	for i := 0; i < length; i++ {
		seed = time.Now().Nanosecond()
		r.Seed(int64(seed))
		code := int(r.Float64()*57 + 65)
		key += string(code)
	}

	return key
}

// Create the current instance in the database
func (a *API) Create(cr *sql.DB) err.Error {
	if a.ID != 0 {
		var e err.Error
		e.Init("API.Create()", "This object already has an ID.")
		return e
	}
	if a.apiKey == "" || a.APIPrefix == "" {
		var e err.Error
		e.Init("API.Create()", "This object does not have an apiKey or APIPrefix.")
		return e
	}

	query := "INSERT INTO api (active, name, create_date, last_update, api_key, api_prefix, local_key, access_rights) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"

	a.CreateDate = time.Now()
	a.LastUpdate = time.Now()
	rights := formatAccessRights(a.AccessRights)

	_, e := cr.Exec(query,
		a.Active,
		a.Name,
		a.CreateDate,
		a.LastUpdate,
		a.apiKey,
		a.APIPrefix,
		a.LocalKey,
		rights,
	)
	if e != nil {
		var err err.Error
		err.Init("API.Create()", e.Error())
		return err
	}

	return err.Error{}
}

// Save the current instance in the database
func (a *API) Save(cr *sql.DB) err.Error {
	if a.ID >= 0 {
		var e err.Error
		e.Init("API.Save()", "ID must be bigger than 0. ("+string(a.ID)+")")
		return e
	}
	if a.apiKey == "" || a.APIPrefix == "" {
		var e err.Error
		e.Init("API.Save()", "This object does not have an apiKey or APIPrefix.")
		return e
	}

	query := "UPDATE api SET active=$2, name=$3, last_update=$4, api_key=5, api_prefix=6, local_key=$7, access_rights=$8 WHERE id=$1;"

	a.LastUpdate = time.Now()
	rights := formatAccessRights(a.AccessRights)

	_, e := cr.Exec(query,
		a.ID,
		a.Active,
		a.Name,
		a.LastUpdate,
		a.apiKey,
		a.APIPrefix,
		a.LocalKey,
		rights,
	)
	if e != nil {
		var err err.Error
		err.Init("API.Save()", e.Error())
		return err
	}

	return err.Error{}
}

// SetAPIKey hashes the given key and sets the value of a.apiKey to hashed(key)
func (a *API) SetAPIKey(key string) err.Error {
	if key == "" {
		var e err.Error
		e.Init("API.SetAPIKey()", "Given key is empty.")
		return e
	}
	if len(key) != 32 {
		var e err.Error
		e.Init("API.SetAPIKey()", "API Key must be 32 characters long.")
		return e
	}

	bhash := sha256.Sum256([]byte(key))
	hash := fmt.Sprintf("%x", bhash)

	a.apiKey = hash

	return err.Error{}
}

// GenerateAPIKey generates an api key and sets it to the variable
func (a *API) GenerateAPIKey() {
	a.APIPrefix = generateRandomKey(6)
	key := generateRandomKey(32)
	fullKey := a.APIPrefix + "." + key
	if a.LocalKey {
		a.apiKey = fullKey
	} else {
		bhash := sha256.Sum256([]byte(fullKey))
		a.apiKey = fmt.Sprintf("%x", bhash)
	}
}

// FindByPrefix takes the given prefix and returns the corresponding API record
func (a *API) FindByPrefix(cr *sql.DB, prefix string) err.Error {
	query := "SELECT id,active,name,create_date,last_update,api_key,api_prefix,local_key,access_rights"
	query += " FROM api WHERE api_prefix=$1;"

	var rights string

	e := cr.QueryRow(query, prefix).Scan(
		&a.ID,
		&a.Active,
		&a.Name,
		&a.CreateDate,
		&a.LastUpdate,
		&a.apiKey,
		&a.APIPrefix,
		&a.LocalKey,
		&rights,
	)
	if e != nil {
		var err err.Error
		err.Init("API.FindByPrefix()", e.Error())
		return err
	}
	a.AccessRights = strings.Split(rights, ";")

	return err.Error{}
}

// FindByID takes the given ID and returns the corresponding API record
func (a *API) FindByID(cr *sql.DB, id int64) err.Error {
	query := "SELECT id,active,name,create_date,last_update,api_key,api_prefix,local_key,access_rights"
	query += " FROM api WHERE id=$1;"

	var rights string

	e := cr.QueryRow(query, id).Scan(
		&a.ID,
		&a.Active,
		&a.Name,
		&a.CreateDate,
		&a.LastUpdate,
		&a.apiKey,
		&a.APIPrefix,
		&a.LocalKey,
		&rights,
	)
	if e != nil {
		var err err.Error
		err.Init("API.FindByID()", e.Error())
		return err
	}
	a.AccessRights = strings.Split(rights, ";")

	return err.Error{}
}

// GetLocalAPIKeys returns all local API Keys
func GetLocalAPIKeys(cr *sql.DB) ([]API, err.Error) {
	var res []API

	query := "SELECT api_prefix FROM api WHERE local_key='t';"
	rows, e := cr.Query(query)
	if e != nil {
		var err err.Error
		err.Init("API.GetLocalKey()", e.Error())
		return nil, err
	}
	for rows.Next() {
		prefix := ""
		if e = rows.Scan(&prefix); e != nil {
			var err err.Error
			err.Init("API.GetLocalKey()", e.Error())
			return nil, err
		}
		var a API
		if err := a.FindByPrefix(cr, prefix); !err.Empty() {
			return nil, err
		}
		res = append(res, a)
	}

	return res, err.Error{}
}

// GetAllAPIKeys returns all existing API keys
func GetAllAPIKeys(cr *sql.DB) ([]API, err.Error) {
	var res []API

	query := "SELECT id FROM api;"
	rows, e := cr.Query(query)
	if e != nil {
		var err err.Error
		err.Init("API.GetAllAPIKeys()", e.Error())
		return nil, err
	}
	for rows.Next() {
		var id int64
		if e = rows.Scan(&id); e != nil {
			var err err.Error
			err.Init("API.GetAllAPIKeys()", e.Error())
			return nil, err
		}
		var a API
		if err := a.FindByID(cr, id); !err.Empty() {
			return nil, err
		}
		res = append(res, a)
	}

	return res, err.Error{}
}
