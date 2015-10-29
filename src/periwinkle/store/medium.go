// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
)

type Medium struct {
	Id string
}

func (o Medium) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

func (o Medium) dbSeed(db *gorm.DB) error {
	return db.Create(&Medium{"email"}).Error
}
