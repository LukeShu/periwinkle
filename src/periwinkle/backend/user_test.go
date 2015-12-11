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

func TestGetAddressByUserAndMedium(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {

		user := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		addrs := GetAddressesByUserAndMedium(tx, user.ID, "email")

		if len(addrs) != len(user.Addresses) {
			t.Error("Number of addresses does not match")
		}
		for _, addr1 := range addrs {
			match := false
			for _, addr2 := range user.Addresses {
				if addr1.Address == addr2.Address {
					match = true
				}
			}
			if !match {
				t.Error("Addresses had no match:" + addr1.Address)
			}
		}
	})
}

func TestUserGetSubscriptions(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		o := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		user := GetUserByID(tx, o.ID)

		subs := user.GetSubscriptions(tx)

		if subs == nil {
			t.Error("(*User).GetSubscriptions returned nil")
		}
	})
}
