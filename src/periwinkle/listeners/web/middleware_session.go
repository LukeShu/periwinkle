// Copyright 2015 Luke Shumaker

package web

import (
	he "httpentity"
	"periwinkle/store"
	"time"
)

type session struct{}

func (p session) Before(req *he.Request) {
	//db := req.Things["db"].(store.Db)
	var sess *store.Session = nil /* TODO */
	if sess != nil {
		sess.Last_used = time.Now()
	}
	req.Things["session"] = sess
}

func (p session) After(req he.Request, res *he.Response) {
	sess, ok := req.Things["session"].(*store.Session)
	if ok && sess != nil {
		sess.Save()
	}
}
