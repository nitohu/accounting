package main

import (
	"net/http"
)

func handleCategoryOverview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/categories/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logWarn("handleCategoryOverview", "Error while creating the context:\n%s", err)
		http.Redirect(w, r, "/index/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; chartset=utf-8")
	ctx["Title"] = "Categories"

	err = tmpl.ExecuteTemplate(w, "categories.html", ctx)

	if err != nil {
		logError("handleCategoryOverview", "%s", err)
	}
}
