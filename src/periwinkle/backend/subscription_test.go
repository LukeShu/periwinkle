// Copyright 2015 Davis Webb

package backend_test

import (
	"periwinkle"
	. "periwinkle/backend"
	"strings"
	"testing"
)

func TestNewSession(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {

		user := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		sess := NewSession(tx, &user, "password")

		switch {
		case sess == nil:
			t.Error("NewSession(returned nil)")
		case !strings.EqualFold(sess.UserID, "JohnDoe"):
			t.Error("NewSession(returned wrong user ID)")
		}

	})
}

func TestGetSessionByID(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {

		user := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		sess := NewSession(tx, &user, "password")

		o := GetSessionByID(tx, sess.ID)

		switch {
		case o == nil:
			t.Error("GetSessionByID(returned nil)")
		case !strings.EqualFold(o.UserID, "JohnDoe"):
			t.Error("GetSessionByID(returned wrong user ID)")
		}

	})
}
