// Copyright 2015 Luke Shumaker

package web

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3" // sqlite3
	he "httpentity"
	"httpentity/util"
	"periwinkle/cfg"
)

type database struct{}

func (p database) Before(req *he.Request) {
	transaction := cfg.DB.Begin()
	req.Things["db"] = transaction
}

func (p database) After(req he.Request, res *he.Response) {
	defer func() {
		if obj := recover(); obj != nil {
			switch err := obj.(type) {
			case sqlite3.Error:
				if err.Code == sqlite3.ErrConstraint {
					*res = he.StatusConflict(heutil.NetString(err.Error()))
					return
				}
			case mysql.MySQLError:
				// TODO: detect constraint falure for MySQL
			}
			// we didn't intercept the error, so pass it along
			panic(obj)
		}
	}()

	if obj := recover(); obj != nil {
		panic(obj)
	}

	transaction := req.Things["db"].(*gorm.DB)
	err := transaction.Commit().Error
	if err != nil {
		panic(err)
	}
}
