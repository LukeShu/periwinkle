// Copyright 2015 Luke Shumaker

package web

import (
	he "httpentity"
)

type postHack struct{}

func (p postHack) Before(req *he.Request) {
	hash, ok := req.Entity.(map[string]interface{})
	if !ok {
		return
	}
	method, ok := hash["_method"].(string)
	delete(hash, "_method")
	if !ok {
		return
	}
	switch method {
	case "POST", "PUT", "PATCH":
		// do nothing
	default:
		req.Entity = nil
	}
	req.Method = method
}

func (p postHack) After(req he.Request, res *he.Response) {}
