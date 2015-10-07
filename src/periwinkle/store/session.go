// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	//"database/sql"
	he "httpentity"
	"time"
	"math/big"
	//"math/rand"
	"crypto/rand"
)

var _ he.NetEntity = &Session{}
var fileSession he.Entity = newFileSession()

// Model /////////////////////////////////////////////////////////////

type Session struct {
	Id        string
	UserId   string
	LastUsed time.Time
}

func NewSession(con DB, username string, password string) *Session {
	user := GetUserByName(con, username)
	if !user.CheckPassword(password) {
		return nil
	}

	ses := &Session{
		Id: createSessionId(),
		UserId:   user.Id,
		LastUsed: Now(),
	}
	return ses
}

func Now() time.Time{
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
			// the below var name is both funny and kinda hot
			badbody := req.StatusBadRequest("submitted body not what expected")
			hash, ok := req.Entity.(map[string]interface{})
			if !ok {
				return badbody
			}
			username, ok := hash["username"].(string)
			if !ok {
				return badbody
			}
			password, ok := hash["password"].(string)
			if !ok {
				return badbody
			}
			if len(hash) != 2 {
				return badbody
			}

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

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var alphabetLen = big.NewInt(int64(len(alphabet)))

func randomByte(size int) string {
		byteSize := size
		var randStr []byte
		for i := 0; i < size; i++ {
			bigint, err := rand.Int(rand.Reader, alphabetLen)
			if err != nil {
				panic(err)
			}
			randStr[i] = alphabet[bigint.Int64()]
		}
		return string(randStr[:])
}

func createSessionId() string {
	return randomByte(24)
}

