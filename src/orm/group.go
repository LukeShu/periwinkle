// Copyright 2015 Davis Webb

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Group struct{
	id int
	name string
}

// get group by its id
func getGroupById(id int)(*Group, error){
	var group Group
	err := con.QueryRow("select * from groups where id=?",id).Scan(&group)
	switch {
		case err == sql.ErrNoRows:
			// group does not exist
			return nil, nil
		case err != nil:
			// error talking to the DB
			return nil, err
		default:
			return &group, nil
	}
}

// get group by the groups name
func getGroupByName(name string)(*Group, error){
	var group Group
	err := con.QueryRow("select * from groups where name=?",name).Scan(&group)
	switch {
		case err == sql.ErrNoRows:
			//user does not exist
			return nil, nil
		case err != nil:
			//error talking to the DB
			return nil, err
		default:
			return &group, nil
	}
}
