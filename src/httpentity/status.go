// Copyright 2015 Luke Shumaker

package httpentity

import (
	"httpentity/util"
	"net/http"
	"net/url"
)

// For when you're returning a document, with nothing special.
func StatusOK(entity NetEntity) Response {
	return Response{
		Status:  200,
		Headers: http.Header{},
		Entity:  entity,
	}
}

// For when you've created a document with a new URL.
func StatusCreated(parent Entity, child_name string, req Request) Response {
	if child_name == "" {
		panic(s("can't call StatusCreated with an empty child name"))
	}
	child := parent.Subentity(child_name, req)
	if child == nil {
		panic(s("called StatusCreated, but the subentity doesn't exist"))
	}
	handler, ok := child.Methods()["GET"]
	if !ok {
		panic(s("called StatusCreated, but can't GET the subentity"))
	}
	response := handler(req)
	response.Headers.Set("Location", url.QueryEscape(child_name))
	if response.Entity == nil {
		panic(s("called StatusCreated, but GET on subentity doesn't return an entity"))
	}
	mimetypes := encoders2mimetypes(response.Entity.Encoders())
	u, _ := url.Parse("") // create a blank dummy url.URL
	return Response{
		Status:  201,
		Headers: response.Headers,
		// XXX: .Entity gets modified by (*Router).route()
		// filled in the rest of the way by Route()
		Entity: mimetypes2net(u, mimetypes),
	}
}

// For when you've received a request, but haven't completed it yet
// (ex, it has been added to a queue).
func StatusAccepted(e NetEntity) Response {
	return Response{
		Status:  202,
		Headers: http.Header{},
		Entity:  e,
	}
}

// For when you've successfully done something, but have no body to
// return.
func StatusNoContent() Response {
	return Response{
		Status:  204,
		Headers: http.Header{},
		Entity:  nil,
	}
}

// The client should reset the form of whatever view it currently has.
func StatusResetContent() Response {
	return Response{
		Status:  205,
		Headers: http.Header{},
		Entity:  nil,
	}
}

// For when you have document in multiple formats, but you're not sure
// which the user wants.
func statusMultipleChoices(u *url.URL, mimetypes []string) Response {
	return Response{
		Status:  300,
		Headers: http.Header{},
		Entity:  mimetypes2net(u, mimetypes),
	}
}

// For when the document the user requested has permantly moved to a
// new address.
func StatusMovedPermanently(u *url.URL) Response {
	return Response{
		Status: 301,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: heutil.NetString(k("301 Moved")),
	}
}

// For when the document the user requested is currently found at
// another address, but that may not be the case in the future.
//
// The client may change a POST to a GET request when trying the new
// location.
func StatusFound(u *url.URL) Response {
	return Response{
		Status: 302,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: heutil.NetString(k("302 Found: ") + u.String()),
	}
}

func StatusSeeOther(u *url.URL) Response {
	return Response{
		Status: 303,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: heutil.NetString(k("303 See Other: ") + u.String()),
	}
}

// For when rhe document the user requested has temporarily moved.
//
// The client must repeate the request exactly the same, except for
// the URL.
func StatusTemporaryRedirect(u *url.URL) Response {
	return Response{
		Status: 307,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: heutil.NetString(k("307 Temporary Redirect: ") + u.String()),
	}
}

// For when the *user* has screwed up a request.
func statusBadRequest(e NetEntity) Response {
	if e == nil {
		e = heutil.NetString(k("400 Bad Request"))
	}
	return Response{
		Status:  400,
		Headers: http.Header{},
		Entity:  e,
	}
}

func StatusForbidden(e NetEntity) Response {
	if e == nil {
		e = heutil.NetString(k("403 Forbidden"))
	}
	return Response{
		Status:  403,
		Headers: http.Header{},
		Entity:  e,
	}
}

func statusNotFound() Response {
	return Response{
		Status:  404,
		Headers: http.Header{},
		Entity:  heutil.NetString(k("404 Not Found")),
	}
}

func statusMethodNotAllowed(methods string) Response {
	return Response{
		Status: 405,
		Headers: http.Header{
			"Allow": {methods},
		},
		Entity: heutil.NetString(k("405 Method Not Allowed")),
	}
}

func statusNotAcceptable(u *url.URL, mimetypes []string) Response {
	return Response{
		Status:  406,
		Headers: http.Header{},
		Entity:  mimetypes2net(u, mimetypes),
	}
}

// For when the user asked us to make a change conflicting with the
// current state of things.
func StatusConflict(entity NetEntity) Response {
	return Response{
		Status:  409,
		Headers: http.Header{},
		Entity:  entity,
	}
}

// For the resource has been deleted, and will never ever return.
func StatusGone(entity NetEntity) Response {
	return Response{
		Status:  410,
		Headers: http.Header{},
		Entity:  entity,
	}
}

func StatusUnsupportedMediaType(e NetEntity) Response {
	if e == nil {
		e = heutil.NetString(k("415 Unsupported Media Type"))
	}
	return Response{
		Status:  415,
		Headers: http.Header{},
		Entity:  e,
	}
}

// TODO: StatusExpectationFailed (417)
// TODO: StatusUpgradeRequired (426)

func statusInternalServerError(err interface{}) Response {
	return Response{
		Status: 500,
		Headers: http.Header{
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		Entity: heutil.NetPrintf(k("500 Internal Server Error: %v"), err),
	}
}

func StatusNotImplemented(e NetEntity) Response {
	return Response{
		Status:  501,
		Headers: http.Header{},
		Entity:  e,
	}
}
