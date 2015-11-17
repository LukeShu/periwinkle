// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package cfg

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"maildir"
	"net/http"
	"os"
//	"periwinkle/email_handlers"
	"postfixpipe"
)

type Cfg struct {
	Mailstore            maildir.Maildir          
	WebUiDir             http.Dir                 
	Debug                bool                     
	TrustForwarded       bool                     // whether to trust X-Forwarded: or Forwarded: HTTP headers
	TwilioAccountId      string                   
	TwilioAuthToken      string                 
	GroupDomain          string                  
	WebRoot              string                  
	DB                   *gorm.DB                
	DomainHandlers       map[string]DomainHandler
	DefaultDomainHandler DomainHandler
}


/*
  in periwinkle.conf
	figure out how to get the io.Reader to give all of the below information
*/
func Parse(in io.Reader) (*Cfg, error) {
	cfg := Cfg{}
	
	b, err := ioutil.ReadAll(in)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}
	/*
	cfg.Mailstore       = "/srv/periwinkle/Maildir"
	cfg.WebUiDir        = "./www"
	cfg.Debug           = true
	cfg.TrustForwarded  = true
	cfg.GroupDomain     = "periwinkle.lol"
	*/
	cfg.TwilioAccountId = os.Getenv("TWILIO_ACCOUNTID")
	cfg.TwilioAuthToken = os.Getenv("TWILIO_TOKEN")
	cfg.WebRoot = getWebroot()
	cfg.DB = getConnection()
	// cfg.DomainHandlers = email_handlers.GetHandlers()
	cfg.DefaultDomainHandler = bounceNoHost

	return &cfg, err
}

type DomainHandler func(io.Reader, string, *gorm.DB, Cfg) uint8

func bounceNoHost(io.Reader, string, *gorm.DB, Cfg) uint8 {
	return postfixpipe.EX_NOHOST
}

func getWebroot() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname + ":8080"
}

func getConnection() *gorm.DB {
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
	db.LogMode(true)
	return &db
}
