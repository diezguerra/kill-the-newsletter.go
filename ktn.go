package main

import (
	"errors"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

// The Backend implements SMTP server methods.
type Backend struct{}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after EHLO.
type Session struct{}

func (s *Session) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return errors.New("invalid username or password")
	}
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Subject: " + env.GetHeader("Subject"))
	log.Println("From: " + env.GetHeader("From"))
	alist, _ := env.AddressList("To")
	for _, addr := range alist {
		log.Println("To: ", addr.Name, addr.Address)
	}
	log.Println("Text: " + strconv.Itoa(len(env.Text)))
	log.Println("HTML: " + strconv.Itoa(len(env.HTML)))

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {
	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = ":2525"
	s.Domain = "ktnrs.com"
	s.ReadTimeout = 30 * time.Second
	s.WriteTimeout = 30 * time.Second
	s.MaxMessageBytes = 2 * 1024 * 1024
	s.MaxRecipients = 1
	s.AuthDisabled = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
