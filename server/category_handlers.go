package main

import (
	"log"
	"net/http"
)

func handleCategoryOverview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/categories/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if !err.Empty() {
		err.AddTraceback("handleCategoryOverview()", "Error while creating the context.")
		log.Println("[ERROR]", err)
		http.Redirect(w, r, "/index/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; chartset=utf-8")
	ctx["Title"] = "Categories"

	e := tmpl.ExecuteTemplate(w, "categories.html", ctx)

	if e != nil {
		err.Init("handleCategoryOverview()", e.Error())
		log.Println("[ERROR]", err)
	}
}
