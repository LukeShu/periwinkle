// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package store

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"httpentity/util" // heutil
	"io"
	"jsonpatch"
	"periwinkle/util" // putil
	"strings"
)

var _ he.Entity = &Group{}
var _ he.NetEntity = &Group{}
var dirGroups he.Entity = newDirGroups()

// Model /////////////////////////////////////////////////////////////

type Group struct {
	Id        string         `json:"group_id"`
	Addresses []GroupAddress `json:"addresses"`
}

func (o Group) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

func (o Group) dbSeed(db *gorm.DB) error {
	errs := []error{}
	errHelper(&errs, db.Create(&Group{"test", []GroupAddress{{0,"test","twilio","add_twilio_phone_number", "test_user"}}}).Error)
	return errorList(errs)
}

type GroupAddress struct {
	Id      int64  `json:"group_address_id"`
	GroupId string `json:"group_id"`
	Medium  string `json:"medium"`
	Address string `json:"address"`
	UserId 	string `json:"user_id"`
}

func (o GroupAddress) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("medium", "media(id)", "RESTRICT", "RESTRICT").
		AddUniqueIndex("uniqueness_idx", "medium", "address").
		Error
}

func GetGroupById(db *gorm.DB, id string) *Group {
	var o Group
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func GetUsersInGroup(db *gorm.DB, groupId string) *[]User {
	var users []User
	err := db.Joins("inner join user_addresses on user_addresses.user_id = users.id").Joins(
		"inner join subscriptions on subscriptions.address_id = user_addresses.id").Where(
		"subscriptions.group_id = ?", groupId).Find(&users)
	if err != nil {
		panic("could not get users in group")
	}
	return &users
}

func GetGroupByAddress(db *gorm.DB, address string) *Group {
	var o Group
	if result := db.Joins("inner join groups on group_addresses.group_id = groups.id").Where("group_addresses.address = ?", address).Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func GetGroupAddressesByMediumAndGroupId(db *gorm.DB, medium string, groupId string) *[]GroupAddress {
	var o []GroupAddress
	if result := db.Where("medium =? and group_id =?", medium, groupId).Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func GetGroupAddressesByMedium(db *gorm.DB, medium string) *[]GroupAddress {
	var o []GroupAddress
	if result := db.Where("medium =?", medium).Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func NewGroup(db *gorm.DB, name string) *Group {
	o := Group{Id: name}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func NewGroupAddress(db *gorm.DB, id int64, group_id string, medium string, address string, user_id string) *GroupAddress {
	o := GroupAddress{
		Id      	:id,
		GroupId 	:group_id,
		Medium  	:medium,
		Address 	:address,
		UserId 		:user_id,
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func (o *Group) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}

func (o *Group) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: API: (*Group).Subentity()")
}

func (o *Group) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return he.StatusOK(o)
		},
		"PUT": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*Session)
			users := GetUsersInGroup(db, o.Id)
			flag := false
			for _, user := range *users {
				if sess.UserId == user.Id {
					flag = true
					break
				}
			}
			if !flag {
				return he.StatusForbidden(heutil.NetString("Unauthorized user"))
			}

			var new_group Group
			err := safeDecodeJSON(req.Entity, &new_group)
			if err != nil {
				return err.Response()
			}
			if o.Id != new_group.Id {
				return he.StatusConflict(heutil.NetString("Cannot change group id"))
			}
			*o = new_group
			o.Save(db)
			return he.StatusOK(o)
		},
		"PATCH": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*Session)
			users := GetUsersInGroup(db, o.Id)
			flag := false
			for _, user := range *users {
				if sess.UserId == user.Id {
					flag = true
					break
				}
			}
			if !flag {
				return he.StatusForbidden(heutil.NetString("Unauthorized user"))
			}

			patch, ok := req.Entity.(jsonpatch.Patch)
			if !ok {
				return putil.HTTPErrorf(415, "PATCH request must have a patch media type").Response()
			}
			var new_group Group
			err := patch.Apply(o, &new_group)
			if err != nil {
				return putil.HTTPErrorf(409, "%v", err).Response()
			}
			if o.Id != new_group.Id {
				return he.StatusConflict(heutil.NetString("Cannot change user id"))
			}
			*o = new_group
			o.Save(db)
			return he.StatusOK(o)
		},
		"DELETE": func(req he.Request) he.Response {
			panic("TODO: API: (*Group).Methods()[\"DELETE\"]")
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *Group) Encoders() map[string]func(io.Writer) error {
	return defaultEncoders(o)
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirGroups struct {
	methods map[string]func(he.Request) he.Response
}

func newDirGroups() t_dirGroups {
	r := t_dirGroups{}
	r.methods = map[string]func(he.Request) he.Response{
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			type postfmt struct {
				Groupname string `json:"groupname"`
			}
			var entity postfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return httperr.Response()
			}

			entity.Groupname = strings.ToLower(entity.Groupname)

			group := NewGroup(db, entity.Groupname)
			if group == nil {
				return he.StatusConflict(heutil.NetString("a group with that name already exists"))
			} else {
				return he.StatusCreated(r, entity.Groupname, req)
			}
		},
	}
	return r
}

func (d t_dirGroups) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirGroups) Subentity(name string, req he.Request) he.Entity {
        name = strings.ToLower(name)
        sess := req.Things["session"].(*Session)
        if sess == nil && req.Method == "POST" {
                group, ok := req.Things["group"].(Group)
                if !ok {
                        return nil
                }
                if group.Id == name {
                        return &group
                }
                return nil
        }
        db := req.Things["db"].(*gorm.DB)
	return GetGroupById(db, name)
}
