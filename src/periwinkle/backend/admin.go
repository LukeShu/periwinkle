// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend

import (
	"github.com/jinzhu/gorm"
)

type Admin struct {
	UserID string `json:"user_id"`
}

func (o Admin) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}
