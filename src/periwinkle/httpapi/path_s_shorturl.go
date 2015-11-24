// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package httpapi

import (
	he "httpentity"
	"net/url"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &ShortUrl{}
var DirShortUrls he.Entity = newDirShortUrls()

type ShortUrl backend.ShortUrl

func (o *ShortUrl) backend() *backend.ShortUrl { return (*backend.ShortUrl)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *ShortUrl) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *ShortUrl) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			u, _ := url.Parse(o.Dest) // TODO: automatic unmarshal
			return he.StatusMovedPermanently(u)
		},
	}
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirShortUrls struct {
	methods map[string]func(he.Request) he.Response
}

func newDirShortUrls() t_dirShortUrls {
	r := t_dirShortUrls{}
	r.methods = map[string]func(he.Request) he.Response{}
	return r
}

func (d t_dirShortUrls) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirShortUrls) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return (*ShortUrl)(backend.GetShortUrlById(db, name))
}
