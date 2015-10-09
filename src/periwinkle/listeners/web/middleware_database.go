// Copyright 2015 Luke Shumaker

package web

import (
	he "httpentity"
	"periwinkle/cfg"
	"github.com/jmoiron/modl"
)

type database struct{}

func (p database) Before(req *he.Request) {
	var transaction modl.SqlExecutor
	var err error
	transaction, err = cfg.DB.Begin()
	if transaction != nil && err == nil {
		req.Things["db"] = transaction
	}
}

func (p database) After(req he.Request, res *he.Response) {
	transaction, ok := req.Things["db"].(*modl.Transaction)
	if !ok {
		return
	}
	err := transaction.Commit()
	if err != nil {
		// TODO: DB: handle the error; it could be either HTTP 500
		// (Internal Server Error) or 409 (Conflict)
	}
}
