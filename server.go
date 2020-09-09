package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var (
	emailErr    = errors.New("invalid email")
	passwordErr = errors.New("invalid password")
)

type Email string

func (e Email) valid() bool { return len(e) > 3 && len(e) < 255 && emailRegex.MatchString(string(e)) }

type Password string

func (p Password) valid() bool { return len(p) > 8 && len(p) < 255 }
func (p Password) HashPassword() (PasswordHash, error) {
	return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
}

type EmailForm struct{ Email Email }

func (f EmailForm) validate() error {
	if !f.Email.valid() {
		return emailErr
	}
	return nil
}

type PasswordForm struct{ Password Password }

func (f PasswordForm) validate() error {
	if !f.Password.valid() {
		return passwordErr
	}
	return nil
}

type PasswordHash []byte

func (h PasswordHash) ComparePassword(p Password) error {
	return bcrypt.CompareHashAndPassword(h, []byte(p))
}

type Registration struct {
	ID           Token
	Email        Email
	PasswordHash PasswordHash
	CreatedAt    time.Time
	VerifiedAt   *time.Time
}

func (r *Registration) Key() Key {
	if r == nil {
		return Key("registration:")
	} else {
		return append([]byte("registration:"), r.ID...)
	}
}
func (r *Registration) Serialize() ([]byte, error) { return serialize(r) }
func (r *Registration) Deserialize(b []byte) error { return deserialize(b, r) }

type CreateRegistration interface {
	CreateRegistration(*Txn, Email, Password) (*Registration, error)
}
type FindRegistrationByID interface {
	FindRegistrationByID(*Txn, Token) (*Registration, error)
}
type FindRegistrationByEmail interface {
	FindRegistrationByEmail(*Txn, Email) (*Registration, error)
}
type RegistrationService interface {
	CreateUser
	FindRegistrationByID
	FindRegistrationByEmail
}
type Registrations struct{ DB *Store }

func (s Registrations) CreateRegistration(txn *Txn, email Email, password Password) (*Registration, error) {
	passwordHash, err := password.HashPassword()
	if err != nil {
		return nil, err
	}
	id, err := GenerateToken(32)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	r := &Registration{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: now}
	err = s.DB.Set(txn, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s Registrations) FindRegistrationByID(*Txn, Token) (*Registration, error) {
	return nil, nil
}

func (s Registrations) FindRegistrationByEmail(*Txn, Email) (*Registration, error) {
	return nil, nil
}

type User struct {
	ID           string
	Email        Email
	PasswordHash PasswordHash
}

func (u *User) Key() Key {
	if u == nil {
		return Key("user:")
	} else {
		return Key("user:" + u.ID)
	}
}

type Session struct {
	ID        Token
	UserID    string
	CreatedAt time.Time
}

type CreateSession interface {
	CreateSession(*Txn, string) (*Session, error)
}
type Sessions struct{}

func (s Sessions) CreateSession(*Txn, string) (*Session, error) {
	return nil, nil
}

type Wall struct {
	ID      string
	Content []interface{}
}

func (w *Wall) Key() Key {
	if w == nil {
		return Key("wall:")
	} else {
		return Key("wall:" + w.ID)
	}
}

type CreateWall interface {
	CreateWall(*Txn, *User) (*Wall, error)
}
type FindWallById interface {
	FindWallById(*Txn, string) (*Wall, error)
}
type DeleteWall interface {
	DeleteWall(*Txn, *Wall) error
}

type FindUserByEmail interface {
	FindUserByEmail(*Txn, Email) (*User, error)
}
type FindUserById interface {
	FindUserById(*Txn, string) (*User, error)
}
type CreateUser interface{ CreateUser(*Txn, *User) error }
type UpdateUser interface{ UpdateUser(*Txn, *User) error }
type DeleteUser interface{ DeleteUser(*Txn, *User) error }
type UserService interface {
	CreateUser
	UpdateUser
	DeleteUser
	FindUserById
	FindUserByEmail
}
type Users struct{ users []*Users }

func (s Users) CreateUser(*Txn, *User) error               { return nil }
func (s Users) UpdateUser(*Txn, *User) error               { return nil }
func (s Users) DeleteUser(*Txn, *User) error               { return nil }
func (s Users) FindUserById(*Txn, string) (*User, error)   { return nil, nil }
func (s Users) FindUserByEmail(*Txn, Email) (*User, error) { return nil, nil }

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

func writeError(w http.ResponseWriter, err error) { writeTmpl(w, "500", err); w.WriteHeader(500) }

func RenderTemplate(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeTmpl(w, name, nil)
	}
}

type Token []byte

