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
	Id        string
	User_id   string
	Last_used time.Time
}

func NewSession(con DB, username string, password string) *Session {
	panic("not implemented")
}

func GetSessionById(con DB, id string) *Session {
	panic("not implemented")
}

func (o *Session) Delete(con DB) {
	panic("not implemented")
}

func (o *Session) Save(con DB) {
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
			db := req.Things["db"].(DB)
			badbody := req.StatusBadRequest("submitted body not what expected")
			hash    , ok := req.Entity.(map[string]interface{}); if !ok { return badbody }
			username, ok := hash["username"].(string)          ; if !ok { return badbody }
			password, ok := hash["password"].(string)          ; if !ok { return badbody }
			if len(hash) != 2                                           { return badbody }

			sess := NewSession(db, username, password)
			if sess == nil {
				return req.StatusUnauthorized(he.NetString("Incorrect username/password"))
			} else {
				ret := req.StatusOK(sess)
				// TODO: set the session_id cookie (in ret.Headers) to sess.Id
				return ret
			}
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(DB)
			sess := req.Things["session"].(*Session)
			if sess != nil {
				sess.Delete(db)
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
