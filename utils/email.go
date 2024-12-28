package utils

import (
	"crypto/tls"

	"github.com/go-mail/mail"
	"github.com/rs/zerolog/log"
)

type ConfigEmail struct {
	Username string
	Password string
	Driver   string
	Name     string
	To       []string
	Port     int
	IsTls    bool
	Sender   string
	Cc       *string
	Subject  string
	Body     string
}

func SendConfirmation(cfg ConfigEmail, attach *string) error {
	m := mail.NewMessage()
	m.SetHeader("From", cfg.Sender)
	m.SetHeader("To", cfg.To...)

	if cfg.Cc != nil {
		m.SetAddressHeader("Cc", *cfg.Cc, *cfg.Cc)
	}
	m.SetHeader("Subject", cfg.Subject)
	m.SetBody("text/html", cfg.Body)

	if attach != nil {
		m.Attach(*attach)
	}

	d := mail.NewDialer(cfg.Driver, cfg.Port, cfg.Username, cfg.Password)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: cfg.IsTls,
	}

	if err := d.DialAndSend(m); err != nil {
		log.Error().Err(err).Msg("[SendMail-1] Failed to Send Email")
		return err
	}

	return nil
}
