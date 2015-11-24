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

			var new_group Group
			httperr := safeDecodeJSON(req.Entity, &new_group)
			if httperr != nil {
				return *httperr
			}
			if o.Id != new_group.Id {
				return he.StatusConflict(heutil.NetString("Cannot change group id"))
			}
			*o = new_group
			o.backend().Save(db)
			return he.StatusOK(o)
		},
		"PATCH": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)

			patch, ok := req.Entity.(jsonpatch.Patch)
			if !ok {
				return he.StatusUnsupportedMediaType(heutil.NetString("PATCH request must have a patch media type"))
			}
			var new_group Group
			err := patch.Apply(o, &new_group)
			if err != nil {
				return he.StatusConflict(heutil.NetPrintf("%v", err))
			}
			if o.Id != new_group.Id {
				return he.StatusConflict(heutil.NetString("Cannot change user id"))
			}
			*o = new_group
			o.backend().Save(db)
			return he.StatusOK(o)
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			if !backend.IsAdmin(db, sess.UserId, *o.backend()) {
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
				groups = backend.GetPublicAndSubscribedGroups(db, *backend.GetUserById(db, sess.UserId))
			}
			generic := make([]interface{}, len(groups))
			type EnumerateGroup struct {
				Id            string                 `json:"id"`
				Existence     string                 `json:"existence"`
				Read          string                 `json:"read"`
				Post          string                 `json:"post"`
				Join          string                 `json:"join"`
				Subscriptions []backend.Subscription `json:"subscriptions"`
			}

			for i, group := range groups {
				var enum EnumerateGroup
				enum.Id = group.Id
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
				Groupname string            `json:"groupname"`
				Existence backend.Existence `json:"existence"`
				Read      backend.Read      `json:"read"`
				Post      backend.Post      `json:"post"`
				Join      backend.Join      `json:"join"`
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
				int(entity.Existence),
				int(entity.Read),
				int(entity.Post),
				int(entity.Join),
			)
			sess := req.Things["session"].(*backend.Session)
			address := backend.GetAddressByIdAndMedium(db, sess.UserId, "noop")
			if address != nil {
				subscription := backend.Subscription{
					AddressId: address.Id,
					GroupId:   group.Id,
					Confirmed: true,
				}
				db.Create(&subscription)
			}
			if group == nil {
				return he.StatusConflict(heutil.NetString("a group with that name already exists"))
			} else {
				return he.StatusCreated(r, group.Id, req)
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
	group := backend.GetGroupById(db, name)
	if group.Read != 1 && !backend.IsSubscribed(db, sess.UserId, *group) {
		return nil
	}
	return (*Group)(group)
}
