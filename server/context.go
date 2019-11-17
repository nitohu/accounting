package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"./models"
	"github.com/gorilla/sessions"
)

// Context for backend/frontend communication
type Context struct {
	Authenticated bool
	User          models.User
}

// EmptyContext ...
func EmptyContext() Context {
	ctx := Context{false, models.EmptyUser()}

	return ctx
}

func createContextFromSession(cr *sql.DB, session *sessions.Session) (Context, error) {
	ctx := Context{false, models.EmptyUser()}

	authenticated, ok := session.Values["authenticated"].(bool)

	if !ok {
		err := fmt.Sprintf("[ERROR] %s CreateContextFromSession(): Authenticated not present in session", time.Now().Local())
		return EmptyContext(), errors.New(err)
	}

	uid, ok := session.Values["uid"].(int)

	if !ok {
		err := fmt.Sprintf("[ERROR] %s CreateContextFromSession(): UserID not present in session", time.Now().Local())
		return EmptyContext(), errors.New(err)
	}

	ctx.Authenticated = authenticated
	ctx.User = models.FindUserByID(cr, uid)

	return ctx, nil
}
