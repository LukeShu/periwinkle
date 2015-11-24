// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Mark Pundmann

package backend

import (
	"github.com/jinzhu/gorm"
)

type Subscription struct {
	Address   UserAddress `json:"addresses"`
	AddressId int64       `json:"-"`
	Group     Group       `json:"group"`
	GroupId   string      `json:"group_id"`
	Confirmed bool        `json:"confirmed"`
}

func (o Subscription) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT").
		AddForeignKey("address_id", "user_addresses(id)", "CASCADE", "RESTRICT").
		Error
}

func GetSubscriptionsGroupById(db *gorm.DB, groupId string) []Subscription {
	var o []Subscription
	if result := db.Where("group_id = ?", groupId).Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return o
}

func IsSubscribed(db *gorm.DB, userid string, group Group) bool {
	subscriptions := GetSubscriptionsGroupById(db, group.Id)
	address_ids := make([]int64, len(subscriptions))
	for i, subscription := range subscriptions {
		address_ids[i] = subscription.AddressId
	}
	var addresses []UserAddress
	if len(address_ids) > 0 {
		if result := db.Where("id IN (?)", address_ids).Find(&addresses); result.Error != nil {
			if !result.RecordNotFound() {
				panic("cant find any subscriptions corresponding user address")
			}
		}
	} else {
		// no subscriptions so user cannot possibly be subscribed
		return false
	}
	for _, address := range addresses {
		if address.UserId == userid {
			return true
		}
	}
	// could not find user in subscribed user addresses, therefore, he/she isn't subscribed
	return false
}

func IsAdmin(db *gorm.DB, userid string, group Group) bool {
	subscriptions := GetSubscriptionsGroupById(db, group.Id)
	address_ids := make([]int64, len(subscriptions))
	for i, subscription := range subscriptions {
		address_ids[i] = subscription.AddressId
	}
	var addresses []UserAddress
	if len(address_ids) > 0 {
		if result := db.Where("id IN (?)", address_ids).Find(&addresses); result.Error != nil {
			if !result.RecordNotFound() {
				panic("cant find any subscriptions corresponding user address")
			}
		}
	} else {
		// no subscriptions so user cannot possibly be subscribed
		return false
	}
	for _, address := range addresses {
		if address.UserId == userid && address.Medium == "admin" {
			return true
		}
	}
	// could not find user in subscribed user addresses, therefore, he/she isn't subscribed
	return false
}
