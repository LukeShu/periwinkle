// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import "database/sql"

type GroupAddr struct {
	id        int
	group_id  int
	medium_id int
	address   string
}

func getGroupAddressById(con DB, id int) (*GroupAddr, error) {
	var g_addr GroupAddr
	err := con.QueryRow("select * from group_addresses where id=?", id).Scan(&g_addr)
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

func getGroupAddressByGroupId(con DB, group_id int) (*GroupAddr, error) {
	var g_addr GroupAddr
	err := con.QueryRow("select * from group_addresses where group_id=?", group_id).Scan(&g_addr)
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

func getGroupAddressByAddress(con DB, address string) (*GroupAddr, error) {
	var g_addr GroupAddr
	err := con.QueryRow("select * from group_addresses where address=?", address).Scan(&g_addr)
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
