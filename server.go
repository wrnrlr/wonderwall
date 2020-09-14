package wonderwall

import (
	"fmt"
	"github.com/julienschmidt/httprouter"

	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger/v2"
	"github.com/rs/xid"
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
	sessionManager := scs.New()
	sessionManager.Lifetime = 15 * 24 * time.Hour

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		panic("failed to create test user")
	}
	store := &Store{db}
	users := &Users{}
	registrations := &Registrations{}
	emails := NewEmailPrinter("http://localhost:9999")
	walls := &Walls{}
	security := Auth{}
	collabConfig := CollabConfig{Debug: false, Workers: 10, Queue: 20, IOTimeout: time.Second * 5}

	LoadTestUser(store, users)

	postRegistration := PostRegistration(store, registrations, users, emails)
	loginHandler := PostLogin(store, users, registrations, sessionManager, emails)
	logoutHandler := Logout(sessionManager)
	postForgotPassword := PostForgotPassword(store, users)
	postVerifyEmail := PostVerifyEmail(store, registrations, users)

	wrapper := noCacheMiddleware

	router := httprouter.New()

	router.HandlerFunc("GET", "/", wrapper(RenderTemplate("index")))
	router.HandlerFunc("GET", "/sandbox", wrapper(RenderTemplate("sandbox")))
	router.HandlerFunc("GET", "/terms", wrapper(RenderTemplate("terms")))

	router.HandlerFunc("GET", "/join", wrapper(RenderTemplate("join")))
	router.HandlerFunc("POST", "/join", wrapper(postRegistration))

	router.HandlerFunc("GET", "/login", wrapper(RenderTemplate("login")))
	router.HandlerFunc("POST", "/login", wrapper(loginHandler))

	router.HandlerFunc("GET", "/logout", wrapper(logoutHandler))

	router.HandlerFunc("GET", "/verify-email", wrapper(RenderTemplate("verify-email")))
	router.HandlerFunc("POST", "/verify-email", wrapper(postVerifyEmail))

	router.HandlerFunc("GET", "/forgot-password", wrapper(RenderTemplate("forgot-password")))
	router.HandlerFunc("POST", "/forgot-password", wrapper(postForgotPassword))

	router.HandlerFunc("GET", "/reset-password", wrapper(RenderTemplate("reset-password")))
	router.HandlerFunc("POST", "/reset-password", wrapper(postRegistration))

	router.HandlerFunc("GET", "/wall", wrapper(WallCollab(collabConfig, store, walls)))
	router.HandlerFunc("POST", "/wall", wrapper(PostWallHandler(security)))
	router.HandlerFunc("GET", "/wall/:id", wrapper(WallCollab(collabConfig, store, walls)))
	router.HandlerFunc("PATCH", "/wall/:id", wrapper(PatchWallHandler(security)))
	router.HandlerFunc("DELETE", "/wall/:id", wrapper(DeleteWallHandler(security)))

	router.HandlerFunc("GET", "/static/*filepath", wrapper(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP))

	log.Fatal(http.ListenAndServe(":9999", router))
}

func LoadTestUser(db *Store, users CreateUser) {
	id := xid.New()
	email := Email("alice@example.com")
	password, _ := Password("Abcd1234").HashPassword()
	u := User{id, email, password, "alice"}
	err := db.Update(func(txn *Txn) error {
		return users.CreateUser(txn, &u)
	})
	if err != nil {
		panic("failed to create test user")
	}
}
