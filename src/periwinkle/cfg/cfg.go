// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package cfg

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"maildir"
	"net/http"
)

const IncomingMail maildir.Maildir = "/srv/periwinkle/Maildir"
const WebUiDir http.Dir = "./www"
const WebAddr string = ":8080"
const Debug bool = true

var DB *gorm.DB = getConnection()

func getConnection() *gorm.DB {
	db, err := gorm.Open("mysql", "periwinkle:periwinkle@/periwinkle?charset=utf8&parseTime=true")
	if err != nil {
		db, err = gorm.Open("sqlite3", "/temp/gorm.db")
		if err != nil {
			panic(err)
		}
	}
	db.LogMode(true)
	return &db
}
