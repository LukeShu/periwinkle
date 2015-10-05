// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import "database/sql"

type Group struct {
	id   int
	name string
}

func getGroupById(con DB, id int) (*Group, error) {
	var group Group
	err := con.QueryRow("select * from groups where id=?", id).Scan(&group)
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

func GetGroupByName(con DB, name string) (*Group, error) {
	var group Group
	err := con.QueryRow("select * from groups where name=?", name).Scan(&group)
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
