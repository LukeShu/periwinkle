// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"httpentity/heutil"
	"io"
	"net/http"
	"strings"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

var _ he.NetEntity = &Session{}
var fileSession he.Entity = newFileSession()

type Session backend.Session

func (o *Session) backend() *backend.Session { return (*backend.Session)(o) }

// View //////////////////////////////////////////////////////////////

func (sess *Session) Encoders() map[string]func(io.Writer) error {
	return defaultEncoders(sess)
}

// File ("Controller") ///////////////////////////////////////////////

type t_fileSession struct {
	methods map[string]func(he.Request) he.Response
}

func newFileSession() t_fileSession {
	r := t_fileSession{}
	r.methods = map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			sess := req.Things["session"].(*backend.Session)
			return he.StatusOK((*Session)(sess))
		},
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			type postfmt struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			var entity postfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return *httperr
			}

			var user *backend.User
			if strings.Contains(entity.Username, "@") {
				user = backend.GetUserByAddress(db, "email", entity.Username)
			} else {
				user = backend.GetUserById(db, entity.Username)
			}

			sess := (*Session)(backend.NewSession(db, user, entity.Password))
			if sess == nil {
				return he.StatusForbidden(heutil.NetString("Incorrect username/password"))
			} else {
				ret := he.StatusOK(sess)
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
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			if sess != nil {
				sess.Delete(db)
			}
			return he.StatusNoContent()
		},
	}
	return r
}

func (d t_fileSession) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_fileSession) Subentity(name string, request he.Request) he.Entity {
	return nil
}
