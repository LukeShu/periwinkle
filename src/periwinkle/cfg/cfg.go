// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package cfg

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"maildir"
	"net/http"
	"os"
)

const IncomingMail maildir.Maildir = "/srv/periwinkle/Maildir"
const WebUiDir http.Dir = "./www"
const Debug bool = true
const TrustForwarded = true // whether to trust X-Forwarded: or Forwarded: HTTP headers

var TwilioAccountId = os.Getenv("TWILIO_ACCOUNTID")
var TwilioAuthToken = os.Getenv("TWILIO_TOKEN")

var GroupDomain = "periwinkle.lol"

var WebRoot = getWebroot()

func getWebroot() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname + ":8080"
}

var DB *gorm.DB = getConnection()

func getConnection() *gorm.DB {
	db, err := gorm.Open("mysql", "periwinkle:periwinkle@/periwinkle?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Falling back to SQLite3\n")
		db, err = gorm.Open("sqlite3", "file:periwinkle.sqlite?cache=shared&mode=rwc")
		if err != nil {
			panic(err)
		}
		db.DB().SetMaxOpenConns(1)
	}
	db.LogMode(Debug)
	return &db
}

type DomainHandler func(io.Reader, string, *gorm.DB) int

var DomainHandlers map[string]DomainHandler // set in email_handlers/init.go because import-cycles

var DefaultDomainHandler DomainHandler = bounce

func bounce(io.Reader, string, *gorm.DB) int {
	return 1
}
