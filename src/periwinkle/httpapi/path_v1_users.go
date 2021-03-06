// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal
// Copyright 2015 Luke Shumaker

package httpapi

import (
	"encoding/json"
	he "httpentity"
	"httpentity/rfc7231"
	"jsonpatch"
	"periwinkle"
	"periwinkle/backend"
)

var _ he.NetEntity = &user{}
var _ he.EntityGroup = &user{}
var _ he.EntityGroup = &dirUsers{}

type user backend.User

func (o *user) backend() *backend.User { return (*backend.User)(o) }

// Model /////////////////////////////////////////////////////////////

func (o *user) Subentity(name string, req he.Request) he.Entity {
	if name == "subscriptions" {
		return &userSubscriptions{*o, "", nil}
	}
	return nil
}

func (d *user) SubentityNotFound(name string, req he.Request) he.Response {
	return rfc7231.StatusNotFound(nil)
}

func (o *user) patchPassword(patch *jsonpatch.Patch) *he.Response {
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
					ret := rfc7231.StatusConflict(he.NetPrintf("old password didn't match"))
					return &ret
				}
				checkedpass = true
			case "replace":
				if !checkedpass {
					ret := rfc7231.StatusUnsupportedMediaType(he.NetPrintf("you must submit and old password (using 'test') before setting a new one"))
					return &ret
				}
				if o.backend().CheckPassword(op.Value) {
					ret := rfc7231.StatusConflict(he.NetPrintf("that new password is the same as the old one"))
					return &ret
				}
				o.backend().SetPassword(op.Value)
			default:
				ret := rfc7231.StatusUnsupportedMediaType(he.NetPrintf("you may only 'set' or 'replace' the password"))
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

func (usr *user) Methods() map[string]func(he.Request) he.Response {
	return map[string]func(he.Request) he.Response{
		"GET": func(req he.Request) he.Response {
			var addresses []backend.UserAddress
			for _, addr := range usr.Addresses {
				if addr.Medium != "noop" && addr.Medium != "admin" {
					addresses = append(addresses, addr)
				}
			}
			usr.Addresses = addresses
			return rfc7231.StatusOK(usr)
		},
		"PUT": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			sess := req.Things["session"].(*backend.Session)
			if sess.UserID != usr.ID {
				return rfc7231.StatusForbidden(he.NetPrintf("Unauthorized user"))
			}
			var newUser user
			httperr := safeDecodeJSON(req.Entity, &newUser)
			if httperr != nil {
				return *httperr
			}
			if usr.ID != newUser.ID {
				return rfc7231.StatusConflict(he.NetPrintf("Cannot change user id"))
			}
			// TODO: this won't play nice with the
			// password hash (because it's private), or
			// with addresses (because the (private) IDs
			// need to be made to match up)
			*usr = newUser
			usr.backend().Save(db)
			return rfc7231.StatusOK(usr)
		},
		"PATCH": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			sess := req.Things["session"].(*backend.Session)
			if sess.UserID != usr.ID {
				return rfc7231.StatusForbidden(he.NetPrintf("Unauthorized user"))
			}
			patch, ok := req.Entity.(jsonpatch.Patch)
			if !ok {
				return rfc7231.StatusUnsupportedMediaType(he.NetPrintf("PATCH request must have a patch media type"))
			}
			httperr := usr.patchPassword(&patch)
			if httperr != nil {
				return *httperr
			}
			var newUser user
			err := patch.Apply(usr, &newUser)
			if err != nil {
				return rfc7231.StatusConflict(he.ErrorToNetEntity(409, err))
			}
			if usr.ID != newUser.ID {
				return rfc7231.StatusConflict(he.NetPrintf("Cannot change user id"))
			}
			if newUser.PwHash == nil || len(newUser.PwHash) == 0 {
				newUser.PwHash = usr.PwHash
			}
			*usr = newUser
			usr.backend().Save(db)
			return rfc7231.StatusOK(usr)
		},
		"DELETE": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
			usr.backend().Delete(db)
			return rfc7231.StatusNoContent()
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *user) Encoders() map[string]he.Encoder {
	return defaultEncoders(o)
}

// Directory ("Controller") //////////////////////////////////////////

type dirUsers struct {
	methods map[string]func(he.Request) he.Response
}

func newDirUsers() dirUsers {
	r := dirUsers{}
	r.methods = map[string]func(he.Request) he.Response{
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(*periwinkle.Tx)
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
				return rfc7231.StatusUnsupportedMediaType(he.NetPrintf("username, email, and password can't be emtpy"))
			}

			if entity.PasswordVerification != "" {
				if entity.Password != entity.PasswordVerification {
					// Passwords don't match
					return rfc7231.StatusConflict(he.NetPrintf("password and password_verification don't match"))
				}
			}

			usr := backend.NewUser(db, entity.Username, entity.Password, entity.Email)
			backend.NewUserAddress(db, usr.ID, "noop", backend.RandomString(20), true)
			backend.NewUserAddress(db, usr.ID, "admin", backend.RandomString(20), true)
			req.Things["user"] = usr
			return rfc7231.StatusCreated(r, usr.ID, req)
		},
	}
	return r
}

func (d dirUsers) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d dirUsers) Subentity(name string, req he.Request) he.Entity {
	// sess := req.Things["session"].(*backend.Session)
	// if sess == nil {
	// 	if req.Method == "POST" {
	// 		usr, ok := req.Things["user"].(backend.User)
	// 		if ok && strings.EqualFold(usr.ID, name) {
	// 			return (*user)(&usr)
	// 		}
	// 	}
	// 	return nil
	// } else if sess.UserID != name {
	// 	return nil
	// }
	db := req.Things["db"].(*periwinkle.Tx)
	return (*user)(backend.GetUserByID(db, name))
}

func (d dirUsers) SubentityNotFound(name string, req he.Request) he.Response {
	return rfc7231.StatusNotFound(nil)
}
