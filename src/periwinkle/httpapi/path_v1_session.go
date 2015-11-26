// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"httpentity/heutil"
	"httpentity/rfc7231"
	"io"
	"net/http"
	"periwinkle/backend"
	"strings"

	"github.com/jinzhu/gorm"
)

var _ he.NetEntity = &session{}
var _ he.Entity = &fileSession{}

type session backend.Session

func (o *session) backend() *backend.Session { return (*backend.Session)(o) }

// View //////////////////////////////////////////////////////////////

func (sess *session) Encoders() map[string]func(io.Writer) error {
	return defaultEncoders(sess)
}

// File ("Controller") ///////////////////////////////////////////////

type fileSession struct {
	methods map[string]func(he.Request) he.Response
}

func newFileSession() fileSession {
	r := fileSession{}
	r.methods = map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			sess := req.Things["session"].(*backend.Session)
			return rfc7231.StatusOK((*session)(sess))
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
				user = backend.GetUserByID(db, entity.Username)
			}

			sess := (*session)(backend.NewSession(db, user, entity.Password))
			if sess == nil {
				return rfc7231.StatusForbidden(heutil.NetString("Incorrect username/password"))
			} else {
				ret := rfc7231.StatusOK(sess)
				cookie := &http.Cookie{
					Name:     "session_id",
					Value:    sess.ID,
					Secure:   req.URL.Scheme == "https",
					HttpOnly: req.URL.Scheme == "http",
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
			return rfc7231.StatusNoContent()
		},
	}
	return r
}

func (d fileSession) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d fileSession) Subentity(name string, request he.Request) he.Entity {
	return nil
}
