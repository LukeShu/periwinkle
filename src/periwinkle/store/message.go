// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	//"database/sql"
	he "httpentity"
)

var _ he.Entity = &Message{}
var _ he.NetEntity = &Message{}
var dirMessages he.Entity = newDirMessages()

// Model /////////////////////////////////////////////////////////////

type Message struct {
	id       string
	group_id int
	filename string
	// cached fields??????
}

func GetMessageById(id string) *Message {
	panic("not implemented")
}

func (o *Message) Subentity(name string, req he.Request) he.Entity {
	panic("not implemented")
}

func (o *Message) Methods() map[string]he.Handler {
	return map[string]he.Handler{
		"GET": func(he.Request) he.Response {
			panic("not implemented")
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *Message) Encoders() map[string]he.Encoder {
	panic("not implemented")
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirMessages struct {
	methods map[string]he.Handler
}

func newDirMessages() t_dirMessages {
	r := t_dirMessages{}
	r.methods = map[string]he.Handler{}
	return r
}

func (d t_dirMessages) Methods() map[string]he.Handler {
	return d.methods
}

func (d t_dirMessages) Subentity(name string, request he.Request) he.Entity {
	return GetMessageById(name)
}
