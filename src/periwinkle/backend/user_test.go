// Copyright 2015 Davis Webb

package backend_test

import (
	"periwinkle"
	. "periwinkle/backend"
	"periwinkle/cfg"
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {

	conf := CreateTempDB()

	user := NewUser(conf.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	switch {
	case !strings.EqualFold(user.ID, "JohnDoe"):
		t.Error("User ID was not properly set to.")
	case user.FullName != "":
		t.Error("User name was not properly set.")
	case user.Addresses[0].Address != "johndoe@purdue.edu":
		t.Error("User address was not preperly set.")
	}
}

func TestGetUserByID(t *testing.T) {
	conf := CreateTempDB()

	user := NewUser(conf.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	o := GetUserByID(conf.DB, user.ID)

	switch {
	case o == nil:
		t.Error("GetUserByID() returned nil")
	case !strings.EqualFold(user.ID, o.ID):
		t.Error("GetUserByID() returned a user with a different ID")
	}
}

func TestNewUserAddress(t *testing.T) {

	conf := CreateTempDB()

	user := NewUser(conf.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	newAddr := NewUserAddress(conf.DB, user.ID, "email", "johndoe2@purdue.edu", false)

	switch {
	case strings.Compare(newAddr.Address, "johndoe2@purdue.edu") != 0:
		t.Error("Error adding new email to user in NewUserAddress()")
	case strings.Compare(newAddr.Medium, "email") != 0:
		t.Error("Error assigning medium type in NewUserAddress()")
	}

	newAddr = NewUserAddress(conf.DB, user.ID, "sms", "7655555555", false)

	switch {
	case strings.Compare(newAddr.Address, "7655555555") != 0:
		t.Error("Error adding new sms to user in NewUserAddress()")
	case strings.Compare(newAddr.Medium, "sms") != 0:
		t.Error("Error assigning medium type in NewUserAddress()")
	}

}

func TestGetUserByAddress(t *testing.T) {
	conf := CreateTempDB()

	user := NewUser(conf.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	o := GetUserByAddress(conf.DB, "email", user.Addresses[0].Address)
	if strings.Compare(user.ID, o.ID) != 0 {
		t.Error("Error in GetUserByAdress()")
	}
}

// func TestSetPassword(t *testing.T) {
// 	t.Error("TODO")
// }

// func TestCheckPassword(t *testing.T) {
// 	t.Error("TODO")
// }

// func TestGetAddressByIDAndMedium(t *testing.T) {
// 	t.Error("TODO")
// }

// func TestGetUserSubscriptions(t *testing.T) {
// 	t.Error("TODO")
// }

func CreateTempDB() *periwinkle.Cfg {
	conf := periwinkle.Cfg{
		Mailstore:      "./Maildir",
		WebUIDir:       "./www",
		Debug:          true,
		TrustForwarded: true,
		GroupDomain:    "localhost",
		WebRoot:        "locahost:8080",
		DB:             nil, // the default DB is set later
	}

	db, err := cfg.OpenDB("sqlite3", "file:temp.sqlite?mode=memory&_txlock=exclusive")
	if err != nil {
		periwinkle.Logf("Error loading sqlite3 database")
	}
	conf.DB = db

	DbSchema(conf.DB)

	return &conf
}
