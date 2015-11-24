// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package httpapi

import (
	he "httpentity"
	"net/url"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &ShortURL{}
var DirShortURLs he.Entity = newDirShortURLs()

type ShortURL backend.ShortURL

func (o *ShortURL) backend() *backend.ShortURL { return (*backend.ShortURL)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *ShortURL) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *ShortURL) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			u, _ := url.Parse(o.Dest) // TODO: automatic unmarshal
			return he.StatusMovedPermanently(u)
		},
	}
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirShortURLs struct {
	methods map[string]func(he.Request) he.Response
}

func newDirShortURLs() t_dirShortURLs {
	r := t_dirShortURLs{}
	r.methods = map[string]func(he.Request) he.Response{}
	return r
}

func (d t_dirShortURLs) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirShortURLs) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return (*ShortURL)(backend.GetShortURLByID(db, name))
}
