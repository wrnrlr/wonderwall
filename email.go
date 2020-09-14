package wonderwall

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"regexp"

	"github.com/markbates/pkger"
)

var (
	emailErr    = errors.New("invalid email")
	passwordErr = errors.New("invalid password")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Email string

func (e Email) valid() bool { return len(e) > 3 && len(e) < 255 && emailRegex.MatchString(string(e)) }

func (e Email) String() string { return string(e) }

type SendEmail interface {
	SendEmail(name string, email Email, i interface{}) error
}

type SendMailFunc func(string, smtp.Auth, string, []string, []byte) error

type EmailConfig struct {
	SenderAddr string
	Username   string
	Password   string
	ServerHost string
	ServerPort string
	Origin     string
}

type EmailClient struct {
	Conf      EmailConfig
	SendMail  SendMailFunc
	Templates map[string]template.Template
}

func NewEmailClient(from, addr, username, password, origin string) SendEmail {
	ts := loadTemplates()
	port := "587"
	conf := EmailConfig{from, username, password, addr, port, origin}
	return &EmailClient{conf, smtp.SendMail, ts}
}

func (s EmailClient) SendEmail(name string, to Email, i interface{}) error {
	t := s.Templates[name]
	name = name + ".eml"
	buf := new(bytes.Buffer)
	err := t.Execute(buf, i)
	if err != nil {
		log.Println("Error executing email template")
		return err
	}
	body := buf.Bytes()
	addr := s.Conf.ServerHost + ":" + s.Conf.ServerPort
	auth := smtp.PlainAuth("", s.Conf.Username, s.Conf.Password, s.Conf.ServerHost)
	return s.SendMail(addr, auth, s.Conf.SenderAddr, []string{string(to)}, body)
}

func loadTemplates() map[string]template.Template {
	d := pkger.Dir("/template/email")
	ts := map[string]template.Template{}
	var b []byte
	names := []string{"activate", "login", "reset", "welcome"}
	for _, name := range names {
		f, err := d.Open(name + ".eml")
		if err != nil {
			log.Fatalf("can't read %v", err)
		}
		b, err = ioutil.ReadAll(f)
		if err != nil {
			log.Fatalf("can't read %v", err)
		}
		ts[name] = *template.Must(template.New(name).Parse(string(b)))
	}
	return ts
}

func NewEmailPrinter(origin string) SendEmail {
	ts := loadTemplates()
	conf := EmailConfig{Origin: origin}
	return &EmailClient{SendMail: printSend, Conf: conf, Templates: ts}
}

func printSend(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	log.Printf(emailFmt, from, to, string(msg))
	return nil
}

const emailFmt = `
From: %s
To: %v
%s`
