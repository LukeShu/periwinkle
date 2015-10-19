// Copyright 2015 Luke Shumaker

package store

import he "httpentity"

var DirRoot he.Entity = newDirRoot()

type t_dirRoot struct {
	methods map[string]func(he.Request) he.Response
}

func newDirRoot() t_dirRoot {
	return t_dirRoot{methods: make(map[string]func(he.Request) he.Response)}
}

func (d t_dirRoot) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d t_dirRoot) Subentity(name string, request he.Request) he.Entity {
	switch name {
	case "captcha":
		return dirCaptchas
	case "groups":
		return dirGroups
	case "msgs":
		return dirMessages
	case "s":
		return dirShortUrls
	case "session":
		return fileSession
	case "users":
		return dirUsers
	}
	return nil
}
