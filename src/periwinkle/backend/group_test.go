// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend_test

import (
	"periwinkle"
	. "periwinkle/backend"
	"strings"
	"testing"
)

func TestNewGroup(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {

		u1 := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		u2 := NewUser(tx, "JaneDoe", "password", "janedoe@purdue.edu")

		sub := []Subscription{{Address: u1.Addresses[0], Confirmed: true}, {Address: u2.Addresses[0], Confirmed: true}}

		existence := []int{2, 2}
		read := []int{2, 2}
		post := []int{1, 1, 1}
		join := []int{1, 1, 1}

		group := NewGroup(tx, "The Doe", existence, read, post, join)

		group.Subscriptions = sub

		switch {
		case !strings.EqualFold("The Doe", group.ID):
			t.Error("ID's do not match")
		}
	})
}

func TestGetGroupByID(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {

		existence := []int{2, 2}
		read := []int{2, 2}
		post := []int{1, 1, 1}
		join := []int{1, 1, 1}

		group := NewGroup(tx, "The Doe", existence, read, post, join)

		o := GetGroupByID(tx, "The Doe")

		switch {
		case o == nil:
			t.Error("GetGroupByID: returned nil")
		case !strings.EqualFold(o.ID, group.ID):
			t.Error("ID does not match requested group")
		}
	})
}

func TestGetGroupsByMember(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		u1 := NewUser(tx, "JohnDoe", "password", "johndoe@purdue.edu")

		u2 := NewUser(tx, "JaneDoe", "password", "janedoe@purdue.edu")

		existence := []int{2, 2}
		read := []int{2, 2}
		post := []int{1, 1, 1}
		join := []int{1, 1, 1}

		err := tx.Create(&Group{
			ID:                 "Purdue",
			ReadPublic:         read[0],
			ReadConfirmed:      read[1],
			ExistencePublic:    existence[0],
			ExistenceConfirmed: existence[1],
			PostPublic:         post[0],
			PostConfirmed:      post[1],
			PostMember:         post[2],
			JoinPublic:         join[0],
			JoinConfirmed:      join[1],
			JoinMember:         join[2],
			Subscriptions: []Subscription{{
				Address:   u1.Addresses[0],
				Confirmed: true,
			},

				{Address: u2.Addresses[0],
					Confirmed: true,
				},
			},
		}).Error
		if err != nil {
			t.Error("Issue creating group")
		}

		u := GetUserByID(tx, u1.ID)

		o := GetGroupsByMember(tx, *u)

		switch {
		case o == nil:
			t.Error("GetGroupsByMember: returned nil")
		case !strings.EqualFold(o[0].ID, "purdue"):
			t.Error("Did not grab correct group")
		}
	})
}

// func TestGetPublicAndSubscribedGroups(t *testing.T) {
// 	t.Log("TODO")
// }

func TestGetAllGroups(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		existence := []int{2, 2}
		read := []int{2, 2}
		post := []int{1, 1, 1}
		join := []int{1, 1, 1}

		group := NewGroup(tx, "The Doe", existence, read, post, join)
		NewGroup(tx, "g2", existence, read, post, join)
		NewGroup(tx, "g3", existence, read, post, join)
		NewGroup(tx, "g4", existence, read, post, join)

		o := GetAllGroups(tx)

		// TODO: not guaranteed to be in this order
		switch {
		case o == nil:
			t.Error("GetAllGroups(returned nil)")
		case !strings.EqualFold(o[0].ID, group.ID):
			t.Error(o[0].ID + " != " + group.ID)
		case !strings.EqualFold(o[1].ID, "g2"):
			t.Error(o[1].ID + " != " + "g2")
		case !strings.EqualFold(o[2].ID, "g3"):
			t.Error(o[2].ID + " != " + "g3")
		case !strings.EqualFold(o[3].ID, "g4"):
			t.Error(o[3].ID + " != " + "g4")
		}
	})
}
