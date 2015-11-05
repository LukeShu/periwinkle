// Copyright 2015 Davis Webb

package handlers

import (
	"io"
	"io/ioutil"
	"net/mail"
	"net/smtp"
)

func HandleEmail(r io.Reader, name string) int {
	m, err := mail.ReadMessage(r)
	if err != nil {
		panic(err)
	}

	header := m.Header
	auth := smtp.PlainAuth("", "TODO@TODO.com", "password", "mail.example.com")
	from := header.Get("From")
	to := []string{header.Get("To")}
	// sbj  := header.Get("Subject")
	body, _ := ioutil.ReadAll(m.Body)
	err = smtp.SendMail("mail.example.com:25", auth, from, []string(to), body)
	if err != nil {
		panic(err)
	}
	panic("TODO")
}
