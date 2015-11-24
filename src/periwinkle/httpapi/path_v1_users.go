// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal
// Copyright 2015 Luke Shumaker

package httpapi

import (
	"encoding/json"
	he "httpentity"
	"httpentity/heutil"
	"io"
	"jsonpatch"
	"periwinkle/backend"
	"strings"

	"github.com/jinzhu/gorm"
)

var _ he.Entity = &User{}
var _ he.NetEntity = &User{}
var dirUsers he.Entity = newDirUsers()

type User backend.User

func (o *User) backend() *backend.User { return (*backend.User)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *User) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *User) patchPassword(patch *jsonpatch.Patch) *he.Response {
	// this is in the running for the grossest code I've ever
	// written, but I think it's the best way to do it --lukeshu
	type patchop struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}
	str, err := json.Marshal(patch)
	if err != nil {
		panic(err)
	}
	var ops []patchop
	err = json.Unmarshal(str, &ops)
	if err != nil {
		return nil
	}
	outOps := make([]patchop, 0, len(ops))
	checkedpass := false
	for _, op := range ops {
		if op.Path == "/password" {
			switch op.Op {
			case "test":
				if !o.backend().CheckPassword(op.Value) {
					ret := he.StatusConflict(heutil.NetString("old password didn't match"))
					return &ret
				}
				checkedpass = true
			case "replace":
				if !checkedpass {
					ret := he.StatusUnsupportedMediaType(heutil.NetString("you must submit and old password (using 'test') before setting a new one"))
					return &ret
				}
				if o.backend().CheckPassword(op.Value) {
					ret := he.StatusConflict(heutil.NetString("that new password is the same as the old one"))
					return &ret
				}
				o.backend().SetPassword(op.Value)
			default:
				ret := he.StatusUnsupportedMediaType(heutil.NetString("you may only 'set' or 'replace' the password"))
				return &ret
			}
		} else {
			outOps = append(outOps, op)
		}
	}
	str, err = json.Marshal(outOps)
	if err != nil {
		panic(err)
	}
	var out jsonpatch.JSONPatch
	err = json.Unmarshal(str, &out)
	if err != nil {
		panic(out)
	}
	*patch = out
	return nil
}

func (user *User) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			return he.StatusOK(user)
		},
		"PUT": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			if sess.UserID != user.ID {
				return he.StatusForbidden(heutil.NetString("Unauthorized user"))
			}
			var newUser User
			httperr := safeDecodeJSON(req.Entity, &newUser)
			if httperr != nil {
				return *httperr
			}
			if user.ID != newUser.ID {
				return he.StatusConflict(heutil.NetString("Cannot change user id"))
			}
			// TODO: this won't play nice with the
			// password hash (because it's private), or
			// with addresses (because the (private) IDs
			// need to be made to match up)
			*user = newUser
			user.backend().Save(db)
			return he.StatusOK(user)
		},
		"PATCH": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*backend.Session)
			if sess.UserID != user.ID {
				return he.StatusForbidden(heutil.NetString("Unauthorized user"))
			}
			patch, ok := req.Entity.(jsonpatch.Patch)
			if !ok {
				return he.StatusUnsupportedMediaType(heutil.NetString("PATCH request must have a patch media type"))
			}
			httperr := user.patchPassword(&patch)
			if httperr != nil {
				return *httperr
			}
			var newUser User
			err := patch.Apply(user, &newUser)
			if err != nil {
				return he.StatusConflict(heutil.NetString(err.Error()))
			}
			if user.ID != newUser.ID {
				return he.StatusConflict(heutil.NetString("Cannot change user id"))
			}
			// some mucking around with private fields to make things match up
			newUser.PwHash = user.PwHash
			deleteAddressIDs := []int64{}
			for o := range user.Addresses {
				oldAddr := &user.Addresses[o]
				match := false
				for n := range newUser.Addresses {
					newAddr := &newUser.Addresses[n]
					if newAddr.Medium == oldAddr.Medium && newAddr.Address == oldAddr.Address {
						newAddr.ID = oldAddr.ID
						match = true
					}
				}
				if !match {
					deleteAddressIDs = append(deleteAddressIDs, oldAddr.ID)
				}
			}
			// save

			*user = newUser
			user.backend().Save(db)
			if len(deleteAddressIDs) > 0 {
				if err = db.Where("id IN (?)", deleteAddressIDs).Delete(backend.UserAddress{}).Error; err != nil {
					panic(err)
				}
			}
			return he.StatusOK(user)
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			db.Delete(user)
			return he.StatusNoContent()
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *User) Encoders() map[string]func(io.Writer) error {
	return defaultEncoders(o)
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirUsers struct {
	methods map[string]func(he.Request) he.Response
}

func newDirUsers() t_dirUsers {
	r := t_dirUsers{}
	r.methods = map[string]func(he.Request) he.Response{
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			type postfmt struct {
				Username             string `json:"username"`
				Email                string `json:"email"`
				Password             string `json:"password"`
				PasswordVerification string `json:"password_verification,omitempty"`
			}
			var entity postfmt
			httperr := safeDecodeJSON(req.Entity, &entity)
			if httperr != nil {
				return *httperr
			}

			if entity.Username == "" || entity.Email == "" || entity.Password == "" {
				return he.StatusUnsupportedMediaType(heutil.NetString("username, email, and password can't be emtpy"))
			}

			if entity.PasswordVerification != "" {
				if entity.Password != entity.PasswordVerification {
					// Passwords don't match
					return he.StatusConflict(heutil.NetString("password and password_verification don't match"))
				}
			}

			entity.Username = strings.ToLower(entity.Username)

			user := backend.NewUser(db, entity.Username, entity.Password, entity.Email)
			backend.NewUserAddress(db, user.ID, "noop", "", true)
			backend.NewUserAddress(db, user.ID, "admin", "", true)
			req.Things["user"] = user
			return he.StatusCreated(r, user.ID, req)
		},
	}
	return r
}

func (d t_dirUsers) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirUsers) Subentity(name string, req he.Request) he.Entity {
	name = strings.ToLower(name)
	sess := req.Things["session"].(*backend.Session)
	if sess == nil && req.Method == "POST" {
		user, ok := req.Things["user"].(User)
		if !ok {
			return nil
		}
		if user.ID == name {
			return &user
		}
		return nil
	} else if sess.UserID != name {
		return nil
	}
	db := req.Things["db"].(*gorm.DB)
	return (*User)(backend.GetUserByID(db, name))
}
