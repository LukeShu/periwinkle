// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
)

type Subscription struct {
	//Id        int
	AddressId int
	GroupId   string
}

func (o Subscription) schema(db *gorm.DB) {
	db.CreateTable(&o).
		AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("address_id", "user_addresses(id)", "RESTRICT", "RESTRICT")
}
