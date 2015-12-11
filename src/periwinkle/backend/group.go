// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package backend

import (
	"locale"
	//"periwinkle/cfg"
	"periwinkle"
	"periwinkle/twilio"
	"strings"

	"github.com/jinzhu/gorm"
)

// A Group or mailing list that users may subscribe to.
type Group struct {
	ID                 string         `json:"group_id"`
	ReadPublic         int            `json:"readpublic"`
	ReadConfirmed      int            `json:"readconfirmed"`
	ExistencePublic    int            `json:"existencepublic"`
	ExistenceConfirmed int            `json:"existenceconfirmed"`
	PostPublic         int            `json:"postpublic"`
	PostConfirmed      int            `json:"postconformed"`
	PostMember         int            `json:"postmember"`
	JoinPublic         int            `json:"joinpublic"`
	JoinConfirmed      int            `json:"joinconfirmed"`
	JoinMember         int            `json:"joinmember"`
	Subscriptions      []Subscription `json:"subscriptions"`
}

func (o Group) dbSchema(db *gorm.DB) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func (o Group) dbSeed(db *gorm.DB) locale.Error {
	existence := [2]int{2, 2}
	read := [2]int{2, 2}
	post := [3]int{1, 1, 1}
	join := [3]int{1, 1, 1}
	return locale.UntranslatedError(db.Create(&Group{
		ID:                 "test",
		ReadPublic:         read[0],
		ReadConfirmed:      read[1],
		ExistencePublic:    existence[0],
		ExistenceConfirmed: existence[1],
		PostPublic:         post[0],
		PostConfirmed:      post[1],
		PostMember:         post[2],
		JoinPublic:         join[0],
		JoinConfirmed:      join[1],
		JoinMember:         join[2],
		Subscriptions:      []Subscription{},
	}).Error)
}

func GetGroupByID(db *gorm.DB, id string) *Group {
	id = strings.ToLower(id)
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
	var groups []Group
	for _, sub := range subscribed {
		// only add group if user is confirmed member or
		// if group allows non confirmed members to see that it exists
		if sub.Confirmed || sub.Group.ExistenceConfirmed == 2 {
			groups = append(groups, sub.Group)
		}
	}

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
	name = strings.ToLower(name)
	subscriptions := make([]Subscription, 0)
	o := Group{
		ID:                 name,
		ReadPublic:         read[0],
		ReadConfirmed:      read[1],
		ExistencePublic:    existence[0],
		ExistenceConfirmed: existence[1],
		PostPublic:         post[0],
		PostConfirmed:      post[1],
		PostMember:         post[2],
		JoinPublic:         join[0],
		JoinConfirmed:      join[1],
		JoinMember:         join[2],
		Subscriptions:      subscriptions,
	}
	if err := db.Create(&o).Error; err != nil {
		dbError(err)
	}
	return &o
}

func (o *Group) Save(db *gorm.DB, cfg *periwinkle.Cfg, moderate bool, userid string) {
	if o.Subscriptions != nil {
		var oldSubscriptions []Subscription
		db.Model(o).Related(&oldSubscriptions)

		for _, oldsub := range oldSubscriptions {
			match := false

			for _, newsub := range o.Subscriptions {
				if newsub.AddressID == oldsub.AddressID {
					match = true
					break
				}
			}
			if !match {
				if err := db.Where("addressid = ? AND groupid = ?", oldsub.AddressID, oldsub.GroupID).Delete(Subscription{}).Error; err != nil {
					dbError(err)
				}
			}
		}

		if cfg != nil {
			for _, newsub := range o.Subscriptions {
				match := false

				for _, oldsub := range oldSubscriptions {
					if newsub.AddressID == oldsub.AddressID {
						match = true
						break
					}
				}
				if !match {
					twilio_num := GetUnusedTwilioNumbersByUser(cfg, db, userid)
					if twilio_num == nil {
						new_num, err := twilio.NewPhoneNum(cfg)
						if err != nil {
							return
						}
						AssignTwilioNumber(db, userid, o.ID, new_num)
					}
					AssignTwilioNumber(db, userid, o.ID, twilio_num[0])

				}

			}
		}

	}

	if err := db.Save(o).Error; err != nil {
		dbError(err)
	}

}
