// Copyright 2015 Luke Shumaker

package main

import (
	he "httpentity"
	"net/url"
	"periwinkle/cfg"
	"periwinkle/util" // putil
)

func MiddlewareDatabase(req he.Request, u *url.URL, handle func(he.Request, *url.URL) he.Response, config cfg.Config) (res he.Response) {
	transaction := config.DB.Begin()
	req.Things["db"] = transaction
	rollback := true
	defer func() {
		if obj := recover(); obj != nil {
			if rollback {
				transaction.Rollback()
			}
			if err, ok := obj.(error); ok {
				perror := putil.ErrorToError(err)
				if perror.HttpCode() != 500 {
					res = putil.ErrorToHTTP(perror)
					return
				}
			}
			// we didn't intercept the error, so pass it along
			panic(obj)
		}
	}()

	res = handle(req, u)

	err := transaction.Commit().Error
	rollback = false
	if err != nil {
		panic(err)
	}

	return
}
