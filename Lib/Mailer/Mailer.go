package Mailer

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(recipient string, subject string, body string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("EMAIL_HOST")
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, 587, from, password)
	d.SSL = false

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
