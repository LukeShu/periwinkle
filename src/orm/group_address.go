// Copyright 2015 Davis Webb

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type GroupAddr struct{
	id int
	group_id int
	medium_id int
	address string
}

func getGroupAddressById(id int)(*GroupAddr, error){
	var g_addr GroupAddr
	err := con.QueryRow("select * from group_addresses where id=?",id).Scan(&g_addr)
	switch {
		case err == sql.ErrNoRows:
			// group does not exist
			return nil, nil
		case err != nil:
			// error talking to the DB
			return nil, err
		default:
			return &g_addr, nil
	}
}

func getGroupAddressByGroupId(group_id int)(*GroupAddr, error){}

func getGroupAddressByAddress(address string)(*GroupAddr, error){}
