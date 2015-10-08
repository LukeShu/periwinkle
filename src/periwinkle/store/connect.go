// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker

package store

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

//  Make sure to call "defer db.Close()" after db is returned
func getConnection() *sql.DB {
	db_user := os.Getenv("DBUSERNAME")
	db_pass := os.Getenv("DBPASSWORD")
	db, err := sql.Open("mysql", fmt.Sprint(db_user, ":", db_pass, "@/periwinkle"))

	if err != nil {
		panic("Could not connect to database")
		return nil
	}
	err = db.Ping()
	if err != nil {
		panic("Could not ping database")
		return nil
	}

	return db
}
