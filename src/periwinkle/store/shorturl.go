// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package store

import (
	"database/sql"
	he "httpentity"
	"net/url"
	"github.com/jmoiron/modl"
)

var _ he.Entity = &ShortUrl{}
var dirShortUrls he.Entity = newDirShortUrls()

// Model /////////////////////////////////////////////////////////////

type ShortUrl struct {
	Id   string
	Dest *url.URL
}

func NewShortURL(con modl.SqlExecutor, u *url.URL) *ShortUrl {
	s := &ShortUrl{
		Id:   randomString(5),
		Dest: u,
	}
	err := con.Insert(s)
	if err != nil {
		return nil
	}
	return s
}

func (s *ShortUrl) Save(con modl.SqlExecutor) error {
	_, err := con.Update(s)
	return err
}

func GetShortUrlById(con modl.SqlExecutor, id string) *ShortUrl {
	var s ShortUrl
	err := con.Get(&s, id)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		return &s
	}
}

func (o *ShortUrl) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *ShortUrl) Methods() map[string]he.Handler {
	return map[string]he.Handler{
		"GET": func(req he.Request) he.Response {
			return req.StatusMoved(o.Dest)
		},
	}
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

func (d t_dirShortUrls) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(modl.SqlExecutor)
	return GetShortUrlById(db, name)
}
