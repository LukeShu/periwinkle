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
	// @/test is the current database we are using which is the test database
	db, err := sql.Open("mariadb", fmt.Sprint(db_user, ":", db_pass, "@/test"))

	if err != nil {
		fmt.Printf("Could not connect to database")
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("Could not ping database")
	}

	return db
}

func main() {
	db := getConnection()
	defer db.Close()
}
