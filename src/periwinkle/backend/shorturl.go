// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package backend

import (
	"net/url"

	"github.com/jinzhu/gorm"
)

type ShortUrl struct {
	Id   string
	Dest string //*url.URL // TODO: figure out how to have (un)marshalling happen automatically
}

func (o ShortUrl) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddUniqueIndex("dest_idx", "dest").
		Error
}

func NewShortURL(db *gorm.DB, u *url.URL) *ShortUrl {
	o := ShortUrl{
		Id:   randomString(5),
		Dest: u.String(), // TODO: automatic marshalling
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
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}
