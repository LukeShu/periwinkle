// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker

package cfg

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"maildir"
	"net/http"
	"github.com/jmoiron/modl"
)

const IncomingMail maildir.Maildir = "/srv/periwinkle/Maildir"
const WebUiDir http.Dir = "./www"
const WebAddr string = ":8080"
const Debug bool = true

var DB *modl.DbMap = getConnection()

func getConnection() *modl.DbMap {
	sql, err := sql.Open("mysql", "periwinkle:periwinkle@/periwinkle?parseTime=true")
	if err != nil {
		panic("Could not connect to database")
	}
	dbMap := modl.NewDbMap(sql, modl.MySQLDialect{"InnoDB", "UTF8"})
	return dbMap
}
