// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker
// Copyright 2015 Guntas Grewal

package httpapi

import (
	he "httpentity"
	"httpentity/rfc7231"
	"jsonpatch"
	"periwinkle"
	"periwinkle/backend"
)

var _ he.EntityGroup = &group{}
var _ he.NetEntity = &group{}
var _ he.EntityGroup = &dirGroups{}

type group backend.Group

type Enumerategroup struct {
	Groupname     string                 `json:"groupname"`
	Post          map[string]string      `json:"post"`
	Join          map[string]string      `json:"join"`
	Read          map[string]string      `json:"read"`
	Existence     map[string]string      `json:"existence"`
	Subscriptions []backend.Subscription `json:"subscriptions"`
}

func EnumerateGroup(o *group) Enumerategroup {
	var enum Enumerategroup
	enum.Groupname = o.ID
	exist := [...]int{o.ExistencePublic, o.ExistenceConfirmed}
	enum.Existence = backend.ReadExist(exist)
	read := [...]int{o.ReadPublic, o.ReadConfirmed}
	enum.Read = backend.ReadExist(read)
	post := [...]int{o.PostPublic, o.PostConfirmed, o.PostMember}
	enum.Post = backend.PostJoin(post)
	join := [...]int{o.JoinPublic, o.JoinConfirmed, o.JoinMember}
	enum.Join = backend.PostJoin(join)
	enum.Subscriptions = o.Subscriptions
	return enum
}

func RenumerateGroup(entity Enumerategroup) group {
	read := make([]int, 2)
	existence := make([]int, 2)
	post := make([]int, 3)
	join := make([]int, 3)

	existence = backend.Reverse(entity.Existence)
	read = backend.Reverse(entity.Read)
	post = backend.Reverse(entity.Post)
	join = backend.Reverse(entity.Join)

	o := group{
		ID:                 entity.Groupname,
		ReadPublic:         read[0],
		ReadConfirmed:      read[1],
		ExistencePublic:    existence[0],
		ExistenceConfirmed: existence[1],
		PostPublic:         post[0],
		PostConfirmed:      post[1],
		PostMember:         post[2],
		JoinPublic:         join[0],
		JoinConfirmed:      join[1],
		JoinMember:         join[2],
		Subscriptions:      entity.Subscriptions,
	}
	return o
}

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
			type EnumerateGroup struct {
				Groupname     string                 `json:"groupname"`
				Post          map[string]string      `json:"post"`
				Join          map[string]string      `json:"join"`
				Read          map[string]string      `json:"read"`
				Existence     map[string]string      `json:"existence"`
				Subscriptions []backend.Subscription `json:"subscriptions"`
			}

			var enum EnumerateGroup
			enum.Groupname = o.ID
			exist := [...]int{o.ExistencePublic, o.ExistenceConfirmed}
			enum.Existence = backend.ReadExist(exist)
			read := [...]int{o.ReadPublic, o.ReadConfirmed}
			enum.Read = backend.ReadExist(read)
			post := [...]int{o.PostPublic, o.PostConfirmed, o.PostMember}
			enum.Post = backend.PostJoin(post)
			join := [...]int{o.JoinPublic, o.JoinConfirmed, o.JoinMember}
			enum.Join = backend.PostJoin(join)
			enum.Subscriptions = o.Subscriptions
			return rfc7231.StatusOK(he.NetJSON{Data: enum})
		},
		"PUT": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)

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
			db := req.Things["db"].(*periwinkle.Tx)
			sess := req.Things["session"].(*backend.Session)
			subscribed := backend.IsSubscribed(db, sess.UserID, *o.backend())
			if !backend.IsAdmin(db, sess.UserID, *o.backend()) {
				if o.JoinPublic == 1 {
					if subscribed == 0 {
						return rfc7231.StatusForbidden(he.NetPrintf("Unauthorized user"))
					}
					if o.JoinConfirmed == 1 && subscribed == 1 {
						return rfc7231.StatusForbidden(he.NetPrintf("Unauthorized user"))
					}
					if o.JoinMember == 1 {
						return rfc7231.StatusForbidden(he.NetPrintf("Unauthorized user"))
					}
				}
			}
			enum := EnumerateGroup(o)
			var newGroup Enumerategroup
			patch, ok := req.Entity.(jsonpatch.Patch)
			if !ok {
				return rfc7231.StatusUnsupportedMediaType(he.NetPrintf("PATCH request must have a patch media type"))
			}
			err := patch.Apply(enum, &newGroup)
			if err != nil {
				return rfc7231.StatusConflict(he.NetPrintf("%v", err))
			}
			if o.ID != newGroup.Groupname {
				return rfc7231.StatusConflict(he.NetPrintf("Cannot change group id"))
			}

			*o = RenumerateGroup(newGroup)
			o.backend().Save(db)
			return rfc7231.StatusOK(o)
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			sess := req.Things["session"].(*backend.Session)
			if !backend.IsAdmin(db, sess.UserID, *o.backend()) {
				return rfc7231.StatusForbidden(he.NetPrintf("Unauthorized user"))
			}
			o.backend().Delete(db)
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
			db := req.Things["db"].(*periwinkle.Tx)
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
				Groupname     string                 `json:"groupname"`
				Post          map[string]string      `json:"post"`
				Join          map[string]string      `json:"join"`
				Read          map[string]string      `json:"read"`
				Existence     map[string]string      `json:"existence"`
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
				enum.Post = backend.PostJoin(post)
				join := [...]int{grp.JoinPublic, grp.JoinConfirmed, grp.JoinMember}
				enum.Join = backend.PostJoin(join)
				enum.Subscriptions = grp.Subscriptions
				data[i] = enum
			}
			return rfc7231.StatusOK(he.NetJSON{Data: data})
		},
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			type Response1 struct {
				Groupname string            `json:"groupname"`
				Post      map[string]string `json:"post"`
				Join      map[string]string `json:"join"`
				Read      map[string]string `json:"read"`
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

			grp := backend.NewGroup(
				db,
				entity.Groupname,
				backend.Reverse(entity.Existence),
				backend.Reverse(entity.Read),
				backend.Reverse(entity.Post),
				backend.Reverse(entity.Join),
			)
			sess := req.Things["session"].(*backend.Session)
			address := backend.GetAddressesByUserAndMedium(db, sess.UserID, "noop")[0]
			backend.NewSubscription(db, address.ID, grp.ID, true)
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
	db := req.Things["db"].(*periwinkle.Tx)
	sess := req.Things["session"].(*backend.Session)
	grp := backend.GetGroupByID(db, name)
	if grp.ReadPublic == 1 {
		subscribed := backend.IsSubscribed(db, sess.UserID, *grp)
		if (grp.ReadConfirmed == 1 && subscribed == 1) || subscribed == 0 {
			return nil
		}
	}
	return (*group)(grp)
}

func (d dirGroups) SubentityNotFound(name string, req he.Request) he.Response {
	return rfc7231.StatusNotFound(nil)
}
