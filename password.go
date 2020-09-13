package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Password string

func (p Password) valid() bool { return len(p) > 8 && len(p) < 255 }

func (p Password) HashPassword() (PasswordHash, error) {
	return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
}

type PasswordHash []byte

func (h PasswordHash) ComparePassword(p Password) error {
	return bcrypt.CompareHashAndPassword(h, []byte(p))
}

type ForgotPasswordForm struct {
	Email Email `json:"email"`
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

func PostResetPassword(db *Store, users FindUserByEmail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
