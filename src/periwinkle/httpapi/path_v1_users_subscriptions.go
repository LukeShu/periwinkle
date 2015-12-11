// Copyright 2015 Luke Shumaker

package httpapi

import (
	"encoding/json"
	he "httpentity"
	"httpentity/rfc7231"
	"periwinkle"
	"periwinkle/backend"
	"strings"
)

var _ he.NetEntity = &userSubscriptions{}
var _ he.EntityGroup = &userSubscriptions{}

type userSubscriptions struct {
	user
	groupID string
	values  []backend.Subscription
}

func (usr *userSubscriptions) Subentity(name string, req he.Request) he.Entity {
	parts := strings.Split(name, ":")
	if len(parts) != 3 {
		return nil
	}
	groupID := strings.ToLower(parts[0])
	medium := parts[1]
	address := parts[2]

	db := req.Things["db"].(*periwinkle.Tx)

	for _, addr := range usr.Addresses {
		if addr.Medium == medium && addr.Address == address {
			for _, sub := range addr.GetSubscriptions(db) {
				if sub.GroupID == groupID {
					ret := userSubscription(sub)
					return &ret
				}
			}
		}
	}
	return nil
}

func (usr *userSubscriptions) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			usr.groupID = req.URL.Query().Get("group_id")
			usr.values = usr.backend().GetFrontEndSubscriptions(db)
			return rfc7231.StatusOK(usr)
		},
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			sess := req.Things["session"].(*backend.Session)
			type postfmt struct {
				GroupID string `json:"group_id"`
				Medium  string `json:"medium,omitempty"`
				Address string `json:"address,omitempty"`
			}
			var entity postfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return *httperr
			}
			entity.GroupID = strings.ToLower(entity.GroupID)

			var address *backend.UserAddress
			if entity.Medium == "" && entity.Address == "" {
				address = &usr.Addresses[0]
				for _, addr := range usr.Addresses {
					if addr.SortOrder < address.SortOrder {
						address = &addr
					}
				}
			} else {
				for _, addr := range usr.Addresses {
					if addr.Medium == entity.Medium && addr.Address == entity.Address {
						address = &addr
						break
					}
				}
			}
			if address == nil {
				return rfc7231.StatusConflict(he.NetPrintf("You don't have that address"))
			}
			backend.NewSubscription(db, address.ID, entity.GroupID, sess != nil && sess.UserID == usr.ID)
			return rfc7231.StatusCreated(usr, entity.GroupID+":"+entity.Medium+":"+entity.Address, req)
		},
	}
}

func (usr *userSubscriptions) Encoders() map[string]he.Encoder {
	return defaultEncoders(usr)
}

func (usr *userSubscriptions) MarshalJSON() ([]byte, error) {
	addressByID := map[int64]backend.UserAddress{}
	for _, addr := range usr.Addresses {
		addressByID[addr.ID] = addr
	}
	type subscriptionfmt struct {
		Medium  string `json:"medium"`
		Address string `json:"address"`
	}
	ret := map[string][]subscriptionfmt{}
	for _, subscription := range usr.values {
		out := subscriptionfmt{
			Medium:  addressByID[subscription.AddressID].Medium,
			Address: addressByID[subscription.AddressID].Address,
		}
		if list, ok := ret[subscription.GroupID]; ok {
			ret[subscription.GroupID] = append(list, out)
		} else {
			ret[subscription.GroupID] = []subscriptionfmt{out}
		}
	}
	if usr.groupID == "" {
		return json.Marshal(ret)
	} else {
		return json.Marshal(ret[usr.groupID])
	}
}

////////////////////////////////////////////////////////////////////////////////////

type userSubscription backend.Subscription

var _ he.Entity = &userSubscription{}

func (o *userSubscription) backend() *backend.Subscription { return (*backend.Subscription)(o) }

func (sub *userSubscription) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			sub.backend().Delete(db)
			return rfc7231.StatusNoContent()
		},
	}
}
