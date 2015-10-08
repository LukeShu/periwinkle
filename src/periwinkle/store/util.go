// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	he "httpentity"
	"io"
	"math/big"
	"github.com/jmoiron/modl"
	"periwinkle/cfg"
)
// Global for database access in the ORM
var dbMap = &modl.DbMap{Db: cfg.DB, Dialect: modl.MySQLDialect{"InnoDB", "UTF8"}}

// The intersection of *sql.DB and *sql.Tx
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Simple dump to JSON, good for most entities
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

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var alphabetLen = big.NewInt(int64(len(alphabet)))

func randomString(size int) string {
	var randStr []byte
	for i := 0; i < size; i++ {
		bigint, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			panic(err)
		}
		randStr[i] = alphabet[bigint.Int64()]
	}
	return string(randStr[:])
}
