// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package cfg

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"periwinkle"
	"periwinkle/domain_handlers"
	"postfixpipe"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	yaml "gopkg.in/yaml.v2"
)

func Parse(in io.Reader) (*periwinkle.Cfg, error) {
	cfg := periwinkle.Cfg{}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	// sets defaults if file is empty
	if string(b) == "" {
		cfg.Mailstore = "/srv/periwinkle/Maildir"
		cfg.WebUiDir = "./www"
		cfg.Debug = true
		cfg.TrustForwarded = true
		cfg.GroupDomain = "periwinkle.lol"
	} else {
		err = yaml.Unmarshal(b, &cfg)
		if err != nil {
			return nil, err
		}
	}
	cfg.TwilioAccountId = os.Getenv("TWILIO_ACCOUNTID")
	cfg.TwilioAuthToken = os.Getenv("TWILIO_TOKEN")
	cfg.DB = getConnection(cfg.Debug) // TODO

	domain_handlers.GetHandlers(&cfg)
	cfg.DefaultDomainHandler = bounceNoHost

	return &cfg, err
}

func bounceNoHost(io.Reader, string, *gorm.DB, *periwinkle.Cfg) postfixpipe.ExitStatus {
	return postfixpipe.EX_NOHOST
}

func getConnection(debug bool) *gorm.DB {
	db, err := gorm.Open("mysql", "periwinkle:periwinkle@/periwinkle?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Falling back to SQLite3...")
		// here to change database load into memory
		db, err = gorm.Open("sqlite3", "file:periwinkle.sqlite?cache=shared&mode=rwc")
		if err != nil {
			panic(err)
		}
		db.DB().SetMaxOpenConns(1)
	}
	db.LogMode(debug)
	return &db
}
