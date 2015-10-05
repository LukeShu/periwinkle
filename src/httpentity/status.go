// Copyright 2015 Luke Shumaker

package httpentity

import (
	"io"
	"net/http"
	"net/url"
)

type netString string

func (s netString) Encoders() map[string]Encoder {
	return map[string]Encoder{"text/plain": s.write}
}

func (s netString) write(w io.Writer) error {
	_, err := w.Write([]byte(s))
	return err
}

func (req Request) StatusOK(entity NetEntity) Response {
	return Response{
		Status:  200,
		Headers: http.Header{},
		Entity:  entity,
	}
}

func (req Request) StatusCreated(parent Entity, child_name string) Response {
	child := parent.Subentity(child_name, req)
	if child == nil {
		panic("called StatusCreated, but the subentity doesn't exist")
	}
	handler, ok := child.Methods()["GET"]
	if !ok {
		panic("called StatusCreated, but can't GET the subentity")
	}
	response := handler(req)
	response.Headers.Set("Location", url.QueryEscape(child_name))
	return Response{
		Status:  201,
		Headers: response.Headers,
		Entity:  response.Entity,
	}
}

func (req Request) statusMultipleChoices(u *url.URL, mimetypes []string) Response {
	panic("not implemented")
}

func (req Request) StatusMoved(url *url.URL) Response {
	return Response{
		Status: 301,
		Headers: http.Header{
			"Location":     {url.String()},
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		Entity: netString("301: Moved"),
	}
}

func (req Request) StatusFound(url *url.URL) Response {
	return Response{
		Status: 302,
		Headers: http.Header{
			"Location":     {url.String()},
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		Entity: netString("302: Found"),
	}
}

func (req Request) statusNotFound() Response {
	return Response{
		Status:  404,
		Headers: http.Header{"Content-Type": {"text/plain; charset=utf-8"}},
		Entity:  netString("404 Not Found"),
	}
}

func (req Request) statusMethodNotAllowed(methods string) Response {
	return Response{
		Status: 405,
		Headers: http.Header{
			"Allow":        {methods},
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		Entity: netString("405 Method Not Allowed"),
	}
}

func (req Request) statusNotAcceptable(url *url.URL, mimetypes []string) Response {
	panic("not implemented")
}

func (req Request) statusInternalServerError() Response {
	return Response{
		Status:  500,
		Headers: http.Header{"Content-Type": {"text/plain; charset=utf-8"}},
		Entity:  netString("500 Internal Server Error"),
	}
}
