package MailSender

import (
	"crypto/tls"

	mail "gopkg.in/gomail.v2"
)

func mailSend(from AuthData, To string, theme string, message string) error {
	m := mail.NewMessage()
	m.SetHeader("From", from.user+"@mail.distet.tech")
	m.SetHeader("To", To)
	m.SetHeader("Subject", theme)
	m.SetBody("text/html", message)
	d := mail.NewDialer("localhost", 25, "distet", "258013006402582")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func MailAuthCode(To string, nickname string, Code string, aboutuser aboutuser) error {
	return mailSend(AuthData{}.NoReply(), To, "Код подтверждения", authHtmlWithData(acceptemail{Name: nickname, Code: Code, Aboutuser: aboutuser}))
}
