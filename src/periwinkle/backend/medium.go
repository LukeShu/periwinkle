// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend

import (
	"locale"

	"github.com/jinzhu/gorm"
)

type Medium struct {
	ID string
}

func (o Medium) dbSchema(db *gorm.DB) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func (o Medium) dbSeed(db *gorm.DB) locale.Error {
	errs := errorList{}
	errHelper(&errs, locale.UntranslatedError(db.Create(&Medium{"email"}).Error))
	errHelper(&errs, locale.UntranslatedError(db.Create(&Medium{"twilio"}).Error))
	if len(errs) > 0 {
		return errs
	}
	return nil
}
