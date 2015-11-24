// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package backend

import (
	"github.com/jinzhu/gorm"
)

type Group struct {
	Id            string         `json:"group_id"`
	Existence     int            `json:"existence"` // 1 -> public, 2 -> confirmed, 3 -> member
	Read          int            `json:"read"`      // 1 -> public, 2 -> confirmed, 3 -> member
	Post          int            `json:"post"`      // 1 -> public, 2 -> confirmed, 3 -> moderator
	Join          int            `json:"join"`      // 1 -> auto join, 2 -> confirm to join
	Subscriptions []Subscription `json:"subscriptions"`
}

func (o Group) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

func (o Group) dbSeed(db *gorm.DB) error {
	return db.Create(&Group{
		Id:            "test",
		Existence:     1,
		Read:          1,
		Post:          1,
		Join:          1,
		Subscriptions: []Subscription{},
	}).Error
}

func GetGroupById(db *gorm.DB, id string) *Group {
	var o Group
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	db.Model(&o).Related(&o.Subscriptions)
	return &o
}

func GetGroupsByMember(db *gorm.DB, user User) []Group {
	// turn the User's list of addresses into a list of address IDs
	user_address_ids := make([]int64, len(user.Addresses))
	for i, user_address := range user.Addresses {
		user_address_ids[i] = user_address.Id
	}
	// use the list of address IDs to get a list of subscriptions
	var subscriptions []Subscription
	if len(user_address_ids) > 0 {
		if result := db.Where("address_id IN (?)", user_address_ids).Find(&subscriptions); result.Error != nil {
			if result.RecordNotFound() {
				return nil
			}
			panic(result.Error)
		}
	} else {
		subscriptions = make([]Subscription, 0)
	}
	// turn the list of subscriptions into a list of group IDs
	group_ids := make([]string, len(subscriptions))
	for i, subscription := range subscriptions {
		group_ids[i] = subscription.GroupId
	}
	// use the list of group IDs to get the groups
	var groups []Group
	if result := db.Where(group_ids).Find(&groups); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	// return them
	return groups
}

func GetPublicAndSubscribedGroups(db *gorm.DB, user User) []Group {
	groups := GetGroupsByMember(db, user)
	// also get public groups
	var publicgroups []Group
	if result := db.Where(&Group{Existence: 1}).Find(&publicgroups); result.Error != nil {
		if !result.RecordNotFound() {
			panic(result.Error)
		}
	}
	// merge public groups and subscribed groups
	for _, publicgroup := range publicgroups {
		for _, group := range groups {
			if group.Id == publicgroup.Id {
				break
			}
		}
		groups = append(groups, publicgroup)
	}

	// return them
	return groups
}

func GetAllGroups(db *gorm.DB) []Group {
	var o []Group
	if result := db.Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return o
}

func NewGroup(db *gorm.DB, name string, existence int, read int, post int, join int) *Group {
	if name == "" {
		panic("name can't be empty")
	}
	subscriptions := make([]Subscription, 0)
	o := Group{
		Id:            name,
		Existence:     CheckInput(existence, 1, 3, 1),
		Read:          CheckInput(read, 1, 3, 1),
		Post:          CheckInput(post, 1, 3, 1),
		Join:          CheckInput(existence, 1, 2, 1),
		Subscriptions: subscriptions,
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func CheckInput(input int, min int, max int, defaultt int) int {
	if input < min || input > max {
		return defaultt
	}
	return input
}

func (o *Group) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}
