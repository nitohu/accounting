package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nitohu/accounting/server/models"
)

func handleStatisticsOverview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/statistics/" {
		handleNotFound(w, r)
		return
	}
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)
	if !err.Empty() {
		err.AddTraceback("handleStatisticOverview()", "Error creating the context.")
		log.Println("[ERROR]", err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset: utf-8")

	ctx["Title"] = "Statistics"
	if ctx["Statistics"], err = models.GetAllStatistics(db); !err.Empty() {
		err.AddTraceback("handleStatisticOverview", "Error while getting statistics.")
		log.Println("[ERROR]", err)
	}

	fmt.Println(ctx["Statistics"].([]models.Statistic)[0].Name)
	fmt.Println(ctx["Statistics"].([]models.Statistic)[0].Value)

	if e := tmpl.ExecuteTemplate(w, "statistics.html", ctx); e != nil {
		log.Println("[ERROR] handleStatisticsOverview(): ", e)
	}
}
