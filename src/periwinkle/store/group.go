// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

import (
	"database/sql"
	"github.com/jmoiron/modl"
	he "httpentity"
	"strings"
)

var _ he.Entity = &Group{}
var _ he.NetEntity = &Group{}
var dirGroups he.Entity = newDirGroups()

// Model /////////////////////////////////////////////////////////////

type Group struct {
	Id        string
	Addresses []GroupAddress
}

func (g *Group) init(con modl.SqlExecutor) error {
	return con.Select(&g.Addresses, "SELECT * FROM group_addresses WHERE group_id=?", g.Id)
}

type GroupAddress struct {
	Id      int64
	groupId string
	Medium  string
	Address string
}

func GetGroupById(con modl.SqlExecutor, name string) *Group {
	var group Group
	err := con.Get(&group, name)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		group.init(con)
		return &group
	}
}

func GetGroupByAddress(con modl.SqlExecutor, medium string, address string) *Group {
	panic("TODO: ORM")
}

func NewGroup(con modl.SqlExecutor, name string) *Group {
	g := &Group{Id: name}
	err := con.Insert(g)
	if err != nil {
		panic(err)
	}
	return g
}

func (o *Group) Subentity(name string, req he.Request) he.Entity {
	panic("TODO: API: (*Group).Subentity()")
}

func (o *Group) Methods() map[string]he.Handler {
	return map[string]he.Handler{
		"GET": func(req he.Request) he.Response {
			panic("TODO: API: (*Group).Methods()[\"GET\"]")
		},
		"PUT": func(req he.Request) he.Response {
			panic("TODO: API: (*Group).Methods()[\"PUT\"]")
		},
		"PATCH": func(req he.Request) he.Response {
			panic("TODO: API: (*Group).Methods()[\"PATCH\"]")
		},
		"DELETE": func(req he.Request) he.Response {
			panic("TODO: API: (*Group).Methods()[\"DELETE\"]")
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *Group) Encoders() map[string]he.Encoder {
	return defaultEncoders(o)
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirGroups struct {
	methods map[string]he.Handler
}

func newDirGroups() t_dirGroups {
	r := t_dirGroups{}
	r.methods = map[string]he.Handler{
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(modl.SqlExecutor)
			badbody := req.StatusBadRequest("submitted body not what expected")
			hash, ok := req.Entity.(map[string]interface{}); if !ok { return badbody }
			groupname, ok := hash["groupname"].(string)    ; if !ok { return badbody }

			groupname = strings.ToLower(groupname)

			group := NewGroup(db, groupname)
			if group == nil {
				return req.StatusConflict(he.NetString("a group with that name already exists"))
			} else {
				return req.StatusCreated(r, groupname)
			}
		},
	}
	return r
}

func (d t_dirGroups) Methods() map[string]he.Handler {
	return d.methods
}

func (d t_dirGroups) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(modl.SqlExecutor)
	return GetGroupById(db, name)
}
