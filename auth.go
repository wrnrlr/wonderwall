package wonderwall

import (
	"crypto/rand"
)

type EmailForm struct{ Email Email }

func (f EmailForm) validate() error {
	if !f.Email.valid() {
		return emailErr
	}
	return nil
}

type PasswordForm struct {
	Password Password
}

func (f PasswordForm) validate() error {
	if !f.Password.valid() {
		return passwordErr
	}
	return nil
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
