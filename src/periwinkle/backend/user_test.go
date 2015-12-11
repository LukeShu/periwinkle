// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend_test

import (
	"periwinkle"
	. "periwinkle/backend"
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {

		user := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		switch {
		case !strings.EqualFold(user.ID, "JohnDoe"):
			t.Error("User ID was not properly set to.")
		case user.FullName != "":
			t.Error("User name was not properly set.")
		case user.Addresses[0].Address != "johndoe@purdue.edu":
			t.Error("User address was not preperly set.")
		}

	})
}

func TestGetUserByID(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		user := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		o := GetUserByID(tx, user.ID)

		switch {
		case o == nil:
			t.Error("GetUserByID() returned nil")
		case !strings.EqualFold(user.ID, o.ID):
			t.Error("GetUserByID() returned a user with a different ID")
		}
	})
}

func TestNewUserAddress(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {

		user := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		newAddr := NewUserAddress(tx, user.ID, "email", "johndoe2@purdue.edu", false)

		switch {
		case newAddr.Address != "johndoe2@purdue.edu":
			t.Error("Error adding new email to user in NewUserAddress()")
		case newAddr.Medium != "email":
			t.Error("Error assigning medium type in NewUserAddress()")
		}

		newAddr = NewUserAddress(tx, user.ID, "sms", "7655555555", false)

		switch {
		case newAddr.Address != "7655555555":
			t.Error("Error adding new sms to user in NewUserAddress()")
		case newAddr.Medium != "sms":
			t.Error("Error assigning medium type in NewUserAddress()")
		}
	})
}

func TestGetUserByAddress(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		user := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		o := GetUserByAddress(tx, "email", user.Addresses[0].Address)
		if !strings.EqualFold(user.ID, o.ID) {
			t.Error("Error in GetUserByAdress()")
		}
	})
}

// func TestSetPassword(t *testing.T) {
// 	t.Error("TODO")
// }

// func TestCheckPassword(t *testing.T) {
// 	t.Error("TODO")
// }

// func TestGetAddressByUserAndMedium(t *testing.T) {
// 	addr := GetAddressByUserAndMedium(tx, user.ID, "email")

// 	switch {
// 	case addr == nil:
// 		t.Error("GetAddressByUserAndMedium() returned nil")
// 	case addr.Address != user.Addresses[0].Address:
// 		t.Error("Addresses do not match: " + user.Addresses[0].Address + " != " + addr.Address)
// 	}
// }

func TestGetUserSubscriptions(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		waste := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		user := GetUserByID(tx, waste.ID)

		subs := user.GetUserSubscriptions(tx)

		if subs == nil {
			t.Error("GetUserSubscriptions returned nil")
		}
	})
}
