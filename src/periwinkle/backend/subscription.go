// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Mark Pundmann

package backend

import (
	"locale"
	"periwinkle"
)

type Subscription struct {
	AddressID int64  `sql:"type:bigint       REFERENCES user_addresses(id) ON DELETE CASCADE  ON UPDATE RESTRICT"`
	GroupID   string `sql:"type:varchar(255) REFERENCES groups(id)         ON DELETE CASCADE  ON UPDATE RESTRICT"`
	Confirmed bool
}

func (o Subscription) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func NewSubscription(db *periwinkle.Tx, addressID int64, groupID string, confirmed bool) Subscription {
	subscription := Subscription{
		AddressID: addressID,
		GroupID:   groupID,
		Confirmed: confirmed,
	}
	if err := db.Create(&subscription).Error; err != nil {
		dbError(err)
	}
	return subscription
}

func IsSubscribed(db *periwinkle.Tx, userID string, group Group) int {
	subscriptions := group.GetSubscriptions(db)
	addressIDs := make([]int64, len(subscriptions))
	for i, subscription := range subscriptions {
		addressIDs[i] = subscription.AddressID
	}
	var addresses []UserAddress
	if len(addressIDs) > 0 {
		if result := db.Where("id IN (?)", addressIDs).Find(&addresses); result.Error != nil {
			if !result.RecordNotFound() {
				panic("cant find any subscriptions corresponding user address")
			}
		}
	} else {
		// no subscriptions so user cannot possibly be subscribed
		return 0 // not subscribed
	}
	for _, address := range addresses {
		if address.UserID == userID {
			for _, subscription := range subscriptions {
				if address.ID == subscription.AddressID {
					if subscription.Confirmed {
						return 2
					}
					return 1
				}
			}
		}
	}
	// could not find user in subscribed user addresses, therefore, he/she isn't subscribed
	return 0
}

func IsAdmin(db *periwinkle.Tx, userID string, group Group) bool {
	subscriptions := group.GetSubscriptions(db)
	addressIDs := make([]int64, len(subscriptions))
	for i, subscription := range subscriptions {
		addressIDs[i] = subscription.AddressID
	}
	var addresses []UserAddress
	if len(addressIDs) > 0 {
		if result := db.Where("id IN (?)", addressIDs).Find(&addresses); result.Error != nil {
			if !result.RecordNotFound() {
				panic("cant find any subscriptions corresponding user address")
			}
		}
	} else {
		// no subscriptions so user cannot possibly be subscribed
		return false
	}
	for _, address := range addresses {
		if address.UserID == userID && address.Medium == "admin" {
			return true
		}
	}
	// could not find user in subscribed user addresses, therefore, he/she isn't subscribed
	return false
}

func (sub *Subscription) Delete(db *periwinkle.Tx) {
	if err := db.Where("address_id = ? AND group_id = ?", sub.AddressID, sub.GroupID).Delete(Subscription{}).Error; err != nil {
		dbError(err)
	}
}
