// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package orm

import "database/sql"

type Medium struct {
	Id string
}

func GetMedium(con DB, id string) (*Medium, error) {
	var med Medium
	err := con.QueryRow("SELECT * FROM group_addresses WHERE id=?", id).Scan(&med)
	switch {
	case err == sql.ErrNoRows:
		// group does not exist
		return nil, nil
	case err != nil:
		// error talking to the DB
		return nil, err
	default:
		// all ok
		return &med, nil
	}
}
