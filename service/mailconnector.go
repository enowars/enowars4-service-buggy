package main

import (
	"github.com/enowars/enowars4-service-buggy/Models"
	gomail "gopkg.in/mail.v2"
	"log"
	"strings"
)

//ToDo: config file for mail sending
func sendMail(msg Models.Msg) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To)
	if msg.Cc != "" {
		//ToDo ensure @ is in string
		m.SetAddressHeader("Cc", msg.Cc, msg.Cc[:strings.IndexByte(msg.Cc, '@')])
	}
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text", msg.Body)
	d := gomail.Dialer{Host: "127.0.0.1", Port: 587}

	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
	}
	return
}
