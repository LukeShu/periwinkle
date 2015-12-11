// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend

import (
	"locale"
	"periwinkle"
)

type Medium struct {
	ID string
}

func (o Medium) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func (o Medium) dbSeed(db *periwinkle.Tx) locale.Error {
	errs := errorList{}
	errHelper(&errs, locale.UntranslatedError(db.Create(&Medium{"email"}).Error))
	errHelper(&errs, locale.UntranslatedError(db.Create(&Medium{"sms"}).Error))
	errHelper(&errs, locale.UntranslatedError(db.Create(&Medium{"mms"}).Error))
	if len(errs) > 0 {
		return errs
	}
	return nil
}
