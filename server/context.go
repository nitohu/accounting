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

	authenticated, ok := session.Values["authenticated"].(bool)

	if !ok {
		err := "Authenticated not present in session"
		return EmptyContext(), errors.New(err)
	}

	uid, ok := session.Values["uid"].(int)

	if !ok {
		err := "UserID not present in session"
		return EmptyContext(), errors.New(err)
	}

	settings, err := models.InitializeSettings(cr)

	if err != nil {
		return ctx, err
	}

	ctx["HumanReadable"] = HumanReadable
	ctx["Authenticated"] = authenticated
	ctx["User"], err = models.FindUserByID(cr, uid)
	ctx["Settings"] = settings

	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
