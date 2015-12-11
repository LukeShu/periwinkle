// Copyright 2015 Luke Shumaker

package periwinkle

import (
	"locale"

	"github.com/jinzhu/gorm"
)

type Conflict struct {
	Err locale.Error
}

func (c Conflict) Error() string {
	return c.Err.Error()
}

func (c Conflict) Locales() []locale.Spec {
	return c.Err.Locales()
}

func (c Conflict) L10NString(s locale.Spec) string {
	return c.Err.L10NString(s)
}

var _ locale.Error = Conflict{}

type Tx struct {
	*gorm.DB
}

type DB struct {
	inner gorm.DB
}

func NewDB(db gorm.DB) *DB {
	return &DB{db}
}

func (db *DB) Do(fn func(tx *Tx)) (conflict locale.Error) {
	transaction := db.inner.Begin()
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

	fn(&Tx{transaction})

	err := transaction.Commit().Error
	rollback = false
	if err != nil {
		panic(err)
	}

	return
}
