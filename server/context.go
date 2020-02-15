package main

import (
	"database/sql"
	"errors"

	"./models"
	"github.com/gorilla/sessions"
)

// Context type represents context
type Context map[string]interface{}

// EmptyContext ...
func EmptyContext() Context {
	ctx := make(map[string]interface{})

	return ctx
}

func createContextFromSession(cr *sql.DB, session *sessions.Session) (Context, error) {
	ctx := EmptyContext()

	// Make sure authenticated flag in session is true
	authenticated, ok := session.Values["authenticated"].(bool)

	if !ok {
		err := "Authenticated not present in session"
		return EmptyContext(), errors.New(err)
	}

	// Check the session key with the one from the database
	var dbSessionKey string
	query := "SELECT session_key FROM settings LIMIT 1;"

	err := db.QueryRow(query).Scan(&dbSessionKey)

	if err != nil {
		logError("createContextFromSession", "Traceback #1")
		return EmptyContext(), err
	}

	sessionKey, ok := session.Values["key"].(string)

	if !ok {
		logWarn("createContextFromSession", "Traceback: Getting session key from session")
		err := "Session key not present in session"
		return EmptyContext(), errors.New(err)
	}

	if sessionKey != dbSessionKey {
		err := "Session key is not equivalent to the one in the database."
		return EmptyContext(), errors.New(err)
	}

	// Get the settings from the database
	settings, err := models.InitializeSettings(cr)

	if err != nil {
		logWarn("createContextFromSession", "Traceback: Initializing settings")
		return ctx, err
	}

	// Add the variables to the context
	ctx["HumanReadable"] = HumanReadable
	ctx["Authenticated"] = authenticated
	ctx["Settings"] = settings

	return ctx, nil
}
