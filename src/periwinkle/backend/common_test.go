// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend_test

import (
	"periwinkle"
	. "periwinkle/backend"
	"periwinkle/cfg"
)

func CreateTempDB() *periwinkle.Cfg {
	conf := periwinkle.Cfg{
		Mailstore:      "./Maildir",
		WebUIDir:       "./www",
		Debug:          true,
		TrustForwarded: true,
		GroupDomain:    "localhost",
		WebRoot:        "locahost:8080",
		DB:             nil, // the default DB is set later
	}

	db, err := cfg.OpenDB("sqlite3", "file:temp.sqlite?mode=memory&_txlock=exclusive")
	if err != nil {
		periwinkle.Logf("Error loading sqlite3 database")
	}
	conf.DB = db

	DbSchema(conf.DB)

	return &conf
}
