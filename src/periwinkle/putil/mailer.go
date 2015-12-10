// Copyright 2015 Luke Shumaker

package putil

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"maildir"
	"net/mail"
	"net/smtp"
	"strings"
)

type RecipientBuilder []mail.Address

func (b RecipientBuilder) String() string {
	s := make([]string, len(b))
	for i, a := range b {
		s[i] = a.String()
	}
	return strings.Join(s, ", ")
}

type MessageBuilder struct {
	Maildir maildir.Maildir
	Headers map[string]string
	Body    string
}

func (b MessageBuilder) Done() {
	b.Headers["MIME-Version"] = "1.0"
	b.Headers["Content-Type"] = "text/plain; charset=\"utf-8\""
	b.Headers["Content-Transfer-Encoding"] = "base64"

	msg822 := []byte{}
	for k := range b.Headers {
		msg822 = append(msg822, []byte(fmt.Sprintf("%s: %s\r\n", k, b.Headers[k]))...)
	}
	msg822 = append(msg822, []byte("\r\n")...)

	var body bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &body)
	encoder.Write([]byte(b.Body))
	encoder.Close()
	msg822 = append(msg822, body.Bytes()...)

	to_addrs, err := mail.ParseAddressList(b.Headers["To"])
	if err != nil {
		panic(err) // FIXME
	}
	to_strs := make([]string, len(to_addrs))
	for i, addr := range to_addrs {
		to_strs[i] = addr.Address
	}

	if len(to_strs) > 0 {
		// send the message out
		err = smtp.SendMail("localhost:25",
			smtp.PlainAuth("", "", "", ""),
			b.Headers["From"],
			to_strs,
			msg822)
		if err != nil {
			panic(err) // FIXME
		}
	}
}
