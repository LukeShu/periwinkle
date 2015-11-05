// Copyright 2015 Luke Shumaker

package putil

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"periwinkle/cfg"
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
	Headers map[string]string
	Body    string
}

func (b MessageBuilder) Done() {
	b.Headers["MIME-Version"] = "1.0"
	b.Headers["Content-Type"] = "text/plain; charset=\"utf-8\""
	b.Headers["Content-Transfer-Encoding"] = "base64"

	writer := cfg.IncomingMail.NewMail()
	for k, v := range b.Headers {
		fmt.Fprintf(writer, "%s: %s\r\n", k, v)
	}
	writer.Write([]byte("\r\n"))
	encoder := base64.NewEncoder(base64.StdEncoding, writer)
	encoder.Write([]byte(b.Body))
	encoder.Close()
	writer.Close()
}
