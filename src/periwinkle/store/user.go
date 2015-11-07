// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal
// Copyright 2015 Luke Shumaker

package store

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	he "httpentity"
	"httpentity/util" // heutil
	"io"
	"jsonpatch"
	"log"
	"strings"
)

var _ he.Entity = &User{}
var _ he.NetEntity = &User{}
var dirUsers he.Entity = newDirUsers()

// Model /////////////////////////////////////////////////////////////

type User struct {
	Id        string        `json:"user_id"`
	FullName  string        `json:"fullname"`
	PwHash    []byte        `json:"-"`
	Addresses []UserAddress `json:"addresses"`
}

func (o User) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

type UserAddress struct {
	// TODO: add a "verified" boolean
	Id            int64          `json:"-"`
	UserId        string         `json:"-"`
	Medium        string         `json:"medium"`
	Address       string         `json:"address"`
	Subscriptions []Subscription `json:"-"`
}

func (o UserAddress) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("medium", "media(id)", "RESTRICT", "RESTRICT").
		AddUniqueIndex("uniqueness_idx", "medium", "address").
		Error
}

func (addr UserAddress) AsEmailAddress() string {
	if addr.Medium == "email" {
		return addr.Address
	} else {
		return addr.Address + "@" + addr.Medium + ".gateway"
	}
}

func (u *User) populate(db *gorm.DB) {
	db.Model(u).Related(&u.Addresses)
	address_ids := make([]int64, len(u.Addresses))
	for i, address := range u.Addresses {
		address_ids[i] = address.Id
	}
	var subscriptions []Subscription
	if len(address_ids) > 0 {
		if result := db.Where("address_id IN (?)", address_ids).Find(&subscriptions); result.Error != nil {
			if !result.RecordNotFound() {
				panic(result.Error)
			}
		}
	} else {
		subscriptions = make([]Subscription, 0)
	}
	for i := range u.Addresses {
		u.Addresses[i].Subscriptions = []Subscription{}
		for _, subscription := range subscriptions {
			if u.Addresses[i].Id == subscription.AddressId {
				u.Addresses[i].Subscriptions = append(u.Addresses[i].Subscriptions, subscription)
			}
		}
	}
}

func GetUserById(db *gorm.DB, id string) *User {
	id = strings.ToLower(id)
	var o User
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	o.populate(db)
	return &o
}

func GetUserByAddress(db *gorm.DB, medium string, address string) *User {
	var o User
	result := db.Joins("inner join user_addresses on user_addresses.user_id=users.id").Where("user_addresses.medium=? and user_addresses.address=?", medium, address).Find(&o)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	o.populate(db)
	return &o
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	u.PwHash = hash
	return err
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PwHash, []byte(password))
	return err == nil
}

func NewUser(db *gorm.DB, name string, password string, email string) User {
	if name == "" {
		panic("name can't be empty")
	}
	o := User{
		Id:        name,
		FullName:  "",
		Addresses: []UserAddress{{Medium: "email", Address: email}},
	}
	if err := o.SetPassword(password); err != nil {
		panic(err)
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return o
}

func (o *User) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}

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
	out_ops := make([]patchop, 0, len(ops))
	checkedpass := false
	for _, op := range ops {
		if op.Path == "/password" {
			switch op.Op {
			case "test":
				if !o.CheckPassword(op.Value) {
					ret := he.StatusConflict(heutil.NetString("old password didn't match"))
					return &ret
				}
				checkedpass = true
			case "replace":
				if !checkedpass {
					ret := he.StatusUnsupportedMediaType(heutil.NetString("you must submit and old password (using 'test') before setting a new one"))
					return &ret
				}
				if o.CheckPassword(op.Value) {
					ret := he.StatusConflict(heutil.NetString("that new password is the same as the old one"))
					return &ret
				}
				o.SetPassword(op.Value)
			default:
				ret := he.StatusUnsupportedMediaType(heutil.NetString("you may only 'set' or 'replace' the password"))
				return &ret
			}
		} else {
			out_ops = append(out_ops, op)
		}
	}
	str, err = json.Marshal(out_ops)
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
			sess := req.Things["session"].(*Session)
			if sess.UserId != user.Id {
				return he.StatusForbidden(heutil.NetString("Unauthorized user"))
			}
			var new_user User
			httperr := safeDecodeJSON(req.Entity, &new_user)
			if httperr != nil {
				return *httperr
			}
			if user.Id != new_user.Id {
				return he.StatusConflict(heutil.NetString("Cannot change user id"))
			}
			// TODO: this won't play nice with the
			// password hash (because it's private), or
			// with addresses (because the (private) IDs
			// need to be made to match up)
			*user = new_user
			user.Save(db)
			return he.StatusOK(user)
		},
		"PATCH": func(req he.Request) he.Response {
			db := req.Things["db"].(*gorm.DB)
			sess := req.Things["session"].(*Session)
			if sess.UserId != user.Id {
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
			var new_user User
			err := patch.Apply(user, &new_user)
			if err != nil {
				return he.StatusConflict(heutil.NetString(err.Error()))
			}
			if user.Id != new_user.Id {
				return he.StatusConflict(heutil.NetString("Cannot change user id"))
			}
			// some mucking around with private fields to make things match up
			new_user.PwHash = user.PwHash
			delete_address_ids := []int64{}
			for o := range user.Addresses {
				old_addr := &user.Addresses[o]
				match := false
				for n := range new_user.Addresses {
					new_addr := &new_user.Addresses[n]
					if new_addr.Medium == old_addr.Medium && new_addr.Address == old_addr.Address {
						new_addr.Id = old_addr.Id
						match = true
					}
				}
				if !match {
					delete_address_ids = append(delete_address_ids, old_addr.Id)
				}
			}
			// save

			*user = new_user
			user.Save(db)
			if len(delete_address_ids) > 0 {
				if err = db.Where("id IN (?)", delete_address_ids).Delete(UserAddress{}).Error; err != nil {
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

			user := NewUser(db, entity.Username, entity.Password, entity.Email)
			req.Things["user"] = user
			return he.StatusCreated(r, user.Id, req)
		},
	}
	return r
}

func (d t_dirUsers) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirUsers) Subentity(name string, req he.Request) he.Entity {
	name = strings.ToLower(name)
	sess := req.Things["session"].(*Session)
	if sess == nil && req.Method == "POST" {
		user, ok := req.Things["user"].(User)
		if !ok {
			return nil
		}
		if user.Id == name {
			return &user
		}
		return nil
	} else if sess.UserId != name {
		return nil
	}
	db := req.Things["db"].(*gorm.DB)
	return GetUserById(db, name)
}
