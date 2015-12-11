// Copyright 2015 Davis Webb

package backend_test

import (
	"net/url"
	"periwinkle"
	. "periwinkle/backend"
	"testing"
)

func TestNewShortURL(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		var u *url.URL

		u, _ = url.Parse("www.example.com")

		url := NewShortURL(tx, u)

		switch {
		case url == nil:
			t.Error("NewShortURL(returned nil)")
		case url.ID == "":
			t.Error("NewShortURL(Did not set ID)")
		case url.Dest == "":
			t.Error("NewShortURL(Did not set Dest)")
		}
	})
}

func TestGetShortURLByID(t *testing.T) {
	conf := CreateTempDB()
	conf.DB.Do(func(tx *periwinkle.Tx) {
		var u *url.URL

		u, _ = url.Parse("www.example.com")

		url := NewShortURL(tx, u)

		o := GetShortURLByID(tx, url.ID)

		switch {
		case o == nil:
			t.Error("GetShortURLByID(returned nil)")
		case o.ID != url.ID:
			t.Error("GetShortURLByID(Did not get the right shortURL)")
		}
	})
}
