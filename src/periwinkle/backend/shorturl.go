// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package backend

import (
	"net/url"

	"github.com/jinzhu/gorm"
)

type ShortURL struct {
	ID   string
	Dest string //*url.URL // TODO: figure out how to have (un)marshalling happen automatically
}

func (o ShortURL) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddUniqueIndex("dest_idx", "dest").
		Error
}

func NewShortURL(db *gorm.DB, u *url.URL) *ShortURL {
	o := ShortURL{
		ID:   randomString(5),
		Dest: u.String(), // TODO: automatic marshalling
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func (o *ShortURL) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}

func GetShortURLByID(db *gorm.DB, id string) *ShortURL {
	var o ShortURL
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}
