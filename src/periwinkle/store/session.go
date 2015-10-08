// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	//"database/sql"
	he "httpentity"
	"net/http"
	"time"
)

var _ he.NetEntity = &Session{}
var fileSession he.Entity = newFileSession()

// Model /////////////////////////////////////////////////////////////

type Session struct {
	Id       string
	UserId   string
	LastUsed time.Time
}

func NewSession(con DB, username string, password string) *Session {
	user := GetUserByName(con, username)
	if !user.CheckPassword(password) {
		return nil
	}

	ses := &Session{
		Id:       randomString(24),
		UserId:   user.Id,
		LastUsed: time.Now(),
	}
	return ses
}

func GetSessionById(con DB, id string) *Session {
	panic("TODO: ORM: GetSessionById()")
}

func (o *Session) Delete(con DB) {
	panic("TODO: ORM: (*Session).Delete()")
}

func (o *Session) Save(con DB) {
	panic("TODO: ORM: (*Session).Save()")
}

// View //////////////////////////////////////////////////////////////

func (sess *Session) Encoders() map[string]he.Encoder {
	dat := map[string]string{
		"session_id": sess.Id,
	}
	return defaultEncoders(dat)
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
			hash, ok := req.Entity.(map[string]interface{}); if !ok { return badbody }
			username, ok := hash["username"].(string)      ; if !ok { return badbody }
			password, ok := hash["password"].(string)      ; if !ok { return badbody }
			if len(hash) != 2                                       { return badbody }

			sess := NewSession(db, username, password)
			if sess == nil {
				return req.StatusUnauthorized(he.NetString("Incorrect username/password"))
			} else {
				ret := req.StatusOK(sess)
				cookie := &http.Cookie{
					Name:     "session_id",
					Value:    sess.Id,
					Secure:   req.Scheme == "https",
					HttpOnly: req.Scheme == "http",
				}
				ret.Headers.Add("Set-Cookie", cookie.String())
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
