// Copyright 2015 Luke Shumaker

package web

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"periwinkle/cfg"
)

type database struct{}

func (p database) Before(req *he.Request) {
	transaction := cfg.DB.Begin()
	req.Things["db"] = transaction
}

func (p database) After(req he.Request, res *he.Response) {
	transaction := req.Things["db"].(*gorm.DB)
	result := transaction.Commit()
	if result.Error != nil {
		// TODO: DB: handle the error; it could be either HTTP 500
		// (Internal Server Error) or 409 (Conflict)
	}
}
