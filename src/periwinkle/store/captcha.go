// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal

package store

import (
	he "httpentity"
	"io"
	"time"

	"github.com/dchest/captcha"
	"github.com/jinzhu/gorm"
)

const (
	// Default number of digits in captcha solution.
	DefaultLen = 6
	// Expiration time of captchas used by default store.
	DefaultExpiration = 20 * time.Minute
	//	Default Captcha Image Width
	DefaultWidth = 640
	//	Default Captcha Image Height
	DefaultHeight = 480
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

func (o Captcha) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

func NewCaptcha(db *gorm.DB) *Captcha {
	o := Captcha{
		Id:    captcha.New(),
		Value: string(captcha.RandomDigits(DefaultLen)),
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func UseCaptcha(db *gorm.DB, token string) bool {
	panic("TODO")
}

func GetCaptchaById(db *gorm.DB, id string) *Captcha {
	var o Captcha
	if result := db.First(&o, "id = ?", id); result.Error != nil {
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

func (o *Captcha) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return he.StatusOK(o)
		},
		/*
			"PUT": func(req he.Request) he.Response {
				db := req.Things["db"].(*gorm.DB)
				sess := req.Things["session"].(*Session)
				var new_captcha Captcha
				httperr := safeDecodeJSON(req.Entity, &new_captcha)
				if httperr != nil {
					return *httperr
				}
			},
		*/
		/*
			"PATCH": func(req he.Request) he.Response {
				panic("TODO: API: (*Captcha).Methods()[\"PATCH\"]")
			},
		*/
	}
}

// View //////////////////////////////////////////////////////////////

func (o *Captcha) Encoders() map[string]func(io.Writer) error {
	return map[string]func(io.Writer) error{
		"image/png": func(w io.Writer) error {
			// TODO: generate PNG and write it to w
			return captcha.WriteImage(w, o.Id, DefaultWidth, DefaultHeight)
		},
		"audio/vnd.wave": func(w io.Writer) error {
			// TODO: generate WAV and write it to w
			return captcha.WriteAudio(w, o.Id, "en")
		},
	}
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirCaptchas struct {
	methods map[string]func(he.Request) he.Response
}

func newDirCaptchas() t_dirCaptchas {
	r := t_dirCaptchas{}
	r.methods = map[string]func(he.Request) he.Response{
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			return he.StatusCreated(r, NewCaptcha(db).Id, req)
		},
	}
	return r
}

func (d t_dirCaptchas) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirCaptchas) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return GetCaptchaById(db, name)
}
