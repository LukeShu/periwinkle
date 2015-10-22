// Copyright 2015 Luke Shumaker

package web

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

	if req.Entity != nil {
		hash, ok := req.Entity.(map[string]interface{})
		if !ok {
			return
		}

		session_id_body, ok := hash["session_id"].(string)
		delete(hash, "session_id")
		if !ok {
			return
		}

		if session_id_body != session_id {
			return
		}
	}

	// It's not worth panicing if we have database errors here.
	db, ok := req.Things["db"].(*gorm.DB)
	if !ok {
		return
	}
	sess = store.GetSessionById(db, session_id)
	if sess == nil {
		return
	}
	sess.Save(db)
}

func (p session) After(req he.Request, res *he.Response) {}
