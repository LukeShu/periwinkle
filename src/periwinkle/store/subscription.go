// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
)

type Subscription struct {
	Address   UserAddress
	AddressId int64 `json:"-"`
	Group     Group
	GroupId   string
}

func (o Subscription) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("address_id", "user_addresses(id)", "RESTRICT", "RESTRICT").
		Error
}
