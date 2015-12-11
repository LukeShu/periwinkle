// Copyright 2015 Davis Webb

package backend_test

import (
	"periwinkle"
	. "periwinkle/backend"
	"testing"
)

func TestNewMessage(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		existence := []int{2, 2}
		read := []int{2, 2}
		post := []int{1, 1, 1}
		join := []int{1, 1, 1}

		group := NewGroup(tx, "Purdue", existence, read, post, join)

		msg := NewMessage(tx, "420.420@example.com", *group, "unique")

		switch {
		case msg.ID == "":
			t.Error("NewMessage(returned nil)")
		case msg.Unique != "unique":
			t.Error("NewMessage(did not properly created a message)")
		}
	})
}

func TestGetMessageByID(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		existence := []int{2, 2}
		read := []int{2, 2}
		post := []int{1, 1, 1}
		join := []int{1, 1, 1}

		group := NewGroup(tx, "Purdue", existence, read, post, join)
		NewMessage(tx, "420.420@example.com", *group, "unique")

		o := GetMessageByID(tx, "420.420@example.com")

		switch {
		case o.ID == "":
			t.Error("NewMessage(returned nil)")
		case o.ID != "420.420@example.com":
			t.Error("NewMessage(did not properly grab the message)")
		}
	})
}
