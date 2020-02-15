package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/sessions"

	"github.com/gorilla/mux"
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
		t := time.Now()
		fmt.Printf("[INFO] %s %s: %s\n", t.Local(), r.URL.Path, r.Method)

		f(w, r)
	}
}

func main() {
	db = dbInit("127.0.0.1", "nitohu", "123", "accounting", 5432)
	defer db.Close()

	tmpl = template.Must(initTemplates(
		"index.html",
		"login.html",
		"accounts.html",
		"account_form.html",
		"transactions.html",
		"transaction_form.html",
		"404.html",
		"settings.html",
		"categories.html",
	))

	staticFiles := http.FileServer(http.Dir("static/"))

	r := mux.NewRouter()

	r.NotFoundHandler = pageNotFoundHandler()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFiles))

	r.HandleFunc("/", logging(handleRoot))
	r.HandleFunc("/settings/", logging(handleSettings))

	// Accounts
	r.HandleFunc("/accounts/", logging(handleAccountOverview))
	r.HandleFunc("/accounts/create/", logging(handleAccountCreation))
	r.HandleFunc("/accounts/edit/{id}", logging(handleAccountEditing))
	// TODO: Will be replaced with API once its done
	r.HandleFunc("/accounts/delete/{id}", logging(handleAccountDeletion))

	// Transactions
	r.HandleFunc("/transactions/", logging(handleTransactionOverview))
	r.HandleFunc("/transactions/create/", logging(handleTransactionForm))
	r.HandleFunc("/transactions/edit/{id}", logging(handleTransactionForm))
	r.HandleFunc("/transactions/delete/{id}", logging(handleTransactionDeletion))

	// Categories
	r.HandleFunc("/categories/", logging(handleCategoryOverview))

	r.HandleFunc("/login/", logging(handleLogin))

	r.HandleFunc("/logout/", logging(handleLogout))

	http.ListenAndServe(":80", r)
}
