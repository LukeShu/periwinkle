// Copyright 2015 Davis Webb

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Medium struct{
	id int
}

func getMediumById(id int)(*Medium, error){
	var med Medium
	err := con.QueryRow("select * from group_addresses where id=?",id).Scan(&med)
	switch {
		case err == sql.ErrNoRows:
			// group does not exist
			return nil, nil
		case err != nil:
			// error talking to the DB
			return nil, err
		default:
			return &med, nil
	}
}
