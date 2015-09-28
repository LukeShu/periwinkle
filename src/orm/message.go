// Copyright 2015 Davis Webb

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Message struct{
	id int
	group_id int
	filename string
	// cached fields??????
}
