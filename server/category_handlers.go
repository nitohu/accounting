package main

import (
	"net/http"

	"./models"
)

func handleCategoryOverview(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logWarn("handleCategoryOverview", "Error while creating the context:\n%s", err)
		http.Redirect(w, r, "/index/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; chartset=utf-8")

	var categories []models.Category

	if categories, err = models.GetAllCategories(db); err != nil {
		logWarn("handleCategoryOverview", "Error while getting all categories:\n%s", err)
	}

	ctx["Categories"] = categories
	ctx["CategoriesLen"] = len(categories)

	err = tmpl.ExecuteTemplate(w, "categories.html", ctx)

	if err != nil {
		logError("handleCategoryOverview", "%s", err)
	}
}
