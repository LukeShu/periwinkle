// Copyright 2015 Luke Shumaker

package main

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"periwinkle/store"
	"time"
)

type session struct{}

func (p session) Before(req *he.Request) {
	var sess *store.Session = nil
	defer func() {
		if sess != nil {
			sess.LastUsed = time.Now()
		}
		req.Things["session"] = sess
	}()

	cookie := req.Cookie("session_id")
	if cookie == nil {
		return
	}
	session_id := cookie.Value

	switch req.Method {
	case "OPTIONS", "GET", "HEAD":
		// do nothing
	default:
		header := req.Headers.Get("X-XSRF-TOKEN")
		if header != session_id {
			return
		}
	}

	// It's not worth panicing if we have database errors here.
	db, ok := req.Things["db"].(*gorm.DB)
	if !ok {
		return
	}
	sess = store.GetSessionById(db, session_id)
	if sess != nil {
		func() {
			defer recover()
			sess.Save(db)
		}()
	}
}

func (p session) After(req he.Request, res *he.Response) {}
