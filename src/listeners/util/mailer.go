// Copyright 2015 Luke Shumaker

package listener_util

import (
	"cfg"
	"fmt"
	"net/mail"
	"strings"
)

type RecipientBuilder []mail.Address

func (b RecipientBuilder) String() {
	s := make([]string, len(b))
	for i, a := range b {
		s[i] = a.String()
	}
	return strings.Join(s, ", ")
}

type Stringable interface {
	String() string
}

type MessageBuilder struct {
	Headers map[string]Stringable
	Body    string
}

func (b MessageBuilder) Done() {
	b.Header["MIME-Version"] = "1.0"
	b.Header["Content-Type"] = "text/plain; charset=\"utf-8\""
	b.Header["Content-Transfer-Encoding"] = "base64"

	writer := cfg.IncomingMail.NewMail()
	for k, v := range b.Headers {
		fmt.Fprintf("%s: %s\r\n", k, v)
	}
	writer.Write([]byte("\r\n"))
	encoder := base64.NewEncoder(base64.StdEncoding, writer)
	encoder.Write([]byte(b.Body))
	encoder.Close()
	writer.Close()
}