func GenerateToken(n int) (Token, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Emails struct{}

func (s Emails) SendEmail(name string, i interface{}) error { return nil }

var duplicateEmailErr = errors.New("duplicate email")

type RegistrationForm struct {
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

func (f RegistrationForm) validate() error {
	if !f.Email.valid() {
		return emailErr
	} else if len(f.Password) < 8 {
		return passwordErr
	}
	return nil
}

func PostRegistration(db *Store, registrations CreateRegistration, users FindUserByEmail, emails *Emails) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			f   RegistrationForm
			u   *User
			reg *Registration
			err error
		)
		if !ContentType(applicationJson, w, r) {
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := f.validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.Update(func(txn *Txn) error {
			u, err = users.FindUserByEmail(txn, f.Email)
			if err != nil {
				return err
			} else if u != nil {
				return duplicateEmailErr
			}
			reg, err = registrations.CreateRegistration(txn, f.Email, f.Password)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		emails.SendEmail("verify-email", reg.ID)
	}
}

type LoginForm struct {
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

func (f LoginForm) validate() error {
	if !f.Email.valid() {
		return emailErr
	} else if len(f.Password) < 8 {
		return passwordErr
	}
	return nil
}

var emailNotVerifiedErr = errors.New("email not verified")
var emailNotRegistered = errors.New("email not registered")

func PostLogin(db *Store, users FindUserByEmail, registrations FindRegistrationByEmail, sessions CreateSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			f   LoginForm
			u   *User
			reg *Registration
			s   *Session
			err error
		)
		if !ContentType(applicationJson, w, r) {
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := f.validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.Update(func(txn *Txn) error {
			u, err = users.FindUserByEmail(txn, f.Email)
			if err != nil {
				return err
			}
			if u != nil {
				reg, err = registrations.FindRegistrationByEmail(txn, f.Email)
				if err != nil {
					return err
				} else if reg != nil {
					return emailNotVerifiedErr
				}
				return emailNotRegistered
			}
			s, err = sessions.CreateSession(txn, u.ID)
			if err != nil {
				return err
			}
			return nil
		})
		if err == emailNotVerifiedErr || err == emailNotRegistered {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO send new login email
		expiration := time.Now().Add(30 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session", Value: string(s.ID), Expires: expiration}
		http.SetCookie(w, &cookie)
	}
}

type ForgotPasswordForm struct {
	Email Email `json:"email`
}

func PostForgotPassword(db *Store, users FindUserByEmail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			f   ForgotPasswordForm
			u   *User
			err error
		)
		if !ContentType(applicationJson, w, r) {
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		db.Update(func(txn *Txn) error {
			u, err = users.FindUserByEmail(txn, f.Email)
			if err != nil {
				return nil
			}
			// TODO Insert reset token
			return nil
		})
		// TODO Send email
	}
}

type Verification struct{ Token Token }

func PostVerifyEmail(db *Store, registrations FindRegistrationByID, users CreateUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			f   Verification
			u   *User
			reg *Registration
			err error
		)
		if !ContentType(applicationJson, w, r) {
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.Update(func(txn *Txn) error {
			reg, err = registrations.FindRegistrationByID(txn, f.Token)
			u = &User{Email: reg.Email, PasswordHash: reg.PasswordHash}
			err = users.CreateUser(txn, u)
			return err
		})
		// TODO send welcome email
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

func main() {
	store := &Store{}
	users := &Users{}
	registrations := &Registrations{}
	sessions := &Sessions{}
	emails := &Emails{}
	postRegistration := PostRegistration(store, registrations, users, emails)
	postLogin := PostLogin(store, users, registrations, sessions)
	loginHandler := PostLogin(store, users, registrations, sessions)
	postForgotPassword := PostForgotPassword(store, users)
	postVerifyEmail := PostVerifyEmail(store, registrations, users)
	getRegistration := RenderTemplate("join")
	wrapper := noCacheMiddleware
	http.HandleFunc("/", wrapper(RenderTemplate("index")))
	http.HandleFunc("/sandbox", wrapper(RenderTemplate("sandbox")))
	http.HandleFunc("/join", wrapper(GetPostRouter(getRegistration, postRegistration)))
	http.HandleFunc("/terms", wrapper(RenderTemplate("terms")))
	http.HandleFunc("/login", wrapper(GetPostRouter(RenderTemplate("login"), postLogin)))
	http.HandleFunc("/logout", wrapper(RenderTemplate("logout")))
	http.HandleFunc("/verify-email", wrapper(GetPostRouter(RenderTemplate("verify-email"), postVerifyEmail)))
	http.HandleFunc("/forgot-password", wrapper(GetPostRouter(RenderTemplate("forgot-password"), postForgotPassword)))
	http.HandleFunc("/reset-password", wrapper(GetPostRouter(RenderTemplate("reset-password"), postRegistration)))
	http.HandleFunc("/walls", wrapper(loginHandler))
	http.HandleFunc("/wall", wrapper(loginHandler))
	http.Handle("/static/", wrapper(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP))
	log.Fatal(http.ListenAndServe(":9999", nil))
}
