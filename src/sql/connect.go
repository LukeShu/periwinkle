package sql

import "os"
import "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

//  Make sure to call "defer db.Close()" after db is returned
func getConnection() *db {
	db_user = os.Getenv("DBUSERNAME")
	db_pass = os.Getenv("DBPASSWORD")
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
	db = getConnection()
	defer db.Close()

}
