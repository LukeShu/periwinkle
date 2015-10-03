import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "golang.org/x/crypto/bcrypt"

cost Int = 8;

/* 
If you want to run a simple server on localhost:8000 and 
you have python:

python -m SimpleHTTPServer 

*/

func setPassword(password string)(string, error){
	// ask luke about byte
	/*Err error 
	if (Err == nil)*/
	hash, err := GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}

func checkPassword(password string){
	// digest given password
	// pull users digested password from database
	// compare
	
	err := con.QueryRow("select password from users where password=?",
			password)
	check := CompareHashAndPassword(err, []byte(password))
	return check
}

func createUser(UserName string, hash string){
	// send user and digested password to DB

	
}




