// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal
// Copyright 2015 Luke Shumaker

package orm

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string
	FullName string
	pwHash   []byte
}

func getUserById(con DB, id int) (*User, error) {
	var user User
	err := con.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user)
	switch {
	case err == sql.ErrNoRows:
		// user does not exist
		return nil, nil
	case err != nil:
		// error talking to the DB
		return nil, err
	default:
		// all ok
		return &user, nil
	}
}

func GetUserByName(con DB, name string) (*User, error) {
	var user User
	err := con.QueryRow("SELECT * FROM users WHERE name=?", name).Scan(&user)
	switch {
	case err == sql.ErrNoRows:
		// user does not exist
		return nil, nil
	case err != nil:
		// error talking to the DB
		return nil, err
	default:
		// all ok
		return &user, nil
	}
}

func GetUserByEmail(con DB, address string) (*User, error) {
	var user User
	err := con.QueryRow(""+
		"SELECT users.* "+
		"FROM users JOIN user_address ON users.id=user_address.user_id "+
		"WHERE user_address.address=? AND user_address=?",
		address,
		"email").Scan(&user)
	switch {
	case err == sql.ErrNoRows:
		// user does not exist
		return nil, nil
	case err != nil:
		// error talking to the DB
		return nil, err
	default:
		// all ok
		return &user, nil
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

func NewUser(con DB, name string, password string) (u *User, err error) {
	u = &User{
		Id:       name,
		FullName: "",
	}
	u.SetPassword(password)
	_, err = con.Exec("INSERT INTO users VALUES (?,?)", u.Id, u.FullName, u.pwHash)
	if err != nil {
		u = nil
	}
	return
}

func (u *User) Save() error {
	// TODO
	panic("not implemented")
}
