package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"net"
	"strconv"
	"strings"

	gomail "gopkg.in/mail.v2"
)

func sendMail() {
	m := gomail.NewMessage()
	m.SetHeader("From", "foobar@apple.com")
	m.SetHeader("To", "bbla")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello")
	m.SetBody("text", "Hello")
	d := gomail.Dialer{Host: "127.0.0.1", Port: 587}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
