// Copyright 2015 Luke Shumaker

package main

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"periwinkle/cfg"
	"periwinkle/util" // putil
)

type database struct{}

func (p database) Before(req *he.Request) {
	transaction := cfg.DB.Begin()
	req.Things["db"] = transaction
}

func (p database) After(req he.Request, res *he.Response) {
	transaction := req.Things["db"].(*gorm.DB)

	defer func() {
		if obj := recover(); obj != nil {
			if err, ok := obj.(error); ok {
				perror := putil.ErrorToError(err)
				if perror.HttpCode() != 500 {
					*res = putil.ErrorToHTTP(perror)
					return
				}
			}
			// we didn't intercept the error, so pass it along
			panic(obj)
		}
	}()

	if obj := recover(); obj != nil {
		transaction.Rollback()
		panic(obj)
	}

	err := transaction.Commit().Error
	if err != nil {
		panic(err)
	}
}
