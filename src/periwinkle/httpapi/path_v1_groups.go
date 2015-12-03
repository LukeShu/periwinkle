// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"jsonpatch"
	"periwinkle/backend"
	"strings"

	"github.com/jinzhu/gorm"
)

var _ he.EntityGroup = &group{}
var _ he.NetEntity = &group{}
var _ he.EntityGroup = &dirGroups{}

type group backend.Group

func (o *group) backend() *backend.Group { return (*backend.Group)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *group) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: API: (*group).Subentity()")
}

func (o *group) SubentityNotFound(name string, req he.Request) he.Response {
	panic("TODO: API: (*group).SubentityNotFound()")
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
				return rfc7231.StatusConflict(he.NetPrintf("Cannot change group id"))
			}
			*o = newGroup
			o.backend().Save(db)
			return rfc7231.StatusOK(o)
		},
		"PATCH": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)

			patch, ok := req.Entity.(jsonpatch.Patch)
			if !ok {
				return rfc7231.StatusUnsupportedMediaType(he.NetPrintf("PATCH request must have a patch media type"))
			}
			var newGroup group
			err := patch.Apply(o, &newGroup)
			if err != nil {
				return rfc7231.StatusConflict(he.NetPrintf("%v", err))
			}
			if o.ID != newGroup.ID {
				return rfc7231.StatusConflict(he.NetPrintf("Cannot change user id"))
			}
			*o = newGroup
			o.backend().Save(db)
			return rfc7231.StatusOK(o)
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			if !backend.IsAdmin(db, sess.UserID, *o.backend()) {
				return rfc7231.StatusForbidden(he.NetPrintf("Unauthorized user"))
			}
			db.Delete(o)
			return rfc7231.StatusNoContent()
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *group) Encoders() map[string]he.Encoder {
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
				//groups = GetAllGroups(db)
				groups = backend.GetPublicAndSubscribedGroups(db, *backend.GetUserByID(db, sess.UserID))
			}
			type EnumerateGroup struct {
                                Groupname string `json:"groupname"`
                                Post map[string]string `json:"post"`
                                Join map[string]string `json:"join"`
                                Read map[string]string `json:"read"`
                                Existence map[string]string `json:"existence"`
				Subscriptions []backend.Subscription `json:"subscriptions"`
			}
			data := make([]EnumerateGroup, len(groups))

			for i, grp := range groups {
				var enum EnumerateGroup
				enum.Groupname = grp.ID
				exist := [...]int{grp.ExistencePublic, grp.ExistenceConfirmed}
				enum.Existence = backend.ReadExist(exist)
                                read := [...]int{grp.ReadPublic, grp.ReadConfirmed}
				enum.Read = backend.ReadExist(read)
                                post := [...]int{grp.PostPublic, grp.PostConfirmed, grp.PostMember}
				enum.Post =  backend.PostJoin(post)
                                join := [...]int{grp.JoinPublic, grp.JoinConfirmed, grp.JoinMember}
				enum.Join = backend.PostJoin(join)
				enum.Subscriptions = grp.Subscriptions
				data[i] = enum
			}
			return rfc7231.StatusOK(he.NetJSON{Data: data})
		},
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			type Response1 struct {
				Groupname string `json:"groupname"`
				Post map[string]string `json:"post"`
				Join map[string]string `json:"join"`
				Read map[string]string `json:"read"`
				Existence map[string]string `json:"existence"`
			}
			var entity Response1
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return *httperr
			}


			if entity.Groupname == "" {
				return rfc7231.StatusUnsupportedMediaType(he.NetPrintf("groupname can't be emtpy"))
			}

			entity.Groupname = strings.ToLower(entity.Groupname)
			grp := backend.NewGroup(
				db,
				entity.Groupname,
				backend.Reverse(entity.Post),
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
				return rfc7231.StatusConflict(he.NetPrintf("a group with that name already exists"))
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
	//sess := req.Things["session"].(*backend.Session)
	grp := backend.GetGroupByID(db, name)
	/*if grp.Read != 1 && !backend.IsSubscribed(db, sess.UserID, *grp) {
		return nil
	}*/
	return (*group)(grp)
}

func (d dirGroups) SubentityNotFound(name string, req he.Request) he.Response {
	return rfc7231.StatusNotFound(nil)
}
