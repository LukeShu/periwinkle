// Copyright 2015 Luke Shumaker

package web

import (
	"database/sql"
	he "httpentity"
	"periwinkle/cfg"
	"periwinkle/store"
)

type database struct{}

func (p database) Before(req *he.Request) {
	var transaction *sql.Tx
	var err error
	transaction, err = cfg.DB.Begin()
	if transaction != nil && err == nil {
		req.Things["db"] = store.DB(transaction)
	}
}

func (p database) After(req he.Request, res *he.Response) {
	transaction, ok := req.Things["db"].(*sql.Tx)
	if !ok {
		return
	}
	err := transaction.Commit()
	if err != nil {
		// TODO: DB: handle the error; it could be either HTTP 500
		// (Internal Server Error) or 409 (Conflict)
	}
}
