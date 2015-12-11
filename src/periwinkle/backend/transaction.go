// Copyright 2015 Luke Shumaker

package backend

import (
	"locale"

	"github.com/jinzhu/gorm"
)

func WithTransaction(db *gorm.DB, fn func(tx *gorm.DB)) (conflict locale.Error) {
	transaction := db.Begin()
	rollback := true

	defer func() {
		if obj := recover(); obj != nil {
			if rollback {
				transaction.Rollback()
			}
			switch err := obj.(type) {
			case Conflict:
				conflict = err
			default:
				panic(obj)
			}
		}
	}()

	fn(transaction)

	err := transaction.Commit().Error
	rollback = false
	if err != nil {
		panic(err)
	}

	return
}
