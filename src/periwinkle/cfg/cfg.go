// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package cfg

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"maildir"
	"net/http"
	"os"
)

const IncomingMail maildir.Maildir = "/srv/periwinkle/Maildir"
const WebUiDir http.Dir = "./www"
const WebAddr string = ":8080"
const Debug bool = true

var DB *gorm.DB = getConnection()

func getConnection() *gorm.DB {
	db, err := gorm.Open("mysql", "periwinkle:periwinkle@/periwinkle?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Falling back to SQLite3\n")
		db, err = gorm.Open("sqlite3", "gorm.db")
		if err != nil {
			panic(err)
		}
	}
	db.LogMode(true)
	return &db
}
