// Copyright 2015 Luke Shumaker

package httpapi

import he "httpentity"

type dirRoot struct {
	methods     map[string]func(he.Request) he.Response
	dirCaptchas he.Entity
	dirGroups   he.Entity
	dirMessages he.Entity
	fileSession he.Entity
	dirUsers    he.Entity
}

func NewDirRoot() he.RootEntity {
	return &dirRoot{
		methods:     make(map[string]func(he.Request) he.Response),
		dirCaptchas: newDirCaptchas(),
		dirGroups:   newDirGroups(),
		dirMessages: newDirMessages(),
		fileSession: newFileSession(),
		dirUsers:    newDirUsers(),
	}
}

func (d dirRoot) Methods() map[string]func(he.Request) he.Response {
	return d.methods
}

func (d dirRoot) Subentity(name string, request he.Request) he.Entity {
	switch name {
	case "captcha":
		return d.dirCaptchas
	case "groups":
		return d.dirGroups
	case "msgs":
		return d.dirMessages
	case "session":
		return d.fileSession
	case "users":
		return d.dirUsers
	}
	return nil
}

func (d dirRoot) SubentityNotFound(name string, request he.Request) he.Response {
	panic("TODO")
}

func (d dirRoot) MethodNotAllowed(request he.Request) he.Response {
	panic("TODO")
}
