// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package httpapi

import (
	he "httpentity"
	"httpentity/heutil"
	"httpentity/rfc7231"
	"io"
	"jsonpatch"
	"periwinkle/backend"
	"strings"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &group{}
var _ he.NetEntity = &group{}
var _ he.Entity = &dirGroups{}

type group backend.Group

func (o *group) backend() *backend.Group { return (*backend.Group)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *group) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: API: (*group).Subentity()")
}

func (o *group) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return rfc7231.StatusOK(o)
		},
		"PUT": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)

			var newGroup group
			httperr := safeDecodeJSON(req.Entity, &newGroup)
			if httperr != nil {
				return *httperr
			}
			if o.ID != newGroup.ID {
				return rfc7231.StatusConflict(heutil.NetString("Cannot change group id"))
			}
			*o = newGroup
			o.backend().Save(db)
			return rfc7231.StatusOK(o)
		},
		"PATCH": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)

			patch, ok := req.Entity.(jsonpatch.Patch)
			if !ok {
				return rfc7231.StatusUnsupportedMediaType(heutil.NetString("PATCH request must have a patch media type"))
			}
			var newGroup group
			err := patch.Apply(o, &newGroup)
			if err != nil {
				return rfc7231.StatusConflict(heutil.NetPrintf("%v", err))
			}
			if o.ID != newGroup.ID {
				return rfc7231.StatusConflict(heutil.NetString("Cannot change user id"))
			}
			*o = newGroup
			o.backend().Save(db)
			return rfc7231.StatusOK(o)
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			if !backend.IsAdmin(db, sess.UserID, *o.backend()) {
				return rfc7231.StatusForbidden(heutil.NetString("Unauthorized user"))
			}
			db.Delete(o)
			return rfc7231.StatusNoContent()
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *group) Encoders() map[string]func(io.Writer) error {
	return defaultEncoders(o)
}

// Directory ("Controller") //////////////////////////////////////////

type dirGroups struct {
	methods map[string]func(he.Request) he.Response
}

func newDirGroups() dirGroups {
	r := dirGroups{}
	r.methods = map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			var groups []backend.Group
			type getfmt struct {
				visibility string
			}
			var entity getfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				entity.visibility = "subscribed"
			}
			if sess == nil {
				groups = []backend.Group{}
			} else if entity.visibility == "subscribed" {
				groups = backend.GetGroupsByMember(db, *backend.GetUserByID(db, sess.UserID))
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

			for i, grp := range groups {
				var enum EnumerateGroup
				enum.ID = grp.ID
				enum.Existence = backend.Existence(grp.Existence).String()
				enum.Read = backend.Read(grp.Read).String()
				enum.Post = backend.Post(grp.Post).String()
				enum.Join = backend.Join(grp.Join).String()
				enum.Subscriptions = grp.Subscriptions
				generic[i] = enum
			}
			return rfc7231.StatusOK(heutil.NetList(generic))
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
				return rfc7231.StatusUnsupportedMediaType(heutil.NetString("groupname can't be emtpy"))
			}

			entity.Groupname = strings.ToLower(entity.Groupname)
			grp := backend.NewGroup(
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
					GroupID:   grp.ID,
					Confirmed: true,
				}
				db.Create(&subscription)
			}
			if grp == nil {
				return rfc7231.StatusConflict(heutil.NetString("a group with that name already exists"))
			} else {
				return rfc7231.StatusCreated(r, grp.ID, req)
			}
		},
	}
	return r
}

func (d dirGroups) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d dirGroups) Subentity(name string, req he.Request) he.Entity {
	name = strings.ToLower(name)
	db := req.Things["db"].(*gorm.DB)
	// TODO: permissions check
	sess := req.Things["session"].(*backend.Session)
	grp := backend.GetGroupByID(db, name)
	if grp.Read != 1 && !backend.IsSubscribed(db, sess.UserID, *grp) {
		return nil
	}
	return (*group)(grp)
}
