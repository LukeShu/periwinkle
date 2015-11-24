// Copyright 2015 Luke Shumaker

// Package httpentity provides a "framework" for exposing resources
// over HTTP.
//
// Within the framework, an "Entity" is simply something that can be
// accessed over HTTP.  A "NetEntity" is something that is capable of
// being transmitted as an HTTP body--most Entities will also be
// NetEntities.  But, things like error messages are NetEntities, but
// aren't Entities, as it's not a normal thing that you can request.
//
// An Entity is accessed using the method allowed in
// (*Entity).Methods().  An entity may also have children--these are
// accessed with (*Entity).Subentities(childName, request).
//
// A Router is the entire Entity tree; it mostly just takes the root
// Entity of the tree.  It handles dispatching Requests to the correct
// Entity, then formatting the Responses to the output stream.
package httpentity

import (
	"io"
	"net/http"
	"net/url"
)

////////////////////////////////////////////////////////////////////////////////

// A Router represents the root of an Entity tree, and handles reading
// and writing messages to the network socket.
type Router struct {
	Prefix      string
	Root        RootEntity
	Decoders    map[string]func(io.Reader, map[string]string) (interface{}, error)
	Middlewares []Middleware

	// Whether to include stacktraces in HTTP 500 responses
	Stacktrace bool

	// Whether to log requests
	LogRequest bool

	// Whether to trust `X-Forwarded-Scheme:` and RFC 7239
	// `Forwarded: proto=`
	TrustForwarded bool

	MethodNotAllowed func(request Request, u *url.URL) Response

	outsideHandler func(Request) Response
	insideHandler  func(Request, Entity) Response
}

////////////////////////////////////////////////////////////////////////////////

// Request is an incoming HTTP request to be handled.
type Request struct {
	Method  string
	URL     *url.URL
	Headers http.Header
	Entity  interface{}
	Things  map[string]interface{}  // Objects added by middlewares
	cookies map[string]*http.Cookie // cached
}

// The Response to an HTTP request.  Create it using the appropriate
// (*Request).StatusDESCRIPTION method.
//
// That is; StatusSomething helper methods exist off of the request
// that you get passed.
type Response struct {
	Status  int16
	Headers http.Header
	Entity  NetEntity
}

// An Entity is some resource that is accessible over HTTP.
type Entity interface {
	// Methods() returns a map of HTTP request methods to Handlers
	// that handle requests for this Entity.
	Methods() map[string]func(Request) Response
}

// 404 Not Found
// 405 Method Not Allowed
// 406 Not Acceptable
// 400 Bad Request
// 500 Internal Server Error

// An EntityGroup is an Entity that also has child entities (i.e., a
// directory or folder).
type EntityGroup interface {
	Entity

	// Subentity(name, request) returns the child of this entity
	// with the name `name`, or nil if a child with that name
	// doesn't exist.
	//
	// The Request is included in the function call so that it can
	// be determined if the user has permission to access that
	// child.
	Subentity(name string, request Request) Entity

	// SubentityNotFound is called if Subentity returns nil.  If
	// the name contains a slash; it indicates that the child was
	// found, but a grandchild was requested, and the child wasn't
	// an EntityGroup.
	SubentityNotFound(name string, request Request) Response
}

// EntityExtra is an Entity that provides a custom HTTP 405 ("Method
// Not Allowed") handler.
type EntityExtra interface {
	Entity
	MethodNotAllowed(request Request) Response
}

// RootEntity is the union of EntityGroup and EntityExtra.
type RootEntity interface {
	Entity
	Subentity(name string, request Request) Entity
	SubentityNotFound(name string, request Request) Response
	MethodNotAllowed(request Request) Response
}

////////////////////////////////////////////////////////////////////////////////

// A NetEntity is just something that is capable of being transmitted
// over the network (in a variety of formats).
type NetEntity interface {
	// Encoders() returns a map of MIME-types to encoders that
	// serialize the NetEntity to that type.
	Encoders() map[string]func(io.Writer) error
}

////////////////////////////////////////////////////////////////////////////////

// A Middleware is something that wraps the request handler.
type Middleware struct {
	// Outside is able to affect the entity that is looked up
	Outside func(Request, func(Request) Response) Response
	// Inside cannot affect the entity that is looked up, but it
	// gets to inspect the entity.
	Inside func(Request, Entity, func(Request, Entity) Response) Response
}
