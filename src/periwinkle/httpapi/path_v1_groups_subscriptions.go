// Copyright 2015 Luke Shumaker

package httpapi

import (
	"encoding/json"
	he "httpentity"
	"httpentity/rfc7231"
	"periwinkle"
)

var _ he.NetEntity = &groupSubscriptions{}
var _ he.Entity = &groupSubscriptions{}

type groupSubscriptions struct {
	group
	values []string
}

func (grp *groupSubscriptions) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			grp.values = grp.backend().GetSubscriberIDs(db)
			return rfc7231.StatusOK(grp)
		},
	}
}

func (grp *groupSubscriptions) Encoders() map[string]he.Encoder {
	return defaultEncoders(grp)
}

func (grp *groupSubscriptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(grp.values)
}
