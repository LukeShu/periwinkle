// Copyright 2015 Luke Shumaker

package store

import (
	"encoding/json"
	he "httpentity"
	"database/sql"
	"io"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func defaultEncoders(o interface{}) map[string]he.Encoder {
	return map[string]he.Encoder{
		"application/json": func(w io.Writer) error {
			bytes, err := json.Marshal(o)
			if err != nil {
				return err
			}
			_, err = w.Write(bytes)
			return err
		},
	}
}
