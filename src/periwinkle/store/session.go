// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"httpentity/util"
	"io"
	"net/http"
	"strings"
	"time"
)

var _ he.NetEntity = &Session{}
var fileSession he.Entity = newFileSession()

// Model /////////////////////////////////////////////////////////////

type Session struct {
	Id       string    `json:"session_id"`
	UserId   string    `json:"user_id"`
	LastUsed time.Time `json:"-"`
}

func (o Session) schema(db *gorm.DB) {
	table := db.CreateTable(&o)
	table.AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
}

func NewSession(db *gorm.DB, user *User, password string) *Session {
	if user == nil || !user.CheckPassword(password) {
		return nil
	}
	o := Session{
		Id:       randomString(24),
		UserId:   user.Id,
		LastUsed: time.Now(),
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func GetSessionById(db *gorm.DB, id string) *Session {
	var o Session
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func (o *Session) Delete(db *gorm.DB) {
	db.Delete(o)
}

func (o *Session) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}

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
			sess := req.Things["session"].(*Session)
			return he.StatusOK(sess)
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
				return httperr.Response()
			}

			var user *User
			if strings.Contains(entity.Username, "@") {
				user = GetUserByAddress(db, "email", entity.Username)
			} else {
				user = GetUserById(db, entity.Username)
			}

			sess := NewSession(db, user, entity.Password)
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
			sess := req.Things["session"].(*Session)
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
