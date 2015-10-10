// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"time"
	//"github.com/dchest/captcha"
)

var _ he.Entity = &Captcha{}
var _ he.NetEntity = &Captcha{}
var dirCaptchas he.Entity = newDirCaptchas()

// Model /////////////////////////////////////////////////////////////

type Captcha struct {
	Id         string
	Value      string
	Token      string
	Expiration time.Time
}

func (o Captcha) schema(db *gorm.DB) {
	db.CreateTable(&o)
}

func NewCaptcha() *Captcha {
	panic("TODO: captcha+ORM: NewCaptcha()")
}

func GetCaptchaById(db *gorm.DB, id string) *Captcha {
	var o Captcha
	if result := db.First(&o, id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func (o *Captcha) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *Captcha) Methods() map[string]he.Handler {
	return map[string]he.Handler{
		"GET": func(req he.Request) he.Response {
			return req.StatusOK(o)
		},
		"PATCH": func(req he.Request) he.Response {
			panic("TODO: API: (*Captcha).Methods()[\"PATCH\"]")
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *Captcha) Encoders() map[string]he.Encoder {
	return defaultEncoders(o)
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirCaptchas struct {
	methods map[string]he.Handler
}

func newDirCaptchas() t_dirCaptchas {
	r := t_dirCaptchas{}
	r.methods = map[string]he.Handler{
		"POST": func(req he.Request) he.Response {
			return req.StatusCreated(r, NewCaptcha().Id)
		},
	}
	return r
}

func (d t_dirCaptchas) Methods() map[string]he.Handler {
	return d.methods
}

func (d t_dirCaptchas) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return GetCaptchaById(db, name)
}
