package wonderwall

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

type LoginForm struct {
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

func (f LoginForm) validate() error {
	if !f.Email.Valid() {
		return EmailErr
	} else if len(f.Password) < 8 {
		return PasswordErr
	}
	return nil
}

var emailNotVerifiedErr = errors.New("email not verified")
var emailNotRegistered = errors.New("email not registered")

func PostLogin(db *Store, users FindUserByEmail, registrations FindRegistrationByEmail, sessionManager *scs.SessionManager, emails SendEmail) http.HandlerFunc {
	origin := "http://localhost:9999"
	type LoginMsg struct {
		Name   string
		Origin string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			f   LoginForm
			u   *User
			reg *Registration
			//s   *Session
			err error
		)
		//if !ContentType(applicationJson, w, r) {
		//	return
		//}
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
			if u == nil {
				reg, err = registrations.FindRegistrationByEmail(txn, f.Email)
				if err != nil {
					return err
				} else if reg != nil {
					return emailNotVerifiedErr
				}
				return emailNotRegistered
			}
			// TODO store session
			//s, err = sessions.CreateSession(txn, u.ID.String())
			//if err != nil {
			//	return err
			//}
			return nil
		})
		if err == emailNotVerifiedErr || err == emailNotRegistered {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Renew token to prevent session fixation
		//err = sessionManager.RenewToken(r.Context())
		//if err != nil {
		//	http.Error(w, err.Error(), 500)
		//	return
		//}
		// Then make the privilege-level change.
		//sessionManager.Put(r.Context(), "userID", u.ID.String())
		//msg := &LoginMsg{Name: u.Name, Origin: origin}
		//emails.SendEmail("login", u.Email, msg)

		expiration := time.Now().Add(30 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session", Value: string("hello"), Expires: expiration}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, origin+"/walls", http.StatusFound)
	}
}

func Logout(sessionManager *scs.SessionManager) http.HandlerFunc {
	origin := "http://localhost:9999"
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager.Remove(r.Context(), "userID")
		http.Redirect(w, r, origin, http.StatusFound)
	}
}
