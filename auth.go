package main

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
