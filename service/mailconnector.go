package main

import (
	"github.com/enowars/enowars4-service-buggy/service/Models"
	gomail "gopkg.in/mail.v2"
	"log"
)

func sendMail(msg MSG) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "foobar@apple.com")
	m.SetHeader("To", "bbla")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello")
	m.SetBody("text", "Hello")
	d := gomail.Dialer{Host: "127.0.0.1", Port: 587}

	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
	}
	return
}
