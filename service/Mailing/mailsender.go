package Mailing

import (
	gomail "gopkg.in/mail.v2"
	"log"
	"strings"
)

func sendMail(msg Msg, conf mailConf) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To)
	if msg.Cc != "" {
		//ToDo ensure @ is in string, maybe inject vuln here?
		m.SetAddressHeader("Cc", msg.Cc, msg.Cc[:strings.IndexByte(msg.Cc, '@')])
	}
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text", msg.Body)
	d := gomail.Dialer{Host: conf.Host, Port: conf.Port}

	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
	}
	return
}
