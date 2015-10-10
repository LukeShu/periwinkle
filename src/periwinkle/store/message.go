// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"maildir"
)

var _ he.Entity = &Message{}
var _ he.NetEntity = &Message{}
var dirMessages he.Entity = newDirMessages()

// Model /////////////////////////////////////////////////////////////

type Message struct {
	Id      string
	GroupId string
	Unique  maildir.Unique
	// cached fields??????
}

func (o Message) schema(db *gorm.DB) {
	db.CreateTable(&o).
		AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT").
		AddUniqueIndex("filename_idx", "unique")
}

func NewMessage(unique maildir.Unique) *Message {
	panic("TODO: SMTP+ORM: NewMessage()")
	// TODO: sprint2: add the message to the outgoing mail queue
}

func GetMessageById(db *gorm.DB, id string) *Message {
	var o Message
	if result := db.First(&o, id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func (o *Message) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: SMTP: (*Message).Subentity()")
}

func (o *Message) Methods() map[string]he.Handler {
	return map[string]he.Handler{
		"GET": func(req he.Request) he.Response {
			// TODO: permission check
			return req.StatusOK(o)
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
	db := req.Things["db"].(*gorm.DB)
	return GetMessageById(db, name)
}
