package main

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var (
	emailErr    = errors.New("invalid email")
	passwordErr = errors.New("invalid password")
)

type Email string

func (e Email) valid() bool { return len(e) > 3 && len(e) < 255 && emailRegex.MatchString(string(e)) }

type Emails struct{}

func (s Emails) SendEmail(name string, i interface{}) error { return nil }
