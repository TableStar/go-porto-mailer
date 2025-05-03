package emailer

import (
	"fmt"
	"net/smtp"
)

type EmailSender interface {
	Send(to, subject, body string) error
}

type smtpSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

//func for new instances of smtp email sender

func NewSmtpSender(host, port, username, password, from string) EmailSender {
	return &smtpSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *smtpSender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := fmt.Sprintf("From %s\nTo: %s\nSubject: %s\n\n%s", s.from, to, subject, body)

	msgBytes := []byte(msg)

	err := smtp.SendMail(addr, auth, s.from, []string{to}, msgBytes)

	if err != nil {
		return fmt.Errorf("smtp.SendMail is failed even if it should not be error but got error: %w", err)

	}
	return nil
}
