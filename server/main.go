package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"./models"
	"github.com/gorilla/sessions"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var tmpl *template.Template
var db *sql.DB

var (
	key   = []byte("087736079f8d9e4c7fc7b642bb4c7afa")
	store = sessions.NewCookieStore(key)
)

func logInfo(funcName, msg string, args ...interface{}) {
	fmt.Printf("[INFO] %s %s: %s\n", time.Now().Local(), funcName, fmt.Sprintf(msg, args...))
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

	tmpl.ExecuteTemplate(w, "index.html", ctx)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		logInfo("handleLogin()", "User is already logged in", r.URL.Path)
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
	db = initDb("127.0.0.1", "nitohu", "123", "accounting", 5432)
	defer db.Close()

	// user := user.Query(db, []user.QueryArgument{
	// 	{Connector: "", Field: "name", Op: "ilike", Value: "niklas%"},
	// 	{Connector: "AND", Field: "id", Op: "=", Value: "1"},
	// })

	tmpl = template.Must(initTemplates("index.html", "login.html"))

	staticFiles := http.FileServer(http.Dir("static/"))

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFiles))

	r.HandleFunc("/", logging(handleRoot))

	r.HandleFunc("/login/", logging(handleLogin))

	r.HandleFunc("/logout/", logging(handleLogout))

	http.ListenAndServe(":80", r)
}
