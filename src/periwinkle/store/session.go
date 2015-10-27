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
			if sess == nil {
				return ret.StatusOK(make(map[string]interface{}))
			} else {
				return ret.StatusOK(sess)
			}
		},
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			badbody := req.StatusBadRequest(heutil.NetString("submitted body not what expected"))
			hash, ok := req.Entity.(map[string]interface{}); if !ok { return badbody }
			username, ok := hash["username"].(string)      ; if !ok { return badbody }
			password, ok := hash["password"].(string)      ; if !ok { return badbody }
			if len(hash) != 2                                       { return badbody }

			var user *User
			if strings.Contains(username, "@") {
				user = GetUserByAddress(db, "email", username)
			} else {
				user = GetUserById(db, username)
			}

			sess := NewSession(db, user, password)
			if sess == nil {
				return req.StatusForbidden(heutil.NetString("Incorrect username/password"))
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
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*Session)
			if sess != nil {
				sess.Delete(db)
			}
			return req.StatusNoContent()
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
