// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Mark Pundmann

package store

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"httpentity/util" // heutil
	"io"
	"strings"
)

type Subscription struct {
	Address   UserAddress `json:"addresses"`
	AddressId int64       `json:"-"`
	Group     Group       `json:"group"`
	GroupId   string      `json:"group_id"`
}

func (o Subscription) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("address_id", "user_addresses(id)", "RESTRICT", "RESTRICT").
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

func (o *Subscription) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return he.StatusOK(o)
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			db.Delete(o)
			return he.StatusNoContent()
		},
	}
}

type t_dirSubscriptions struct {
	methods map[string]func(he.Request) he.Response
}

func newDirSubscriptions() t_dirSubscriptions {
	r := t_dirSubscriptions{}
	r.methods = map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			type getfmt struct {
				GroupId string `json:"groupid"`
			}
			var entity getfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return *httperr
			}

			if entity.GroupId == "" {
				return he.StatusUnsupportedMediaType(heutil.NetString(k("groupname can't be emtpy")))
			}
			entity.GroupId = strings.ToLower(entity.GroupId)
			var subscriptions []Subscription
			subscriptions = GetSubscriptionsGroupById(db, entity.GroupId)
			generic := make([]interface{}, len(subscriptions))
			for i, subscription := range subscriptions {
				generic[i] = subscription.Address.Address
			}

			return he.StatusOK(heutil.NetList(generic))
		},
	}
	return r
}

func (o *Subscription) Encoders() map[string]func(io.Writer) error {
	return defaultEncoders(o)
}

func (d t_dirSubscriptions) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirSubscriptions) Subentity(user_id string, group_name string, req he.Request) he.Entity {
	//group_name = strings.ToLower(group_name)
	//db := req.Things["db"].(*gorm.DB)
	panic("TODO")
}
