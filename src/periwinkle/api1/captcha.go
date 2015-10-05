// Copyright 2015 Luke Shumaker

package api1

import (
	he "httpentity"
	//"github.com/dchest/captcha"
	//"orm"
)

type t_dirCaptcha struct {
	methods map[string]he.Handler
}

var dirCaptcha he.Entity = newDirCaptcha()

func newDirCaptcha() t_dirCaptcha {
	r := t_dirCaptcha{methods: make(map[string]he.Handler)}
	r.methods["POST"] = r.newCaptcha
	return r
}

func (d t_dirCaptcha) Methods() map[string]he.Handler {
	return d.methods
}

func (d t_dirCaptcha) Subentity(name string, request he.Request) he.Entity {
	panic("not implemented")
}

func (d t_dirCaptcha) newCaptcha(req he.Request) he.Response {
	panic("not implemented")
}
