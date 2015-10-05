// Copyright 2015 Luke Shumaker

package httpentity

import (
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Headers http.Header
	Query   url.Values
	Entity  interface{}
}

type Response struct {
	status  int16
	Headers http.Header
	entity  NetEntity
}

type Encoder func(out io.Writer) error

type NetEntity interface {
	Encoders() map[string]Encoder
}

type Handler func(request Request) Response

type Entity interface {
	Methods() map[string]Handler
	Subentity(name string, request Request) Entity
}
