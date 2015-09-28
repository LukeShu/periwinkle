// Copyright 2015 Davis Webb

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Session struct{
	id int
	user_id int
	last_used string // not 100% on this
}
