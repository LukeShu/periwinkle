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
			badbody := req.StatusBadRequest("submitted body not what expected")
			hash    , ok := req.Entity.(map[string]interface{}); if !ok { return badbody }
			username, ok := hash["username"].(string)          ; if !ok { return badbody }
			password, ok := hash["password"].(string)          ; if !ok { return badbody }
			if len(hash) != 2                                           { return badbody }

			sess := NewSession(username, password)
			if sess == nil {
				return req.StatusUnauthorized(he.NetString("Incorrect username/password"))
			} else {
				ret := req.StatusOK(sess)
				// TODO: set the session_id cookie (in ret.Headers) to sess.Id
				return ret
			}
		},
		"DELETE": func(req he.Request) he.Response {
			sess := req.Things["session"].(*Session)
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
