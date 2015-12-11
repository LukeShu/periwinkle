// Copyright 2015 Luke Shumaker

package rfc7231

import (
	he "httpentity"
	"net/http"
	"net/url"
)

// For when you're returning a document, with nothing special.
func StatusOK(entity he.NetEntity) he.Response {
	return he.Response{
		Status:  200,
		Headers: http.Header{},
		Entity:  entity,
	}
}

// For when you've created a document with a new URL.
func StatusCreated(parent he.EntityGroup, childName string, req he.Request) he.Response {
	if childName == "" {
		panic("can't call StatusCreated with an empty child name")
	}
	// find the child
	child := parent.Subentity(childName, req)
	if child == nil {
		panic("called StatusCreated, but the subentity doesn't exist")
	}
	// prepare the response
	u, _ := req.URL.Parse(url.QueryEscape(childName))
	return he.Response{
		Status: 201,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity:                 he.NetPrintf("%s", u.String()),
		InhibitNotAcceptable:   true,
		InhibitMultipleChoices: true,
	}
}

// For when you've received a request, but haven't completed it yet
// (ex, it has been added to a queue).
func StatusAccepted(e he.NetEntity) he.Response {
	return he.Response{
		Status:  202,
		Headers: http.Header{},
		Entity:  e,
	}
}

// For when you've successfully done something, but have no body to
// return.
func StatusNoContent() he.Response {
	return he.Response{
		Status:  204,
		Headers: http.Header{},
		Entity:  nil,
	}
}

// The client should reset the form of whatever view it currently has.
func StatusResetContent() he.Response {
	return he.Response{
		Status:  205,
		Headers: http.Header{},
		Entity:  nil,
	}
}

// For when you have document in multiple formats, but you're not sure
// which the user wants.
func StatusMultipleChoices(u *url.URL, mimetypes []string) he.Response {
	return he.Response{
		Status:  300,
		Headers: http.Header{},
		Entity:  mimetypes2net(u, mimetypes),
	}
}

// For when the document the user requested has permantly moved to a
// new address.
func StatusMovedPermanently(u *url.URL) he.Response {
	return he.Response{
		Status: 301,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: he.NetPrintf("301 Moved"),
	}
}

// For when the document the user requested is currently found at
// another address, but that may not be the case in the future.
//
// The client may change a POST to a GET request when trying the new
// location.
func StatusFound(u *url.URL) he.Response {
	return he.Response{
		Status: 302,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: he.NetPrintf("302 Found: %v", u),
	}
}

func StatusSeeOther(u *url.URL) he.Response {
	return he.Response{
		Status: 303,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: he.NetPrintf("303 See Other: %v", u),
	}
}

// For when rhe document the user requested has temporarily moved.
//
// The client must repeate the request exactly the same, except for
// the URL.
func StatusTemporaryRedirect(u *url.URL) he.Response {
	return he.Response{
		Status: 307,
		Headers: http.Header{
			"Location": {u.String()},
		},
		Entity: he.NetPrintf("307 Temporary Redirect: %v", u),
	}
}

// For when the *user* has screwed up a request.
func StatusBadRequest(e he.NetEntity) he.Response {
	if e == nil {
		e = he.NetPrintf("400 Bad Request")
	}
	return he.Response{
		Status:  400,
		Headers: http.Header{},
		Entity:  e,
	}
}

func StatusForbidden(e he.NetEntity) he.Response {
	if e == nil {
		e = he.NetPrintf("403 Forbidden")
	}
	return he.Response{
		Status:  403,
		Headers: http.Header{},
		Entity:  e,
	}
}

func StatusNotFound(e he.NetEntity) he.Response {
	if e == nil {
		e = he.NetPrintf("404 Not Found")
	}
	return he.Response{
		Status:                 404,
		Headers:                http.Header{},
		Entity:                 e,
		InhibitNotAcceptable:   true,
		InhibitMultipleChoices: true,
	}
}

func StatusMethodNotAllowed(entity he.Entity, request he.Request) he.Response {
	return he.Response{
		Status: 405,
		Headers: http.Header{
			"Allow": {methods2string(entity.Methods())},
		},
		Entity:                 he.NetPrintf("405 Method Not Allowed"),
		InhibitNotAcceptable:   true,
		InhibitMultipleChoices: true,
	}
}

func StatusNotAcceptable(u *url.URL, mimetypes []string) he.Response {
	return he.Response{
		Status:               406,
		Headers:              http.Header{},
		Entity:               mimetypes2net(u, mimetypes),
		InhibitNotAcceptable: true,
	}
}

// For when the user asked us to make a change conflicting with the
// current state of things.
func StatusConflict(entity he.NetEntity) he.Response {
	return he.Response{
		Status:  409,
		Headers: http.Header{},
		Entity:  entity,
	}
}

// For the resource has been deleted, and will never ever return.
func StatusGone(entity he.NetEntity) he.Response {
	return he.Response{
		Status:  410,
		Headers: http.Header{},
		Entity:  entity,
	}
}

func StatusUnsupportedMediaType(e he.NetEntity) he.Response {
	if e == nil {
		e = he.NetPrintf("415 Unsupported Media Type")
	}
	return he.Response{
		Status:  415,
		Headers: http.Header{},
		Entity:  e,
	}
}

// TODO: StatusExpectationFailed (417)
// TODO: StatusUpgradeRequired (426)

func StatusInternalServerError(err interface{}) he.Response {
	return he.Response{
		Status: 500,
		Headers: http.Header{
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		Entity: he.NetPrintf("500 Internal Server Error: %v", err),
	}
}

func StatusNotImplemented(e he.NetEntity) he.Response {
	return he.Response{
		Status:  501,
		Headers: http.Header{},
		Entity:  e,
	}
}
