// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"database/sql"
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

func GetMessageById(con DB, id string) *Message {
	var mes Message
	err := con.QueryRow("select * from message where id=?", id).Scan(&mes)
	switch {
	case err == sql.ErrNoRows:
		// message does not exist
		return nil
	case err != nil:
		// error talking to the DB
		panic(err)
	default:
		return &mes
	}
}

func (o *Message) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: SMTP: (*Message).Subentity()")
}

func (o *Message) Methods() map[string]he.Handler {
	return map[string]he.Handler{
		"GET": func(he.Request) he.Response {
			panic("TODO: API: (*Message).Subentity()")
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *Message) Encoders() map[string]he.Encoder {
	panic("TODO: API: (*Message).Encoders()")
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

func (d t_dirMessages) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(DB)
	return GetMessageById(db, name)
}
