// Copyright 2015 Luke Shumaker

package main

import (
	"github.com/jinzhu/gorm"
	he "httpentity"
	"periwinkle/store"
	"time"
)

func MiddlewareSession(req he.Request, handle func(he.Request) he.Response) he.Response {
	cookie := req.Cookie("session_id")
	if cookie == nil {
		return handle(req)
	}
	session_id := cookie.Value

	switch req.Method {
	case "OPTIONS", "GET", "HEAD":
		// do nothing
	default:
		header := req.Headers.Get("X-XSRF-TOKEN")
		if header != session_id {
			return handle(req)
		}
	}

	// It's not worth panicing if we have database errors here.
	if db, dbok := req.Things["db"].(*gorm.DB); dbok {
		sess := store.GetSessionById(db, session_id)
		if sess != nil {
			sess.LastUsed = time.Now()
			req.Things["session"] = sess
			func() {
				defer recover()
				sess.Save(db)
			}()
		}
	}
	return handle(req)
}
