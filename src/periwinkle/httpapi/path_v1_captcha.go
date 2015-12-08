// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"io"
	"locale"
	"periwinkle/backend"
	"time"

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

		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			type postfmt struct {
				Value      string    `json:"value"`
				Expiration time.Time `json:"password"`
			}
			var entity postfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return *httperr
			}

			o := (*captcha)(backend.NewCaptcha(db))

			if o == nil {
				return rfc7231.StatusForbidden(he.NetPrintf("Captcha generation failed"))
			} else {
				ret := rfc7231.StatusOK(o)
				return ret
			}
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

func (o *captcha) IsText() bool {
	return false
}

func (o *captcha) Locales() []locale.Spec {
	return []locale.Spec{}
}

type captchaPNG struct{ *captcha }

func (o captchaPNG) Write(w io.Writer, l locale.Spec) locale.Error {
	return o.backend().MarshalPNG(w)
}

type captchaWAV struct{ *captcha }

func (o captchaWAV) Write(w io.Writer, l locale.Spec) locale.Error {
	return o.backend().MarshalWAV(w)
}

func (o *captcha) Encoders() map[string]he.Encoder {
	return map[string]he.Encoder{
		"image/png":      captchaPNG{o},
		"audio/vnd.wave": captchaWAV{o},
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
