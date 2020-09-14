package wonderwall

import (
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"github.com/rs/xid"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	applicationJson = "application/json"
)

func ContentType(t string, w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Content-Type") != t {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return false
	} else {
		return true
	}
}

func writeTmpl(w http.ResponseWriter, name string, i interface{}) {
	indexTmpl, err := template.ParseFiles(fmt.Sprintf("./template/%s.html", name))
	if err != nil {
		panic(err)
	}
	if err = indexTmpl.Execute(w, nil); err != nil {
		panic(err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	writeTmpl(w, "500", err)
	w.WriteHeader(500)
}

func RenderTemplate(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeTmpl(w, name, nil)
	}
}

func GetPostRouter(get, post http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			get(w, r)
		} else if r.Method == "POST" {
			post(w, r)
		}
	}
}

func StartServer() {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		panic("failed to create test user")
	}
	store := &Store{db}
	users := &Users{}
	registrations := &Registrations{}
	sessions := &Sessions{}
	emails := NewEmailPrinter("http://localhost:9999")
	walls := &Walls{}
	collabConfig := CollabConfig{Debug: false, Workers: 10, Queue: 20, IOTimeout: time.Second * 5}

	LoadTestUser(store, users)

	postRegistration := PostRegistration(store, registrations, users, emails)
	loginHandler := PostLogin(store, users, registrations, sessions, emails)
	postForgotPassword := PostForgotPassword(store, users)
	postVerifyEmail := PostVerifyEmail(store, registrations, users)
	getRegistration := RenderTemplate("join")

	wrapper := noCacheMiddleware

	http.HandleFunc("/", wrapper(RenderTemplate("index")))
	http.HandleFunc("/sandbox", wrapper(RenderTemplate("sandbox")))
	http.HandleFunc("/terms", wrapper(RenderTemplate("terms")))

	http.HandleFunc("/join", wrapper(GetPostRouter(getRegistration, postRegistration)))
	http.HandleFunc("/login", wrapper(GetPostRouter(RenderTemplate("login"), loginHandler)))
	http.HandleFunc("/logout", wrapper(RenderTemplate("logout")))
	http.HandleFunc("/verify-email", wrapper(GetPostRouter(RenderTemplate("verify-email"), postVerifyEmail)))
	http.HandleFunc("/forgot-password", wrapper(GetPostRouter(RenderTemplate("forgot-password"), postForgotPassword)))
	http.HandleFunc("/reset-password", wrapper(GetPostRouter(RenderTemplate("reset-password"), postRegistration)))

	http.HandleFunc("/wall", wrapper(WallCollab(collabConfig, store, walls)))
	http.Handle("/static/", wrapper(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP))
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func LoadTestUser(db *Store, users CreateUser) {
	id := xid.New()
	email := Email("alice@example.com")
	password, _ := Password("abc").HashPassword()
	u := User{id, email, password, "alice"}
	err := db.Update(func(txn *Txn) error {
		return users.CreateUser(txn, &u)
	})
	if err != nil {
		panic("failed to create test user")
	}
}
