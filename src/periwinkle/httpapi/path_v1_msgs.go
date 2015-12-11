// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"periwinkle"
	"periwinkle/backend"
)

var _ he.Entity = &message{}
var _ he.NetEntity = &message{}
var _ he.EntityGroup = dirMessages{}

type message backend.Message

func (o *message) backend() *backend.Message { return (*backend.Message)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *message) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: SMTP: (*message).Subentity()")
}

func (o *message) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return rfc7231.StatusOK(o)
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *message) Encoders() map[string]he.Encoder {
	panic("TODO: API: (*message).Encoders()")
}

func (o *message) SubentityNotFound(name string, req he.Request) he.Response {
	panic("TODO: SMTP: (*message).SubentityNotFound()")
}

// Directory ("Controller") //////////////////////////////////////////

type dirMessages struct {
	methods map[string]func(he.Request) he.Response
}

func newDirMessages() dirMessages {
	r := dirMessages{}
	r.methods = map[string]func(he.Request) he.Response{}
	return r
}

func (d dirMessages) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d dirMessages) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(*periwinkle.Tx)
	message := (*message)(backend.GetMessageByID(db, name))
	sess := req.Things["session"].(*backend.Session)
	grp := backend.GetGroupByID(db, message.GroupID)
	if grp.ReadPublic == 1 {
		subscribed := backend.IsSubscribed(db, sess.UserID, *grp)
		if (grp.ReadConfirmed == 1 && subscribed == 1) || subscribed == 0 {
			return nil
		}
	}
	return message
}

func (d dirMessages) SubentityNotFound(name string, req he.Request) he.Response {
	return rfc7231.StatusNotFound(nil)
}
