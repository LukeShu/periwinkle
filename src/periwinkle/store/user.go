// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal
// Copyright 2015 Luke Shumaker

package store

import (
	"database/sql"
	"github.com/jmoiron/modl"
	"golang.org/x/crypto/bcrypt"
	he "httpentity"
	"strings"
)

var _ he.Entity = &User{}
var _ he.NetEntity = &User{}
var dirUsers he.Entity = newDirUsers()

// Model /////////////////////////////////////////////////////////////

type User struct {
	Id        string
	FullName  string
	pwHash    []byte
	Addresses []UserAddress
}

func (u *User) init(con modl.SqlExecutor) error {
	return con.Select(&u.Addresses, "SELECT * FROM user_addresses WHERE user_id=?", u.Id)
}

type UserAddress struct {
	id      int64
	userId  string
	Medium  string
	Address string
}

func GetUserById(con modl.SqlExecutor, id string) *User {
	var user User
	err := con.Get(&user, id)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		user.init(con)
		return &user
	}
}

func GetUserByAddress(con modl.SqlExecutor, medium string, address string) *User {
	var user User
	err := con.SelectOne(&user,
		"SELECT users.* "+
			"FROM users JOIN user_address ON users.id=user_address.user_id "+
			"WHERE user_address.address=? AND user_address=?",
		address,
		"email")
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		user.init(con)
		return &user
	}
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	u.pwHash = hash
	return err
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.pwHash, []byte(password))
	return err != nil
}

func NewUser(con modl.SqlExecutor, name string, password string, email string) *User {
	u := &User{
		Id:       name,
		FullName: "",
	}
	err := u.SetPassword(password)
	if err != nil {
		panic(err)
	}
	err = con.Insert(u)
	if err != nil {
		panic(err)
	}
	a := &UserAddress{
		userId:  u.Id,
		Medium:  "email",
		Address: email,
	}
	con.Insert(a)
	if err != nil {
		panic(err)
	}
	u.init(con)
	return u
}

func (u *User) Save(con modl.SqlExecutor) error {
	_, err := con.Update(u)
	return err
}

func (o *User) Subentity(name string, req he.Request) he.Entity {
	return nil
}

func (o *User) Methods() map[string]he.Handler {
	return map[string]he.Handler{
		"GET": func(req he.Request) he.Response {
			badbody := req.StatusBadRequest("submitted body not what expected")
			sess, ok := req.Things["session"].(*Session); if !ok { return badbody }
			if !ok || sess.UserId != o.Id {
				return req.StatusUnauthorized(o)
			}
			return req.StatusOK(o)
		},
		"PUT": func(req he.Request) he.Response {
			panic("TODO: API: (*User).Methods()[\"PUT\"]")
		},
		"PATCH": func(req he.Request) he.Response {
			panic("TODO: API: (*User).Methods()[\"PATCH\"]")
		},
		"DELETE": func(req he.Request) he.Response {
			panic("TODO: API: (*User).Methods()[\"DELETE\"]")
		},
	}
}

// View //////////////////////////////////////////////////////////////

func (o *User) Encoders() map[string]he.Encoder {
	return defaultEncoders(o)
}

// Directory ("Controller") //////////////////////////////////////////

type t_dirUsers struct {
	methods map[string]he.Handler
}

func newDirUsers() t_dirUsers {
	r := t_dirUsers{}
	r.methods = map[string]he.Handler{
		"POST": func(req he.Request) he.Response {
			db := req.Things["db"].(modl.SqlExecutor)
			badbody := req.StatusBadRequest("submitted body not what expected")
			hash, ok := req.Entity.(map[string]interface{}); if !ok { return badbody }
			username, ok := hash["username"].(string)      ; if !ok { return badbody }
			email   , ok := hash["email"].(string)         ; if !ok { return badbody }
			password, ok := hash["password"].(string)      ; if !ok { return badbody }

			if password2, ok := hash["password_verification"].(string); ok {
				if password != password2 {
					// Passwords don't match
					return req.StatusConflict(he.NetString("password and password_verification don't match"))
				}
			}

			username = strings.ToLower(username)

			user := NewUser(db, username, password, email)
			if user == nil {
				return req.StatusConflict(he.NetString("either that username or password is already taken"))
			} else {
				return req.StatusCreated(r, username)
			}
		},
	}
	return r
}

func (d t_dirUsers) Methods() map[string]he.Handler {
	return d.methods
}

func (d t_dirUsers) Subentity(name string, req he.Request) he.Entity {
	db := req.Things["db"].(modl.SqlExecutor)
	return GetUserById(db, name)
}
