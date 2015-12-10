// Copyright 2015 Davis Webb

package backend

import (
	"testing"
)

var u1 User // JohnDoe
var u2 User // JaneDoe

var group *Group

func TestNewGroup(t *testing.T) {

	t.Log("Starting Group Tests")

	cfg = CreateTempDB()

	u1 = NewUser(cfg.DB, "JohnDoe", "password", "johndoe@purdue.edu")

	u2 = NewUser(cfg.DB, "JaneDoe", "password", "janedoe@purdue.edu")

	sub := []Subscription{{Address: u1.Addresses[0], Confirmed: true}, {Address: u2.Addresses[0], Confirmed: true}}

	existence := []int{2, 2}
	read := []int{2, 2}
	post := []int{1, 1, 1}
	join := []int{1, 1, 1}

	group = NewGroup(cfg.DB, "The Doe", existence, read, post, join)

	group.Subscriptions = sub

	switch {
	case "The Doe" != group.ID:
		t.Error("ID's do not match")
	}
}

func TestGetGroupByID(t *testing.T) {

	o := GetGroupByID(cfg.DB, "The Doe")

	switch {
	case o == nil:
		t.Error("GetGroupByID: returned nil")
	case o.ID != group.ID:
		t.Error("ID does not match requested group")
	}
}

func TestGetGroupsByMember(t *testing.T) {

	o := GetGroupsByMember(cfg.DB, u1)

	switch {
	case o == nil:
		t.Error("GetGroupsByMember: returned nil")
	case o[0].ID != group.ID:
		t.Error("Did not grab correct group")
	}
}

func TestGetPublicAndSubscribedGroups(t *testing.T) {
	t.Log("TODO")
}

func TestGetAllGroups(t *testing.T) {

	existence := []int{2, 2}
	read := []int{2, 2}
	post := []int{1, 1, 1}
	join := []int{1, 1, 1}

	NewGroup(cfg.DB, "g2", existence, read, post, join)
	NewGroup(cfg.DB, "g3", existence, read, post, join)
	NewGroup(cfg.DB, "g4", existence, read, post, join)

	o := GetAllGroups(cfg.DB)

	switch {
	case o == nil:
		t.Error("GetAllGroups(returned nil)")
	case o[0].ID != group.ID:
		t.Error(o[0].ID + " != " + group.ID)
	case o[1].ID != "g2":
		t.Error(o[1].ID + " != " + "g2")
	case o[2].ID != "g3":
		t.Error(o[2].ID + " != " + "g3")
	case o[3].ID != "g4":
		t.Error(o[3].ID + " != " + "g4")
	}
}

func TestCloseGroupDB(t *testing.T) {
	t.Log("Finishing Group Tests")
	cfg.DB.Close()
}
