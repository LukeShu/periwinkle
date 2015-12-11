// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal
// Copyright 2015 Mark Pundmann

package backend

import (
	"locale"
	"periwinkle"
	"strings"
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

func (o Group) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func (o Group) dbSeed(db *periwinkle.Tx) locale.Error {
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

func GetGroupByID(db *periwinkle.Tx, id string) *Group {
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

func GetGroupsByMember(db *periwinkle.Tx, user User) []Group {
	subscriptions := user.GetSubscriptions(db)
	var groupIDs []string
	for _, sub := range subscriptions {
		groupIDs = append(groupIDs, sub.GroupID)
	}
	var groups []Group
	if len(groupIDs) > 0 {
		if err := db.Where("id IN (?)", groupIDs).Find(&groups).Error; err != nil {
			dbError(err)
		}
	} else {
		groups = make([]Group, 0)
	}
	groupsByID := make(map[string]Group)
	for _, group := range groups {
		groupsByID[group.ID] = group
	}

	groupset := map[string]Group{}
	for _, sub := range subscriptions {
		// only add group if user is confirmed member or
		// if group allows non confirmed members to see that it exists
		if sub.Confirmed || groupsByID[sub.GroupID].ExistenceConfirmed == 2 {
			groupset[sub.GroupID] = groupsByID[sub.GroupID]
		}
	}
	var finalgroups []Group
	for _, group := range groupset {
		finalgroups = append(finalgroups, group)
	}

	return finalgroups
}

func GetPublicAndSubscribedGroups(db *periwinkle.Tx, user User) []Group {
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

func GetAllGroups(db *periwinkle.Tx) []Group {
	var o []Group
	if result := db.Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return o
}

func NewGroup(db *periwinkle.Tx, name string, existence []int, read []int, post []int, join []int) *Group {
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

func (grp *Group) GetSubscriptions(db *periwinkle.Tx) []Subscription {
	var subscriptions []Subscription
	if result := db.Where("group_id = ?", grp.ID).Find(&subscriptions); result.Error != nil {
		if result.RecordNotFound() {
			return []Subscription{}
		}
		dbError(result.Error)
	}
	return subscriptions
}

func (grp *Group) GetSubscriberIDs(db *periwinkle.Tx) []string {
	subscriptions := grp.GetSubscriptions(db)
	addressIDs := make([]int64, len(subscriptions))
	for i, sub := range subscriptions {
		addressIDs[i] = sub.AddressID
	}
	var addresses []UserAddress
	if len(addressIDs) > 0 {
		if err := db.Where("id IN (?)", addressIDs).Find(&addresses).Error; err != nil {
			dbError(err)
		}
	} else {
		addresses = []UserAddress{}
	}
	userIDSet := map[string]bool{}
	for _, address := range addresses {
		userIDSet[address.UserID] = true
	}
	userIDs := make([]string, len(userIDSet))
	i := 0
	for userID := range userIDSet {
		userIDs[i] = userID
		i++
	}
	return userIDs
}

func (o *Group) Save(db *periwinkle.Tx) {
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
				var o UserAddress
				db.First(&o, "id = ?", oldsub.AddressID)
				if o.Medium != "noop" && o.Medium != "admin" {
					if err := db.Where("address_id = ? AND group_id = ?", oldsub.AddressID, oldsub.GroupID).Delete(Subscription{}).Error; err != nil {
						dbError(err)
					}
				}
			}
		}

	}

	if err := db.Save(o).Error; err != nil {
		dbError(err)
	}

}

func (grp *Group) Delete(db *periwinkle.Tx) {
	if err := db.Delete(grp).Error; err != nil {
		dbError(err)
	}
}
