// Copyright 2015 Luke Shumaker

package main

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"net/url"
	"periwinkle/store"
	"time"
)

func getsession(req he.Request) *store.Session {
	cookie := req.Cookie("session_id")
	if cookie == nil {
		return nil
	}
	session_id := cookie.Value

	switch req.Method {
	case "OPTIONS", "GET", "HEAD":
		// do nothing
	default:
		header := req.Headers.Get("X-XSRF-TOKEN")
		if header != session_id {
			return nil
		}
	}

	// It's not worth panicing if we have database errors here.
	db, ok := req.Things["db"].(*gorm.DB)
	if !ok {
		return nil
	}
	sess := store.GetSessionById(db, session_id)
	if sess != nil {
		sess.LastUsed = time.Now()
		func() {
			defer recover()
			sess.Save(db)
		}()
	}
	return sess
}

func MiddlewareSession(req he.Request, u *url.URL, handle func(he.Request, *url.URL) he.Response) he.Response {
	req.Things["session"] = getsession(req)
	return handle(req, u)
}
