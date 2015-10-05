// Copyright 2015 Luke Shumaker

package store

import (
	//"database/sql"
	he "httpentity"
	"net/url"
)

var _ he.Entity = &ShortUrl{}
var dirShortUrls he.Entity = newDirShortUrls()

// Model /////////////////////////////////////////////////////////////

type ShortUrl struct {
	Id   string
	Dest *url.URL
}

func GetShortUrlById(con DB, id string) *ShortUrl {
	panic("not implemented")
}

func (o *ShortUrl) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *ShortUrl) Methods() map[string]he.Handler {
	panic("not implemented")
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirShortUrls struct {
	methods map[string]he.Handler
}

func newDirShortUrls() t_dirShortUrls {
	r := t_dirShortUrls{}
	r.methods = map[string]he.Handler{}
	return r
}

func (d t_dirShortUrls) Methods() map[string]he.Handler {
	return d.methods
}

func (d t_dirShortUrls) Subentity(name string, request he.Request) he.Entity {
	return GetShortUrlById(nil /*TODO*/, name)
}
