// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	//"database/sql"
	he "httpentity"
	"time"
)

var _ he.NetEntity = &Session{}
var fileSession he.Entity = newFileSession()

// Model /////////////////////////////////////////////////////////////

type Session struct {
	id        string
	user_id   string
	last_used time.Time
}

func NewSession(username string, password string) *Session {
	panic("not implemented")
}

func GetSessionById(id string) *Session {
	panic("not implemented")
}

func (o *Session) Delete() {
	panic("not implemented")
}

// View //////////////////////////////////////////////////////////////

func (o *Session) Encoders() map[string]he.Encoder {
	panic("not implemented")
}

// File ("Controller") ///////////////////////////////////////////////

type t_fileSession struct {
	methods map[string]he.Handler
}

func newFileSession() t_fileSession {
	r := t_fileSession{}
	r.methods = map[string]he.Handler{
		"POST": func(req he.Request) he.Response {
			username := "" /*TODO*/
			password := "" /*TODO*/
			sess := NewSession(username, password)
			if sess == nil {
				return req.StatusUnauthorized(he.NetString("Incorrect username/password"))
			} else {
				return req.StatusOK(sess)
			}
		},
		"DELETE": func(req he.Request) he.Response {
			sess := GetSessionById("" /*TODO*/)
			if sess != nil {
				sess.Delete()
			}
			return req.StatusNoContent()
		},
	}
	return r
}

func (d t_fileSession) Methods() map[string]he.Handler {
	return d.methods
}

func (d t_fileSession) Subentity(name string, request he.Request) he.Entity {
	return nil
}
