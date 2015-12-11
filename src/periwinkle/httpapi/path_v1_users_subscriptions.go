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
	values []backend.Subscription
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
			usr.values = usr.backend().GetSubscriptions(db)
			return rfc7231.StatusOK(usr)
		},
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			sess := req.Things["session"].(*backend.Session)
			type postfmt struct {
				GroupID string `json:"group_id"`
				Medium  string `json:"medium"`
				Address string `json:"address"`
			}
			var entity postfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return *httperr
			}
			if entity.GroupID == "" || entity.Medium == "" || entity.Address == "" {
				return rfc7231.StatusUnsupportedMediaType(he.NetPrintf("group_id, medium, and address can't be emtpy"))
			}
			entity.GroupID = strings.ToLower(entity.GroupID)

			for _, addr := range usr.Addresses {
				if addr.Medium == entity.Medium && addr.Address == entity.Address {
					backend.NewSubscription(db, addr.ID, entity.GroupID, sess != nil && sess.UserID == usr.ID)
					return rfc7231.StatusCreated(usr, entity.GroupID+":"+entity.Medium+":"+entity.Address, req)
				}
			}
			return rfc7231.StatusConflict(he.NetPrintf("You don't have that address"))
		},
	}
}

func (usr *userSubscriptions) Encoders() map[string]he.Encoder {
	return defaultEncoders(usr)
}

func (usr *userSubscriptions) MarshalJSON() ([]byte, error) {
	ret := map[string][]backend.Subscription{}
	for _, subscription := range usr.values {
		if list, ok := ret[subscription.GroupID]; ok {
			ret[subscription.GroupID] = append(list, subscription)
		} else {
			ret[subscription.GroupID] = []backend.Subscription{subscription}
		}
	}
	return json.Marshal(ret)
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
