// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package store

import (
	"github.com/jinzhu/gorm"
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

func (o ShortUrl) schema(db *gorm.DB) {
	db.CreateTable(&o).
		AddUniqueIndex("dest_idx", "dest")
}

func NewShortURL(db *gorm.DB, u *url.URL) *ShortUrl {
	o := ShortUrl{
		Id:   randomString(5),
		Dest: u,
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func (o *ShortUrl) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}

func GetShortUrlById(db *gorm.DB, id string) *ShortUrl {
	var o ShortUrl
	if result := db.First(&o, id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func (o *ShortUrl) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *ShortUrl) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return req.StatusMoved(o.Dest)
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
	return GetShortUrlById(db, name)
}
