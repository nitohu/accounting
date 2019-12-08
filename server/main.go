package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"./models"
	"github.com/gorilla/sessions"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var tmpl *template.Template

var (
	key   = []byte("087736079f8d9e4c7fc7b642bb4c7afa")
	store = sessions.NewCookieStore(key)
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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

	err = tmpl.ExecuteTemplate(w, "index.html", ctx)

	if err != nil {
		logError("handleLogin", "%s", err)
	}
}

func handleAccountOverview(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountOverview", "%s", err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

	ctx["Title"] = "All Accounts"

	err = tmpl.ExecuteTemplate(w, "accounts.html", ctx)

	if err != nil {
		logError("handleAccountOverview", "%s", err)
	}
}

func handleAccountCreation(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountCreation", "%s", err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

	ctx["Title"] = "Create Account"
	ctx["Header"] = "Create a new account"
	ctx["Btn"] = "Create Account"
	ctx["Account"] = models.EmptyAccount()

	if r.Method != http.MethodPost {
		err = tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		if err != nil {
			logError("handleLogin", "%s", err)
		}
		return
	}

	account := models.EmptyAccount()

	account.Name = r.FormValue("name")
	account.Balance, _ = strconv.ParseFloat(r.FormValue("balance"), 64)
	account.BalanceForecast = account.Balance
	account.Iban = r.FormValue("iban")
	account.BankCode = r.FormValue("bankCode")
	account.AccountNr = r.FormValue("accountNumber")
	account.BankName = r.FormValue("bankName")
	account.BankType = r.FormValue("accountType")
	account.UserID = ctx["User"].(models.User).ID

	err = account.Create(db)

	if err != nil {
		logError("handleAccountCreation", "%s", err)
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func handleAccountEditing(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	ctx, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountEditing", "%s", err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])

	account := models.EmptyAccount()
	account.FindByID(db, accountID)

	ctx["Title"] = "Edit Account: " + account.Name
	ctx["Header"] = "Edit " + account.Name
	ctx["Btn"] = "Save Account"
	ctx["Account"] = account

	if r.Method != http.MethodPost {
		err := tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		if err != nil {
			logError("handleAccountEditing", "%s", err)
		}
		return
	}

	account.Name = r.FormValue("name")
	account.Balance, _ = strconv.ParseFloat(r.FormValue("balance"), 64)
	account.BalanceForecast = account.Balance
	account.Iban = r.FormValue("iban")
	account.BankCode = r.FormValue("bankCode")
	account.AccountNr = r.FormValue("accountNumber")
	account.BankType = r.FormValue("accountType")

	if account.BankType == "bank" {
		account.BankName = r.FormValue("bankName")
	} else {
		account.BankName = r.FormValue("providerName")
	}

	if account.BankType != "bank" && account.BankType != "online" {
		logError("handleAccountEditing", "Wrong value for bank type: %s", account.BankType)
		ctx["Error"] = "There was an error in the form."
		err = tmpl.ExecuteTemplate(w, "account_form.html", ctx)
		return
	}

	err = account.Save(db)

	if err != nil {
		logError("handleAccountEditing", "%s", err)
		http.Redirect(w, r, "/accounts/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/accounts/", http.StatusSeeOther)
}

func handleAccountDeletion(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	_, err := createContextFromSession(db, session)

	if err != nil {
		logError("handleAccountDeletion", "%s", err)
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		logError("handleAccountDeletion", "%s", err)
		return
	}

	account := models.FindAccountByID(db, id)

	account.Delete(db)

}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		logInfo(r.URL.Path, "User is already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	u := models.EmptyUser()

	mail := r.FormValue("email")
	pw := r.FormValue("password")

	logInfo("handleLogin()", "Login try %s", mail)

	err := u.Login(db, mail, pw)

	if err != nil {
		fmt.Printf("[ERROR] %s %s", time.Now().Local(), err)
		tmpl.ExecuteTemplate(w, "login.html", map[string]string{"err": "Wrong credentials."})
		return
	}

	logInfo("handleLogin()", "Login try of %s was successful", mail)

	logInfo("handleLogin()", "UID: %d\n", u.ID)

	// Successfully logged in
	session.Values["authenticated"] = true
	session.Values["email"] = u.Email
	session.Values["uid"] = u.ID

	err = session.Save(r, w)

	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	session.Values["authenticated"] = false
	session.Values["user"] = models.EmptyUser()

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	db = dbInit("127.0.0.1", "nitohu", "123", "accounting", 5432)
	defer db.Close()

	// user := user.Query(db, []user.QueryArgument{
	// 	{Connector: "", Field: "name", Op: "ilike", Value: "niklas%"},
	// 	{Connector: "AND", Field: "id", Op: "=", Value: "1"},
	// })

	tmpl = template.Must(initTemplates("index.html", "login.html", "accounts.html", "account_form.html"))

	staticFiles := http.FileServer(http.Dir("static/"))

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFiles))

	r.HandleFunc("/", logging(handleRoot))

	r.HandleFunc("/accounts/", logging(handleAccountOverview))
	r.HandleFunc("/accounts/create/", logging(handleAccountCreation))
	r.HandleFunc("/accounts/edit/{id}", logging(handleAccountEditing))
	// TODO: Will be replaced with API once its done
	r.HandleFunc("/accounts/delete/{id}", logging(handleAccountDeletion))

	r.HandleFunc("/login/", logging(handleLogin))

	r.HandleFunc("/logout/", logging(handleLogout))

	http.ListenAndServe(":80", r)
}
