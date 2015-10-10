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
