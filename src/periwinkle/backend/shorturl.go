// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package backend

import (
	"locale"
	"net/url"
	"periwinkle"
)

type ShortURL struct {
	ID   string
	Dest string //*url.URL // TODO: figure out how to have (un)marshalling happen automatically
}

func (o ShortURL) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).
		AddUniqueIndex("dest_idx", "dest").
		Error)
}

func NewShortURL(db *periwinkle.Tx, u *url.URL) *ShortURL {
	o := ShortURL{
		ID:   randomString(5),
		Dest: u.String(), // TODO: automatic marshalling
	}
	if err := db.Create(&o).Error; err != nil {
		dbError(err)
	}
	return &o
}

func (o *ShortURL) Save(db *periwinkle.Tx) {
	if err := db.Save(o).Error; err != nil {
		dbError(err)
	}
}

func GetShortURLByID(db *periwinkle.Tx, id string) *ShortURL {
	var o ShortURL
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return &o
}
