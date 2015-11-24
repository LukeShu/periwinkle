// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package httpapi

import (
	he "httpentity"
	"httpentity/heutil"
	"io"
	"jsonpatch"
	"periwinkle/backend"
	"strings"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &Group{}
var _ he.NetEntity = &Group{}
var dirGroups he.Entity = newDirGroups()

type Group backend.Group

func (o *Group) backend() *backend.Group { return (*backend.Group)(o) }

// Model /////////////////////////////////////////////////////////////

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

			var newGroup Group
			httperr := safeDecodeJSON(req.Entity, &newGroup)
			if httperr != nil {
				return *httperr
			}
			if o.ID != newGroup.ID {
				return he.StatusConflict(heutil.NetString("Cannot change group id"))
			}
			*o = newGroup
			o.backend().Save(db)
			return he.StatusOK(o)
		},
		"PATCH": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)

			patch, ok := req.Entity.(jsonpatch.Patch)
			if !ok {
				return he.StatusUnsupportedMediaType(heutil.NetString("PATCH request must have a patch media type"))
			}
			var newGroup Group
			err := patch.Apply(o, &newGroup)
			if err != nil {
				return he.StatusConflict(heutil.NetPrintf("%v", err))
			}
			if o.ID != newGroup.ID {
				return he.StatusConflict(heutil.NetString("Cannot change user id"))
			}
			*o = newGroup
			o.backend().Save(db)
			return he.StatusOK(o)
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			if !backend.IsAdmin(db, sess.UserID, *o.backend()) {
				return he.StatusForbidden(heutil.NetString("Unauthorized user"))
			}
			db.Delete(o)
			return he.StatusNoContent()
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
		"GET": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			var groups []backend.Group
			if sess == nil {
				groups = []backend.Group{}
			} else {
				groups = backend.GetPublicAndSubscribedGroups(db, *backend.GetUserByID(db, sess.UserID))
			}
			generic := make([]interface{}, len(groups))
			type EnumerateGroup struct {
				ID            string                 `json:"id"`
				Existence     string                 `json:"existence"`
				Read          string                 `json:"read"`
				Post          string                 `json:"post"`
				Join          string                 `json:"join"`
				Subscriptions []backend.Subscription `json:"subscriptions"`
			}

			for i, group := range groups {
				var enum EnumerateGroup
				enum.ID = group.ID
				enum.Existence = backend.Existence(group.Existence).String()
				enum.Read = backend.Read(group.Read).String()
				enum.Post = backend.Post(group.Post).String()
				enum.Join = backend.Join(group.Join).String()
				enum.Subscriptions = group.Subscriptions
				generic[i] = enum
			}
			return he.StatusOK(heutil.NetList(generic))
		},
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			type postfmt struct {
				Groupname string `json:"groupname"`
				Existence string `json:"existence"`
				Read      string `json:"read"`
				Post      string `json:"post"`
				Join      string `json:"join"`
			}
			var entity postfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return *httperr
			}

			if entity.Groupname == "" {
				return he.StatusUnsupportedMediaType(heutil.NetString("groupname can't be emtpy"))
			}

			entity.Groupname = strings.ToLower(entity.Groupname)
			group := backend.NewGroup(
				db,
				entity.Groupname,
				backend.Reverse(entity.Existence),
				backend.Reverse(entity.Read),
				backend.Reverse(entity.Post),
				backend.Reverse(entity.Join),
			)
			sess := req.Things["session"].(*backend.Session)
			address := backend.GetAddressByIDAndMedium(db, sess.UserID, "noop")
			if address != nil {
				subscription := backend.Subscription{
					AddressID: address.ID,
					GroupID:   group.ID,
					Confirmed: true,
				}
				db.Create(&subscription)
			}
			if group == nil {
				return he.StatusConflict(heutil.NetString("a group with that name already exists"))
			} else {
				return he.StatusCreated(r, group.ID, req)
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
	db := req.Things["db"].(*gorm.DB)
	// TODO: permissions check
	sess := req.Things["session"].(*backend.Session)
	group := backend.GetGroupByID(db, name)
	if group.Read != 1 && !backend.IsSubscribed(db, sess.UserID, *group) {
		return nil
	}
	return (*Group)(group)
}
