// Copyright 2015 Luke Shumaker

package web

import (
	"database/sql"
	he "httpentity"
	"periwinkle/store"
)

type database struct{}

func (p database) Before(req *he.Request) {
	var transaction *sql.Tx = nil /* TODO: DB */
	req.Things["db"] = store.DB(transaction)
}

func (p database) After(req he.Request, res *he.Response) {
	transaction := req.Things["db"].(*sql.Tx)
	err := transaction.Commit()
	if err != nil {
		// TODO: DB: handle the error; it could be either HTTP 500
		// (Internal Server Error) or 409 (Conflict)
	}
}
