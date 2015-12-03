// Copyright 2015 Davis Webb
package backend

import (
	"periwinkle"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

func CreateTempCfg() *periwinkle.Cfg {
	cfg := periwinkle.Cfg{
		Mailstore:      "./Maildir",
		WebUIDir:       "./www",
		Debug:          true,
		TrustForwarded: true,
		GroupDomain:    "localhost",
		WebRoot:        "locahost:8080",
		DB:             nil, // the default DB is set later
	}

	db, err := gorm.Open("sqlite3", "file:temp.sqlite?mode=memory&_txlock=exclusive")

	if err != nil {
		periwinkle.Logf("Error loading sqlite3 database")
	}

	cfg.DB = &db

	DbSchema(cfg.DB)

	return &cfg
}

func TestNewUser(t *testing.T) {
	t.Log("Testing NewUser")

	var cfg *periwinkle.Cfg

	cfg = CreateTempCfg()

	user := NewUser(cfg.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	switch {
	case strings.Compare(user.ID, "JohnDoe") != 0:
		t.Error("User ID was not properly set to the user object")
	case strings.Compare(user.FullName, "") != 0:
		t.Error("User name was not properly set to the user object")
	}
}
