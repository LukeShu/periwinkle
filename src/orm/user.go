// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type User struct{
	id int
	user_id int
	name string
	medium int
	address string
}

// grab User object from database based on id
func getUserById(id int) (*User, error) {
	var user User
	err := con.QueryRow("select * from users where id=?", id).Scan(&user)
	switch {
	case err == sql.ErrNoRows:
		// user does not exist
		return nil, nil
	case err != nil:
		// error talking to the DB
		return nil, err
	default:
		return &user, nil
	}
}

// grab User object by name
func getUserByName(name string)(*User, error){
	var user User
	err := con.QueryRow("select * from users where name=?",name).Scan(&user)
	switch {
		case err == sql.ErrNoRows:
			//user does not exist
			return nil, nil
		case err != nil:
			//error talking to the DB
			return nil, err
		default:
			return &user, nil
	}
}

// grab User by user_id
func getUserByUserId(user_id int)(*User, error){
	var user User
	err := con.QueryRow("select * from users where user_id=?",user_id).Scan(&user)
	switch {
		case err == sql.ErrNoRows:
			//user does not exist
			return nil, nil
		case err != nil:
			//error talking to the DB
			return nil, err
		default:
			return &user, nil
	}
}

func getUserByAddress(address string)(*User, error){
	var user User
	err := con.QueryRow("select * from users where address=?",address).Scan(&user)
	switch {
		case err == sql.ErrNoRows:
			//user does not exist
			return nil, nil
		case err != nil:
			//error talking to the DB
			return nil, err
		default:
			return &user, nil
	}
}
