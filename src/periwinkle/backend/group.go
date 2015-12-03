// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package backend

import (
	"locale"

	"github.com/jinzhu/gorm"
)

// A Group or mailing list that users may subscribe to.
type Group struct {
	ID            string         `json:"group_id"`
	ReadPublic int `json:"readpublic"`
	ReadConfirmed int `json:"readconfirmed"`
	ExistencePublic int `json:"existencepublic"`
        ExistenceConfirmed int `json:"existenceconfirmed"`
	PostPublic int `json:"postpublic"`
        PostConfirmed int `json:"postconformed"`
        PostMember int `json:"postmember"`
	JoinPublic int `json:"joinpublic"`
        JoinConfirmed int `json:"joinconfirmed"`
        JoinMember int `json:"joinmember"`
	Subscriptions []Subscription `json:"subscriptions"`
}

func (o Group) dbSchema(db *gorm.DB) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func GetGroupByID(db *gorm.DB, id string) *Group {
	var o Group
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	db.Model(&o).Related(&o.Subscriptions)
	return &o
}

func GetGroupsByMember(db *gorm.DB, user User) []Group {
	subscribed := user.GetUserSubscriptions(db)
	var subscriptions map[string]int
	subscriptions = make(map[string]int)
	for _, sub := range subscribed {
		subscriptions[sub.GroupID] = 1
	}
	groupids := make([]string, 0, len(subscriptions))
	for key := range subscriptions {
		groupids = append(groupids, key)
	}
	// use the list of group IDs to get the groups
	var groups []Group
	if len(groupids) > 0 {
		if result := db.Where(groupids).Find(&groups); result.Error != nil {
			if result.RecordNotFound() {
				return nil
			}
			dbError(result.Error)
		}
	}
	// return them
	return groups
}

func GetPublicAndSubscribedGroups(db *gorm.DB, user User) []Group {
	groups := GetGroupsByMember(db, user)
	// also get public groups
	var publicgroups []Group
	if result := db.Where(&Group{ExistencePublic: 1}).Find(&publicgroups); result.Error != nil {
		if !result.RecordNotFound() {
			dbError(result.Error)
		}
	}
	// merge public groups and subscribed groups
	for _, publicgroup := range publicgroups {
		for _, group := range groups {
			if group.ID == publicgroup.ID {
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
		dbError(result.Error)
	}
	return o
}

func NewGroup(db *gorm.DB, name string, existence []int, read []int, post []int, join []int) *Group {
	if name == "" {
		programmerError("Group name can't be empty")
	}
	subscriptions := make([]Subscription, 0)
	o := Group{
		ID:            name,
	        ReadPublic: read[0],
	        ReadConfirmed: read[1],
	        ExistencePublic: existence[0],
	        ExistenceConfirmed: existence[1],
	        PostPublic: post[0],
	        PostConfirmed: post[1],
	        PostMember: post[2],
	        JoinPublic: join[0],
	        JoinConfirmed: join[1],
	        JoinMember: join[2],
		Subscriptions: subscriptions,
	}
	if err := db.Create(&o).Error; err != nil {
		dbError(err)
	}
	return &o
}

// TODO: we should have the database do this.
func checkInput(input int, min int, max int, defaultt int) int {
	if input < min || input > max {
		return defaultt
	}
	return input
}

func (o *Group) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		dbError(err)
	}
}
