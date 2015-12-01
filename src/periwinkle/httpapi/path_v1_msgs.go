// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &message{}
var _ he.NetEntity = &message{}
var _ he.EntityGroup = dirMessages{}

type message backend.Message

func (o *message) backend() *backend.Message { return (*backend.Message)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *message) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: SMTP: (*message).Subentity()")
}

func (o *message) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			// TODO: permission check
			return rfc7231.StatusOK(o)
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *message) Encoders() map[string]he.Encoder {
	panic("TODO: API: (*message).Encoders()")
}

func (o *message) SubentityNotFound(name string, req he.Request) he.Response {
	panic("TODO: SMTP: (*message).SubentityNotFound()")
}

// Directory ("Controller") //////////////////////////////////////////

type dirMessages struct {
	methods map[string]func(he.Request) he.Response
}

func newDirMessages() dirMessages {
	r := dirMessages{}
	r.methods = map[string]func(he.Request) he.Response{}
	return r
}

func (d dirMessages) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d dirMessages) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return (*message)(backend.GetMessageByID(db, name))
}

func (d dirMessages) SubentityNotFound(name string, req he.Request) he.Response {
	return rfc7231.StatusNotFound(nil)
}
