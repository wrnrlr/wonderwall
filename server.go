package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var (
	emailErr    = errors.New("invalid email")
	passwordErr = errors.New("invalid password")
)

type Email string

func (e Email) valid() bool { return len(e) > 3 && len(e) < 255 && emailRegex.MatchString(string(e)) }

type Password string

func (e Password) valid() bool { return len(e) > 8 && len(e) < 255 }

type AuthForm struct {
	Email    Email
	Password Password
}

func (f AuthForm) validate() error {
	if !f.Email.valid() {
		return emailErr
	} else if len(f.Password) < 8 {
		return passwordErr
	}
	return nil
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
func main() {
	writeTmpl := func(w http.ResponseWriter, name string, i interface{}) {
		indexTmpl, err := template.ParseFiles(fmt.Sprintf("./template/%s.html", name))
		if err != nil {
			panic(err)
		}
		if err = indexTmpl.Execute(w, nil); err != nil {
			panic(err)
		}
	}
	writeError := func(w http.ResponseWriter, err error) { writeTmpl(w, "500", err); w.WriteHeader(500) }
	indexHandler := func(w http.ResponseWriter, r *http.Request) { writeTmpl(w, "index", nil) }
	sandboxHandler := func(w http.ResponseWriter, r *http.Request) { writeTmpl(w, "sandbox", nil) }
	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			writeTmpl(w, "login", nil)
			return
		}
		if r.Method == "POST" {
			var form AuthForm
			if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
				writeError(w, err)
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "session"})
		}
	}
	logoutHandler := func(w http.ResponseWriter, r *http.Request) {}
	forgotPasswordHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			writeTmpl(w, "forgot-password", nil)
			return
		}
		if r.Method == "POST" {
			var form EmailForm
			if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
				writeError(w, err)
				return
			}
			/*TODO*/
		}
	}
	registrationHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			writeTmpl(w, "join", nil)
			return
		}
		if r.Method == "POST" {
			var form AuthForm
			if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
				writeError(w, err)
				return
			}
		}
	}
	termsHandler := func(w http.ResponseWriter, r *http.Request) { writeTmpl(w, "terms", nil) }
	wrapper := noCacheMiddleware
	http.HandleFunc("/", wrapper(indexHandler))
	http.HandleFunc("/sandbox", wrapper(sandboxHandler))
	http.HandleFunc("/join", wrapper(registrationHandler))
	http.HandleFunc("/terms", wrapper(termsHandler))
	http.HandleFunc("/login", wrapper(loginHandler))
	http.HandleFunc("/logout", wrapper(logoutHandler))
	http.HandleFunc("/forgot-password", wrapper(forgotPasswordHandler))
	http.HandleFunc("/reset-password", wrapper(loginHandler))
	http.HandleFunc("/walls", wrapper(loginHandler))
	http.HandleFunc("/wall", wrapper(loginHandler))
	http.Handle("/static/", wrapper(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP))
	log.Fatal(http.ListenAndServe(":9999", nil))
}
