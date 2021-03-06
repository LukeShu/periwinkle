// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"net/url"
	"periwinkle"
	"periwinkle/backend"
)

var _ he.Entity = &shortURL{}

type shortURL backend.ShortURL

func (o *shortURL) backend() *backend.ShortURL { return (*backend.ShortURL)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *shortURL) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			u, _ := url.Parse(o.Dest) // TODO: automatic unmarshal
			return rfc7231.StatusMovedPermanently(u)
		},
	}
}

// Directory ("Controller") //////////////////////////////////////////

type dirShortURLs struct {
	methods map[string]func(he.Request) he.Response
}

func NewDirShortURLs() he.RootEntity {
	return &dirShortURLs{
		methods: map[string]func(he.Request) he.Response{},
	}
}

func (d dirShortURLs) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d dirShortURLs) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*periwinkle.Tx)
	return (*shortURL)(backend.GetShortURLByID(db, name))
}

func (d dirShortURLs) SubentityNotFound(name string, request he.Request) he.Response {
	panic("TODO")
}

func (d dirShortURLs) MethodNotAllowed(request he.Request) he.Response {
	panic("TODO")
}
