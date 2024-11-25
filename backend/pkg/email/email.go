package email

import (
	"crypto/tls"

	gomail "gopkg.in/mail.v2"

	"github.com/Cheasezz/anSpace/backend/config"
)

type Sender interface {
	Send(to, message string) error
}

type EmailSender struct {
	dialer *gomail.Dialer
}

func NewSender(cfg config.EmailSender) (*EmailSender, error) {
	d := gomail.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.From, cfg.Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &EmailSender{
		dialer: d,
	}, nil
}

func (s *EmailSender) Send(to, message string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.dialer.Username)

	m.SetHeader("To", to)

	m.SetHeader("Subject", "Gomail test subject")

	m.SetBody("text/plain", message)

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}