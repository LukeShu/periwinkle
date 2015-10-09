// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"database/sql"
	"github.com/jmoiron/modl"
)

type Medium struct {
	Id string
}

func GetMedium(con *modl.Transaction, id string) *Medium {
	var med Medium
	err := con.Get(&med, id)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		return &med
	}
}
