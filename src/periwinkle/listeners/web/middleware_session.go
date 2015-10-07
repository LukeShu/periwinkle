// Copyright 2015 Luke Shumaker

package web

import (
	he "httpentity"
	"periwinkle/store"
	"time"
)

type session struct{}

func (p session) Before(req *he.Request) {
	db := req.Things["db"].(store.DB)
	var sess *store.Session = nil

	var session_id1 string
	var session_id2 string

	hash    , ok := req.Entity.(map[string]interface{}); if !ok { goto end }
	session_id1, ok = hash["session_id"].(string); if !ok { goto end }
	delete(hash, "session_id")
	session_id2 = "" // TODO: get from req.Headers cookies

	if session_id1 != session_id2 {
		goto end
	}

	sess = store.GetSessionById(db, session_id1)
end:
	if sess != nil {
		sess.LastUsed = time.Now()
	}
	req.Things["session"] = sess
}

func (p session) After(req he.Request, res *he.Response) {
	db := req.Things["db"].(store.DB)
	sess, ok := req.Things["session"].(*store.Session)
	if ok && sess != nil {
		sess.Save(db)
	}
}
