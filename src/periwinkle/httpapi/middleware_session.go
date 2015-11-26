// Copyright 2015 Luke Shumaker

package httpapi

import (
	he "httpentity"
	"periwinkle/backend"
	"time"

	"github.com/jinzhu/gorm"
)

func getsession(req he.Request) *backend.Session {
	cookie := req.Cookie("session_id")
	if cookie == nil {
		return nil
	}
	sessionID := cookie.Value

	switch req.Method {
	case "OPTIONS", "GET", "HEAD":
		// do nothing
	default:
		header := req.Headers.Get("X-XSRF-TOKEN")
		if header != sessionID {
			return nil
		}
	}

	// It's not worth panicing if we have database errors here.
	db, ok := req.Things["db"].(*gorm.DB)
	if !ok {
		return nil
	}
	sess := backend.GetSessionByID(db, sessionID)
	if sess != nil {
		sess.LastUsed = time.Now()
		func() {
			defer recover()
			sess.Save(db)
		}()
	}
	return sess
}

var MiddlewareSession = he.Middleware{
	Outside: func(req he.Request, handle func(he.Request) he.Response) he.Response {
		req.Things["session"] = getsession(req)
		return handle(req)
	},
}
