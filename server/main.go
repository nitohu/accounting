package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/sessions"

	_ "github.com/lib/pq"
)

var tmpl *template.Template

var (
	key   = []byte("087736079f8d9e4c7fc7b642bb4c7afa")
	store = sessions.NewCookieStore(key)

	// datetime form layout
	dtFormLayout   = "Monday 02 January 2006 - 15:04"
	dateFormLayout = "Monday 02 January 2006"
	dbTimeLayout   = "2006-01-02 15:04:00"
)

func logInfo(funcName, msg string, args ...interface{}) {
	fmt.Printf("[INFO] %s %s: %s\n", time.Now().Local(), funcName, fmt.Sprintf(msg, args...))
}

func logWarn(funcName, msg string, args ...interface{}) {
	fmt.Printf("[WARN] %s %s: %s\n", time.Now().Local(), funcName, fmt.Sprintf(msg, args...))
}

func logError(funcName, msg string, args ...interface{}) {
	fmt.Printf("[ERROR] %s %s: %s\n", time.Now().Local(), funcName, fmt.Sprintf(msg, args...))
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s: %s\n", r.URL.Path, r.Method)

		f(w, r)
	}
}

func init() {
	db = dbInit("127.0.0.1", "nitohu", "123", "accounting", 5432)
	tmpl = template.Must(template.ParseGlob("./templates/*"))
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

	http.HandleFunc("/", logging(handleRoot))
	http.HandleFunc("/settings", logging(handleSettings))

	// Accounts
	http.HandleFunc("/accounts", logging(handleAccountOverview))
	http.HandleFunc("/accounts/form", logging(handleAccountForm))
	// TODO: Will be replaced with API once its done
	http.HandleFunc("/accounts/delete/{id}", logging(handleAccountDeletion))

	// Transactions
	http.HandleFunc("/transactions", logging(handleTransactionOverview))
	http.HandleFunc("/transactions/create", logging(handleTransactionForm))
	http.HandleFunc("/transactions/edit/{id}", logging(handleTransactionForm))
	http.HandleFunc("/transactions/delete/{id}", logging(handleTransactionDeletion))

	// Categories
	http.HandleFunc("/categories", logging(handleCategoryOverview))

	http.HandleFunc("/login", logging(handleLogin))

	http.HandleFunc("/logout", logging(handleLogout))

	http.ListenAndServe(":80", nil)
}
