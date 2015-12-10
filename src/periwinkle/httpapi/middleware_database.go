// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"periwinkle"
	"periwinkle/backend"

	"github.com/jinzhu/gorm"
)

func MiddlewareDatabase(config *periwinkle.Cfg) he.Middleware {
	return he.Middleware{
		Outside: func(req he.Request, handle func(he.Request) he.Response) he.Response {
			var res he.Response
			conflict := backend.WithTransaction(config.DB, func(transaction *gorm.DB) {
				req.Things["db"] = transaction
				res = handle(req)
			})
			if conflict != nil {
				res = rfc7231.StatusConflict(he.NetStringer{conflict})
			}
			return res
		},
	}
}
