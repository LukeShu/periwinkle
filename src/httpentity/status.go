// Copyright 2015 Luke Shumaker

package httpentity

import (
	"fmt"
	"httpentity/util"
	"net/http"
	"net/url"
)

// For when you're returning a document, with nothing special.
func (req Request) StatusOK(entity NetEntity) Response {
	return Response{
		Status:  200,
		Headers: http.Header{},
		Entity:  entity,
	}
}

// For when you've created a document with a new URL.
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
	if response.Entity == nil {
		panic("called StatusCreated, but GET on subentity doesn't return an entity")
	}
	mimetypes := encoders2mimetypes(response.Entity.Encoders())
	u, _ := url.Parse("")
	return Response{
		Status:  201,
		Headers: response.Headers,
		// XXX: .entity gets modified by (*Router).route()
		// filled in the rest of the way by Route()
		Entity: mimetypes2net(u, mimetypes),
	}
}

// For when you've received a request, but haven't completed it yet
// (ex, it has been added to a queue).
func (req Request) StatusAccepted(e NetEntity) Response {
	return Response{
		Status:  202,
		Headers: http.Header{},
		Entity:  e,
	}
}

// For when you've successfully done something, but have no body to
// return.
func (req Request) StatusNoContent() Response {
	return Response{
		Status:  204,
		Headers: http.Header{},
		Entity:  nil,
	}
}

// The client should reset the form of whatever view it currently has.
func (req Request) StatusResetContent() Response {
	return Response{
		Status:  205,
		Headers: http.Header{},
		Entity:  nil,
	}
}

// For when you have document in multiple formats, but you're not sure
// which the user wants.
func (req Request) statusMultipleChoices(u *url.URL, mimetypes []string) Response {
	return Response{
		Status:  300,
		Headers: http.Header{},
		Entity:  mimetypes2net(u, mimetypes),
	}
}

// For when the document the user requested has permantly moved to a
// new address.
func (req Request) StatusMovedPermanently(u *url.URL) Response {
	return Response{
		Status: 301,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: heutil.NetString("301 Moved"),
	}
}

// For when the document the user requested is currently found at
// another address, but that may not be the case in the future.
//
// The client may change a POST to a GET request when trying the new
// location.
func (req Request) StatusFound(u *url.URL) Response {
	return Response{
		Status: 302,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: heutil.NetString("302 Found: " + u.String()),
	}
}

func (req Request) StatusSeeOther(u *url.URL) Response {
	return Response{
		Status: 303,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: heutil.NetString("303 See Other: " + u.String()),
	}
}

// For when rhe document the user requested has temporarily moved.
//
// The client must repeate the request exactly the same, except for
// the URL.
func (req Request) StatusTemporaryRedirect(u *url.URL) Response {
	return Response{
		Status: 307,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: heutil.NetString("307 Temporary Redirect: " + u.String()),
	}
}

// For when the *user* has screwed up a request.
func (req Request) statusBadRequest(e NetEntity) Response {
	if e == nil {
		e = heutil.NetString("400 Bad Request")
	}
	return Response{
		Status:  400,
		Headers: http.Header{},
		Entity:  e,
	}
}

func (req Request) StatusForbidden(e NetEntity) Response {
	if e == nil {
		e = heutil.NetString("403 Forbidden")
	}
	return Response{
		Status:  403,
		Headers: http.Header{},
		Entity:  e,
	}
}

func (req Request) statusNotFound() Response {
	return Response{
		Status:  404,
		Headers: http.Header{},
		Entity:  heutil.NetString("404 Not Found"),
	}
}

func (req Request) statusMethodNotAllowed(methods string) Response {
	return Response{
		Status: 405,
		Headers: http.Header{
			"Allow": {methods},
		},
		Entity: heutil.NetString("405 Method Not Allowed"),
	}
}

func (req Request) statusNotAcceptable(u *url.URL, mimetypes []string) Response {
	return Response{
		Status:  406,
		Headers: http.Header{},
		Entity:  mimetypes2net(u, mimetypes),
	}
}

// For when the user asked us to make a change conflicting with the
// current state of things.
func (req Request) StatusConflict(entity NetEntity) Response {
	return Response{
		Status:  409,
		Headers: http.Header{},
		Entity:  entity,
	}
}

// For the resource has been deleted, and will never ever return.
func (req Request) StatusGone(entity NetEntity) Response {
	return Response{
		Status:  410,
		Headers: http.Header{},
		Entity:  entity,
	}
}

func (req Request) StatusUnsupportedMediaType(e NetEntity) Response {
	if e == nil {
		e = heutil.NetString("415 Unsupported Media Type")
	}
	return Response{
		Status:  415,
		Headers: http.Header{},
		Entity:  e,
	}
}

// TODO: StatusExpectationFailed (417)
// TODO: StatusUpgradeRequired (426)

func (req Request) statusInternalServerError(err interface{}) Response {
	return Response{
		Status: 500,
		Headers: http.Header{
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		Entity: heutil.NetString(fmt.Sprintf("500 Internal Server Error: %v", err)),
	}
}

func (req Request) StatusNotImplemented(e NetEntity) Response {
	return Response{
		Status:  501,
		Headers: http.Header{},
		Entity:  e,
	}
}
