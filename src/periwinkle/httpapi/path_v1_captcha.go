// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal

package httpapi

import (
	he "httpentity"
	"io"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &Captcha{}
var _ he.NetEntity = &Captcha{}
var dirCaptchas he.Entity = newDirCaptchas()

type Captcha backend.Captcha

func (o *Captcha) backend() *backend.Captcha { return (*backend.Captcha)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *Captcha) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *Captcha) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return he.StatusOK(o)
		},
		"PUT": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			var newCaptcha Captcha
			httperr := safeDecodeJSON(req.Entity, &newCaptcha)
			if httperr != nil {
				return *httperr
			}
			*o = newCaptcha
			o.backend().Save(db)
			return he.StatusOK(o)
		},
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
			return o.backend().MarshalPNG(w)
		},
		"audio/vnd.wave": func(w io.Writer) error {
			return o.backend().MarshalWAV(w)
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
			return he.StatusCreated(r, backend.NewCaptcha(db).ID, req)
		},
	}
	return r
}

func (d t_dirCaptchas) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirCaptchas) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return (*Captcha)(backend.GetCaptchaByID(db, name))
}
