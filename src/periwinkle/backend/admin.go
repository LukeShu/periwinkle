// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend

import (
	"locale"

	"github.com/jinzhu/gorm"
)

type Admin struct {
	UserID string `json:"user_id"`
}

func (o Admin) dbSchema(db *gorm.DB) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}
