// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"io"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &Message{}
var _ he.NetEntity = &Message{}
var dirMessages he.Entity = newDirMessages()

type Message backend.Message

func (o *Message) backend() *backend.Message { return (*backend.Message)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *Message) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: SMTP: (*Message).Subentity()")
}

func (o *Message) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			// TODO: permission check
			return he.StatusOK(o)
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *Message) Encoders() map[string]func(io.Writer) error {
	panic("TODO: API: (*Message).Encoders()")
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirMessages struct {
	methods map[string]func(he.Request) he.Response
}

func newDirMessages() t_dirMessages {
	r := t_dirMessages{}
	r.methods = map[string]func(he.Request) he.Response{}
	return r
}

func (d t_dirMessages) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirMessages) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return (*Message)(backend.GetMessageById(db, name))
}
