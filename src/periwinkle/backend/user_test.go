// Copyright 2015 Davis Webb

package backend

import (
	"periwinkle"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestNewUser(t *testing.T) {

	cfg := CreateTempDB()

	user := NewUser(cfg.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	switch {
	case strings.Compare(user.ID, "JohnDoe") != 0:
		t.Error("User ID was not properly set to.")
	case strings.Compare(user.FullName, "") != 0:
		t.Error("User name was not properly set.")
	case strings.Compare(user.Addresses[0].Address, "johndoe@purdue.edu") != 0:
		t.Error("User address was not preperly set.")
	}
}

func TestGetUserByID(t *testing.T) {
	cfg := CreateTempDB()

	user := NewUser(cfg.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	o := GetUserByID(cfg.DB, user.ID)

	switch {
	case o == nil:
		t.Error("GetUserByID() returned nil")
	case strings.Compare(user.ID, o.ID) != 0:
		t.Error("Error in GetUserByID()")
	}
}

func TestNewUserAddress(t *testing.T) {

	cfg := CreateTempDB()

	user := NewUser(cfg.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	newAddr := NewUserAddress(cfg.DB, user.ID, "email", "johndoe2@purdue.edu", false)

	switch {
	case strings.Compare(newAddr.Address, "johndoe2@purdue.edu") != 0:
		t.Error("Error adding new email to user in NewUserAddress()")
	case strings.Compare(newAddr.Medium, "email") != 0:
		t.Error("Error assigning medium type in NewUserAddress()")
	}

	newAddr = NewUserAddress(cfg.DB, user.ID, "sms", "7655555555", false)

	switch {
	case strings.Compare(newAddr.Address, "7655555555") != 0:
		t.Error("Error adding new sms to user in NewUserAddress()")
	case strings.Compare(newAddr.Medium, "sms") != 0:
		t.Error("Error assigning medium type in NewUserAddress()")
	}

}

func TestGetUserByAddress(t *testing.T) {
	cfg := CreateTempDB()

	user := NewUser(cfg.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	o := GetUserByAddress(cfg.DB, "email", user.Addresses[0].Address)
	if strings.Compare(user.ID, o.ID) != 0 {
		t.Error("Error in GetUserByAdress()")
	}
}

func TestSetPassword(t *testing.T) {
	t.Error("TODO")
}

func TestCheckPassword(t *testing.T) {
	t.Error("TODO")
}

func TestGetAddressByIDAndMedium(t *testing.T) {
	t.Error("TODO")
}

func TestGetUserSubscriptions(t *testing.T) {
	t.Error("TODO")
}

func CreateTempDB() *periwinkle.Cfg {
	cfg := periwinkle.Cfg{
		Mailstore:      "./Maildir",
		WebUIDir:       "./www",
		Debug:          true,
		TrustForwarded: true,
		GroupDomain:    "localhost",
		WebRoot:        "locahost:8080",
		DB:             nil, // the default DB is set later
	}

	db, err := gorm.Open("sqlite3", "file:temp.sqlite?mode=memory&_txlock=exclusive")

	if err != nil {
		periwinkle.Logf("Error loading sqlite3 database")
	}

	cfg.DB = &db

	DbSchema(cfg.DB)

	return &cfg
}
