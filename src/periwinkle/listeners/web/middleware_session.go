// Copyright 2015 Luke Shumaker

package web

import (
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

	hash, ok := req.Entity.(map[string]interface{}); if !ok { return }
	session_id1, ok := hash["session_id"].(string) ; if !ok { return }
	delete(hash, "session_id")
	cookie := req.Cookie("session_id")
	session_id2 := ""
	if cookie != nil {
		session_id2 = cookie.Value
	}

	if session_id1 != session_id2 {
		return
	}

	db, ok := req.Things["db"].(store.DB); if !ok { return }
	sess = store.GetSessionById(db, session_id1)
	sess.Save(db)
}

func (p session) After(req he.Request, res *he.Response) {}
