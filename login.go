package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

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
