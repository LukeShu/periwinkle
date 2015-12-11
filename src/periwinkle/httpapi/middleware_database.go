// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"periwinkle"
)

func MiddlewareDatabase(config *periwinkle.Cfg) he.Middleware {
	return he.Middleware{
		Outside: func(req he.Request, handle func(he.Request) he.Response) he.Response {
			var res he.Response
			conflict := config.DB.Do(func(transaction *periwinkle.Tx) {
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
