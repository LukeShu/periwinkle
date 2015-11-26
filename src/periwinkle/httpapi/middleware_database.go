// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"periwinkle"
	"periwinkle/putil"
)

func MiddlewareDatabase(config *periwinkle.Cfg) he.Middleware {
	return he.Middleware{
		Outside: func(req he.Request, handle func(he.Request) he.Response) (res he.Response) {
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
						if perror.HTTPCode() != 500 {
							res = putil.ErrorToHTTP(perror)
							return
						}
					}
					// we didn't intercept the error, so pass it along
					panic(obj)
				}
			}()

			res = handle(req)

			err := transaction.Commit().Error
			rollback = false
			if err != nil {
				panic(err)
			}

			return
		},
	}
}
