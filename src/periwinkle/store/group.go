// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package store

import (
	"fmt"
	"github.com/jinzhu/gorm"
	he "httpentity"
	"httpentity/util"
	"io"
	"strings"
)

var _ he.Entity = &Group{}
var _ he.NetEntity = &Group{}
var dirGroups he.Entity = newDirGroups()

// Model /////////////////////////////////////////////////////////////

type Group struct {
	Id        string
	Addresses []GroupAddress
}

func (o Group) schema(db *gorm.DB) {
	db.CreateTable(&o)
}

type GroupAddress struct {
	Id      int64
	GroupId string
	Medium  string
	Address string
}

func (o GroupAddress) schema(db *gorm.DB) {
	table := db.CreateTable(&o)
	table.AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	table.AddForeignKey("medium", "media(id)", "RESTRICT", "RESTRICT")
	table.AddUniqueIndex("uniqueness_idx", "medium", "address")
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

func getGroupAddressesByMediumAndGroupId(db *gorm.DB, medium string, groupId string) *GroupAddress {
	var o GroupAddress
	if result := db.Where("medium =? and group_id =?", medium, groupId).Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func getGroupAddressesByMedium(db *gorm.DB, medium string) *GroupAddress {
	var o GroupAddress
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

func (o *Group) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: API: (*Group).Subentity()")
}

func (o *Group) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			// TODO: permission check
			return req.StatusOK(o)
		},
		"PUT": func(req he.Request) he.Response {
			panic("TODO: API: (*Group).Methods()[\"PUT\"]")
		},
		"PATCH": func(req he.Request) he.Response {
			panic("TODO: API: (*Group).Methods()[\"PATCH\"]")
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
			badbody := req.StatusBadRequest(heutil.NetString(fmt.Sprintf("submitted body not what expected")))
			hash, ok := req.Entity.(map[string]interface{}); if !ok { return badbody }
			groupname, ok := hash["groupname"].(string)    ; if !ok { return badbody }

			groupname = strings.ToLower(groupname)

			group := NewGroup(db, groupname)
			if group == nil {
				return req.StatusConflict(heutil.NetString("a group with that name already exists"))
			} else {
				return req.StatusCreated(r, groupname)
			}
		},
	}
	return r
}

func (d t_dirGroups) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirGroups) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*gorm.DB)
	return GetGroupById(db, name)
}
