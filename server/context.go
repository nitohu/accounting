package main

import (
	"database/sql"

	"github.com/nitohu/err"

	"github.com/nitohu/accounting/server/models"
	"github.com/gorilla/sessions"
)

// Context type represents context
type Context map[string]interface{}

// EmptyContext ...
func EmptyContext() Context {
	ctx := make(map[string]interface{})

	return ctx
}

func createContextFromSession(cr *sql.DB, session *sessions.Session) (Context, err.Error) {
	ctx := EmptyContext()

	// Make sure authenticated flag in session is true and exists
	var authenticated bool
	if authenticated, ok := session.Values["authenticated"].(bool); !authenticated || !ok {
		var err err.Error
		err.Init("createContextFromSession()", "Authenticated not present in session or false.")
		return ctx, err
	}

	// Check the session key with the one from the database
	var dbSessionKey string
	query := "SELECT session_key FROM settings LIMIT 1;"

	if e := db.QueryRow(query).Scan(&dbSessionKey); e != nil {
		var err err.Error
		err.Init("createContextFromSession()", e.Error())
		return ctx, err
	}

	if sessionKey, ok := session.Values["key"].(string); !ok || (sessionKey != dbSessionKey) {
		var err err.Error
		err.Init("createContextFromSession()", "Error validating session key with the one from db.")
		return ctx, err
	}

	// Get the settings from the database
	settings, e := models.InitializeSettings(cr)

	if !e.Empty() {
		e.AddTraceback("createContextFromSession", "Error initializing settings.")
		return ctx, e
	}

	// Add the variables to the context
	ctx["HumanReadable"] = HumanReadable
	ctx["Authenticated"] = authenticated
	ctx["Settings"] = settings

	return ctx, err.Error{}
}
