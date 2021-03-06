package wonderwall

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

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
func (r *Registration) Serialize() ([]byte, error) { return Serialize(r) }
func (r *Registration) Deserialize(b []byte) error { return Deserialize(b, r) }

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

var DuplicateEmailErr = errors.New("duplicate email")

type RegistrationForm struct {
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

func (f RegistrationForm) Validate() error {
	if !f.Email.Valid() {
		return EmailErr
	} else if len(f.Password) < 8 {
		return PasswordErr
	}
	return nil
}

func PostRegistration(db *Store, registrations CreateRegistration, users FindUserByEmail, emails SendEmail) http.HandlerFunc {

	type ActivateMsg struct {
		Name   string
		Origin string
		Token  string
	}

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
		if err := f.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.Update(func(txn *Txn) error {
			u, err = users.FindUserByEmail(txn, f.Email)
			if err != nil {
				return err
			} else if u != nil {
				return DuplicateEmailErr
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
		emails.SendEmail("verify-email", u.Email, reg.ID)
	}
}

type Verification struct{ Token Token }

func PostVerifyEmail(db *Store, registrations FindRegistrationByID, users CreateUser) http.HandlerFunc {
	type WelcomeMsg struct {
		Name   string
		Origin string
	}
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
