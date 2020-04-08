package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/sessions"

	_ "github.com/lib/pq"
)

var tmpl *template.Template

var (
	port = ":80"

	key   = []byte("087736079f8d9e4c7fc7b642bb4c7afa")
	store = sessions.NewCookieStore(key)

	// datetime form layout
	dateSettingsLayout = "Monday 02 January 2006"
	dtLayout           = "02.01.2006 - 15:04"
	dateLayout         = "02.01.2006"
	dbTimeLayout       = "2006-01-02 15:04:00"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s: %s\n", r.URL.Path, r.Method)

		f(w, r)
	}
}

func init() {
	data, err := getCmdLineArgs(os.Args)
	if !err.Empty() {
		log.Fatalln(err)
	}
	db = dbInit(data["dbhost"], data["dbuser"], data["dbpassword"], data["dbdatabase"], data["dbport"])

	tmpl = template.Must(template.ParseGlob("./templates/*"))

	if val, ok := data["port"]; ok {
		port = ":" + val
	}
}

func main() {
	defer db.Close()
	http.Handle(
		"/static/", http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static/")),
		),
	)

	var api API

	http.Handle("/api/", api)

	// General
	http.HandleFunc("/", logging(handleRoot))
	http.HandleFunc("/settings/", logging(handleSettings))
	http.HandleFunc("/login/", logging(handleLogin))
	http.HandleFunc("/logout/", logging(handleLogout))

	// Accounts
	http.HandleFunc("/accounts/", logging(handleAccountOverview))
	http.HandleFunc("/accounts/form/", logging(handleAccountForm))

	// Transactions
	http.HandleFunc("/transactions/", logging(handleTransactionOverview))
	http.HandleFunc("/transactions/form/", logging(handleTransactionForm))
	http.HandleFunc("/transactions/delete/{id}/", logging(handleTransactionDeletion))

	// Statistics
	http.HandleFunc("/statistics/", logging(handleStatisticsOverview))

	// Categories
	http.HandleFunc("/categories/", logging(handleCategoryOverview))

	log.Fatalln(http.ListenAndServe(port, nil))
}
