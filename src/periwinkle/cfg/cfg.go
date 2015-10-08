// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker

package cfg

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"maildir"
	"net/http"
)

const IncomingMail maildir.Maildir = "/srv/periwinkle/Maildir"
const WebUiDir http.Dir = "./www"
const WebAddr string = ":8080"
const Debug bool = true

var DB *sql.DB = getConnection()

func getConnection() *sql.DB {
	db_user := "periwinkle"
	db_pass := "periwinkle"
	db_name := "periwinkle"

	db, err := sql.Open("mysql", db_user+":"+db_pass+"@/"+db_name)
	if err != nil {
		panic("Could not connect to database")
	}

	return db
}
