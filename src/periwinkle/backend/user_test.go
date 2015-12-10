// Copyright 2015 Davis Webb

package backend

import (
	"periwinkle"
	"testing"

	"github.com/jinzhu/gorm"
)

var cfg *periwinkle.Cfg // I figure we dont need to make a DB each time
var user User           // same goes for our JohnDoe

func TestNewUser(t *testing.T) {

	t.Log("Starting User Tests")

	cfg = CreateTempDB()

	user = NewUser(cfg.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	switch {
	case user.ID != "JohnDoe":
		t.Error("User ID was not properly set to.")
	case user.FullName != "":
		t.Error("User name was not properly set.")
	case user.Addresses[0].Address != "johndoe@purdue.edu":
		t.Error("User address was not preperly set.")
	}
}

func TestGetUserByID(t *testing.T) {

	o := GetUserByID(cfg.DB, user.ID)

	switch {
	case o == nil:
		t.Error("GetUserByID() returned nil")
	case user.ID != o.ID:
		t.Error("Error in GetUserByID()")
	}
}

func TestNewUserAddress(t *testing.T) {

	newAddr := NewUserAddress(cfg.DB, user.ID, "email", "johndoe2@purdue.edu", false)

	switch {
	case newAddr.Address != "johndoe2@purdue.edu":
		t.Error("Error adding new email to user in NewUserAddress()")
	case newAddr.Medium != "email":
		t.Error("Error assigning medium type in NewUserAddress()")
	}

	newAddr = NewUserAddress(cfg.DB, user.ID, "sms", "7655555555", false)

	switch {
	case newAddr.Address != "7655555555":
		t.Error("Error adding new sms to user in NewUserAddress()")
	case newAddr.Medium != "sms":
		t.Error("Error assigning medium type in NewUserAddress()")
	}
}

func TestGetUserByAddress(t *testing.T) {

	o := GetUserByAddress(cfg.DB, "email", user.Addresses[0].Address)
	if user.ID != o.ID {
		t.Error("Error in GetUserByAdress()")
	}
}

func TestSetPassword(t *testing.T) {
	t.Log("TODO")
}

func TestCheckPassword(t *testing.T) {
	t.Log("TODO")
}

func TestGetAddressByIDAndMedium(t *testing.T) {
	addr := GetAddressByIDAndMedium(cfg.DB, user.ID, "email")

	switch {
	case addr == nil:
		t.Error("GetAddressByIDAndMedium() returned nil")
	case addr.Address != user.Addresses[0].Address:
		t.Error("Addresses do not match: " + user.Addresses[0].Address + " != " + addr.Address)
	}
	cfg.DB.Close()
}

func TestGetUserSubscriptions(t *testing.T) {
	subs := user.GetUserSubscriptions(cfg.DB)
	if subs == nil {
		t.Error("GetUserSubscriptions returned nil")
	}
}

func TestCloseUserDB(t *testing.T) {
	t.Log("Finishing User Tests")
	cfg.DB.Close()
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
