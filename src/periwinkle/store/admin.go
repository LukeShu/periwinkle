// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
)

type Admin struct {
	UserId string `json:"user_id"`
}

func (o Admin) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}
