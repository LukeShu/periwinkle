// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
)

type Medium struct {
	Id string
}

func (o Medium) schema(db *gorm.DB) {
	db.CreateTable(&o)
}

func GetMedium(db *gorm.DB, id string) *Medium {
	var o Medium
	if result := db.First(&o, id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}
