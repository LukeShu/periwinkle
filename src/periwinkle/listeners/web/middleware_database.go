// Copyright 2015 Luke Shumaker

package web

import (
	he "httpentity"
	"periwinkle/store"
	"database/sql"
)

type database struct{}

func (p database) Before(req *he.Request) {
	var transaction *sql.Tx = nil /* TODO */
	req.Things["db"] = store.DB(transaction)
}

func (p database) After(req he.Request, res *he.Response) {
	transaction := req.Things["db"].(*sql.Tx)
	err := transaction.Commit()
	if err != nil {
		// TODO: handle the error; it could be either HTTP 500
		// (Internal Server Error) or 409 (Conflict)
	}
}
