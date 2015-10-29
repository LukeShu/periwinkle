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
	transaction := req.Things["db"].(*gorm.DB)

	defer func() {
		if obj := recover(); obj != nil {
			switch err := obj.(type) {
			case sqlite3.Error:
				if err.Code == sqlite3.ErrConstraint {
					*res = he.StatusConflict(heutil.NetString(err.Error()))
					return
				}
			case *mysql.MySQLError:
				// TODO: this list of error numbers might not be complete
				// see https://mariadb.com/kb/en/mariadb/mariadb-error-codes/
				switch err.Number {
				case 1022, 1062, 1169, 1216, 1217, 1451, 1452, 1557, 1761, 1762, 1834:
					*res = he.StatusConflict(heutil.NetString(err.Error()))
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
