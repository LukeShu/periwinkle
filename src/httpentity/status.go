// Copyright 2015 Luke Shumaker

package httpentity

import (
	"net/http"
	"net/url"
	"fmt"
)

func (req Request) StatusOK(entity NetEntity) Response {
	return Response{
		status:  200,
		Headers: http.Header{},
		entity:  entity,
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
	if response.entity == nil {
		panic("called StatusCreated, but GET on subentity doesn't return an entity")
	}
	mimetypes := encoders2mimetypes(response.entity.Encoders())
	u, _ := url.Parse("")
	return Response{
		status:  201,
		Headers: response.Headers,
		entity:  mimetypes2net(u, mimetypes), // gets modified by Route() filled in the rest of the way by Route()
	}
}

func (req Request) StatusNoContent() Response {
	return Response{
		status:  204,
		Headers: http.Header{},
		entity:  nil,
	}
}

func (req Request) statusMultipleChoices(u *url.URL, mimetypes []string) Response {
	return Response{
		status: 300,
		Headers: http.Header{},
		entity: mimetypes2net(u, mimetypes),
	}
}

func (req Request) StatusMoved(u *url.URL) Response {
	return Response{
		status: 301,
		Headers: http.Header{
			"Location":     {u.String()},
		},
		entity: NetString("301: Moved"),
	}
}

func (req Request) StatusFound(u *url.URL) Response {
	return Response{
		status: 302,
		Headers: http.Header{
			"Location":     {u.String()},
		},
		entity: NetString("302: Found"),
	}
}

func (req Request) StatusUnauthorized(e NetEntity) Response {
	return Response{
		status:  401,
		Headers: http.Header{
			"WWW-Authenticate": {"TODO"},
		},
		entity:  e,
	}
}

func (req Request) statusNotFound() Response {
	return Response{
		status:  404,
		Headers: http.Header{},
		entity:  NetString("404 Not Found"),
	}
}

func (req Request) statusMethodNotAllowed(methods string) Response {
	return Response{
		status: 405,
		Headers: http.Header{},
		entity: NetString("405 Method Not Allowed"),
	}
}

func (req Request) statusNotAcceptable(u *url.URL, mimetypes []string) Response {
	return Response{
		status: 406,
		Headers: http.Header{},
		entity: mimetypes2net(u, mimetypes),
	}
}

func (req Request) statusInternalServerError(err interface{}) Response {
	return Response{
		status:  500,
		Headers: http.Header{
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		entity:  NetString(fmt.Sprintf("500 Internal Server Error: %v", err)),
	}
}
