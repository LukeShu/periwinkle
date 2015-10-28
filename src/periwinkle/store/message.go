// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"io"
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
	table := db.CreateTable(&o)
	table.AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	table.AddUniqueIndex("filename_idx", "unique")
}

func NewMessage(unique maildir.Unique) *Message {
	panic("TODO: SMTP+ORM: NewMessage()")
	// TODO: sprint2: add the message to the outgoing mail queue
}

func GetMessageById(db *gorm.DB, id string) *Message {
	var o Message
	if result := db.First(&o, "id = ?", id); result.Error != nil {
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
	return GetMessageById(db, name)
}
