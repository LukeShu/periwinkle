// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package cfg

import (
	"fmt"
	"io"
	"io/ioutil"
	"maildir"
	"net/http"
	"os"
	"periwinkle"
	"periwinkle/domain_handlers"
	"postfixpipe"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	yaml "gopkg.in/yaml.v2"
)

func Parse(in io.Reader) (cfgptr *periwinkle.Cfg, err error) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if !ok {
				panic(r)
			}
			cfgptr = nil
			err = e
		}
	}()

	// these are the defaults
	cfg := periwinkle.Cfg{
		Mailstore:            "./Maildir",
		WebUIDir:             "./www",
		Debug:                true,
		TrustForwarded:       true,
		TwilioAccountID:      os.Getenv("TWILIO_ACCOUNTID"),
		TwilioAuthToken:      os.Getenv("TWILIO_TOKEN"),
		GroupDomain:          "localhost",
		WebRoot:              "locahost:8080",
		DB:                   nil, // the default DB is set later
		DefaultDomainHandler: bounceNoHost,
	}

	datstr, err := ioutil.ReadAll(in)
	if err != nil {
		panic(err)
	}

	var datint interface{}
	err = yaml.Unmarshal(datstr, &datint)
	if err != nil {
		panic(err)
	}

	datmap, ok := datint.(map[interface{}]interface{})
	if !ok {
		panic(err)
	}

	for key, val := range datmap {
		switch key {
		case "Mailstore":
			cfg.Mailstore = maildir.Maildir(getString(key.(string), val))
		case "WebUIDir":
			cfg.WebUIDir = http.Dir(getString(key.(string), val))
		case "Debug":
			cfg.Debug = getBool(key.(string), val)
		case "TrustForwarded":
			cfg.TrustForwarded = getBool(key.(string), val)
		case "TwilioAccountID":
			cfg.TwilioAccountID = getString(key.(string), val)
		case "TwilioAuthToken":
			cfg.TwilioAuthToken = getString(key.(string), val)
		case "GroupDomain":
			cfg.GroupDomain = getString(key.(string), val)
		case "WebRoot":
			cfg.WebRoot = getString(key.(string), val)
		case "DB":
			m, ok := val.(map[interface{}]interface{})
			if !ok {
				panic(fmt.Errorf("value for %q is not a map", key.(string)))
			}
			var driver string
			var source string
			for key, val := range m {
				switch key {
				case "driver":
					driver = getString("DB."+key.(string), val)
				case "source":
					source = getString("DB."+key.(string), val)
				default:
					panic(fmt.Errorf("unknown field: %v", "DB."+key.(string)))
				}
			}
			db, err := gorm.Open(driver, source)
			if err != nil {
				panic(err)
			}
			cfg.DB = &db
		default:
			panic(fmt.Errorf("unknown field: %v", key))
		}
	}

	// Set the default database
	if cfg.DB == nil {
		fmt.Fprintln(os.Stderr, "DB not configured, trying MySQL periwinkle:periwinkle@localhost/periwinkle")
		db, err := gorm.Open("mysql", "periwinkle:periwinkle@/periwinkle?charset=utf8&parseTime=True")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "Failed to connect to MySQL, trying SQLite3 file:periwinkle.sqlite")
			db, err = gorm.Open("sqlite3", "file:periwinkle.sqlite?cache=shared&mode=rwc")
			if err != nil {
				panic(err)
			}
		}
		cfg.DB = &db
	}

	cfg.DB.LogMode(cfg.Debug)

	domain_handlers.GetHandlers(&cfg)

	return &cfg, err
}

func getString(key string, val interface{}) string {
	str, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("value for %q is not a string", key))
	}
	return str
}

func getBool(key string, val interface{}) bool {
	b, ok := val.(bool)
	if !ok {
		panic(fmt.Errorf("value for %q is not a Boolean", key))
	}
	return b
}

func bounceNoHost(io.Reader, string, *gorm.DB, *periwinkle.Cfg) postfixpipe.ExitStatus {
	return postfixpipe.EX_NOHOST
}
