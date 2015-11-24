// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package httpapi

import (
	he "httpentity"
	"net/url"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &shortURL{}
var _ he.Entity = &dirShortURLs{}

type shortURL backend.ShortURL

func (o *shortURL) backend() *backend.ShortURL { return (*backend.ShortURL)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *shortURL) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *shortURL) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			u, _ := url.Parse(o.Dest) // TODO: automatic unmarshal
			return he.StatusMovedPermanently(u)
		},
	}
}

// Directory ("Controller") //////////////////////////////////////////

type dirShortURLs struct {
	methods map[string]func(he.Request) he.Response
}

func NewDirShortURLs() he.Entity {
	return &dirShortURLs{
		methods: map[string]func(he.Request) he.Response{},
	}
}

func (d dirShortURLs) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d dirShortURLs) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return (*shortURL)(backend.GetShortURLByID(db, name))
}
