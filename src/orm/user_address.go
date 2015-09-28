// Copyright 2015 Davis Webb

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type UserAddress struct{
	id int
	user_id int
	medium_id int
	address string
}
