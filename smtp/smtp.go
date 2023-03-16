package ktn

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
	log "github.com/sirupsen/logrus"
)

var emailCounter int = 0

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
	log.Trace("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Trace("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		log.Error(err)
		return nil
	}

	emailCounter += 1

	addresses := ""

	alist, _ := env.AddressList("To")
	for _, addr := range alist {
		addresses += fmt.Sprintf("%s <%s>", addr.Name, addr.Address)
	}
	log.WithFields(log.Fields{
		"from":     env.GetHeader("From"),
		"subject":  env.GetHeader("Subject"),
		"to":       addresses,
		"len_txt":  strconv.Itoa(len(env.Text)),
		"len_html": strconv.Itoa(len(env.HTML)),
		"counter":  emailCounter,
	}).Info("Email received")

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func SmtpServer() *smtp.Server {
	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = "0.0.0.0:2525"
	s.Domain = "ktnrs.com"
	s.ReadTimeout = 30 * time.Second
	s.WriteTimeout = 30 * time.Second
	s.MaxMessageBytes = 2 * 1024 * 1024
	s.MaxRecipients = 1
	s.AuthDisabled = true

	return s

}
