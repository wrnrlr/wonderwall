package wonderwall

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
)

func StartServer() {
	sessionManager := scs.New()
	sessionManager.Lifetime = 15 * 24 * time.Hour

	store, err := NewStore(StoreConfig{"", true})
	if err != nil {
		panic("failed to create test user")
	}

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

	router.HandlerFunc("GET", "/wall", wrapper(GetWallHandler(security)))
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
	u := User{ID: id, Email: email, PasswordHash: password, Name: "alice"}
	err := db.Update(func(txn *Txn) error {
		return users.CreateUser(txn, &u)
	})
	if err != nil {
		panic("failed to create test user")
	}
}
