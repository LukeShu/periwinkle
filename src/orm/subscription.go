// Copyright 2015 Davis Webb

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Subscription struct{
	id int
	address_id int
	group_id int
}
