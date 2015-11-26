// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"io"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &captcha{}
var _ he.NetEntity = &captcha{}
var _ he.EntityGroup = &dirCaptchas{}

type captcha backend.Captcha

func (o *captcha) backend() *backend.Captcha { return (*backend.Captcha)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *captcha) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return rfc7231.StatusOK(o)
		},
		"PUT": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			var newCaptcha captcha
			httperr := safeDecodeJSON(req.Entity, &newCaptcha)
			if httperr != nil {
				return *httperr
			}
			*o = newCaptcha
			o.backend().Save(db)
			return rfc7231.StatusOK(o)
		},
		/*
			"PATCH": func(req he.Request) he.Response {
				panic("TODO: API: (*captcha).Methods()[\"PATCH\"]")
			},
		*/
	}
}

// View //////////////////////////////////////////////////////////////

func (o *captcha) Encoders() map[string]func(io.Writer) error {
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

type dirCaptchas struct {
	methods map[string]func(he.Request) he.Response
}

func newDirCaptchas() dirCaptchas {
	r := dirCaptchas{}
	r.methods = map[string]func(he.Request) he.Response{
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			return rfc7231.StatusCreated(r, backend.NewCaptcha(db).ID, req)
		},
	}
	return r
}

func (d dirCaptchas) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d dirCaptchas) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return (*captcha)(backend.GetCaptchaByID(db, name))
}

func (d dirCaptchas) SubentityNotFound(name string, req he.Request) he.Response {
	return rfc7231.StatusNotFound(nil)
}
