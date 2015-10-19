// Copyright 2015 Luke Shumaker

package httpentity

import (
	"fmt"
	"strings"
)

func methods2string(methods map[string]func(request Request) Response) string {
	set := make(map[string]bool, len(methods)+2)
	for method := range methods {
		set[method] = true
	}
	set["OPTIONS"] = true
	if _, get := set["GET"]; get {
		set["HEAD"] = true
	}
	list := make([]string, len(set))
	i := uint(0)
	for method := range set {
		list[i] = method
		i++
	}
	return strings.Join(list, ", ")
}

// Takes the normalized path without the leading slash
func route(entity Entity, req Request, upath string) (ret Response) {
	if entity == nil {
		ret = req.statusNotFound()
	} else if upath == "" {
		callmethod := req.Method
		if callmethod == "HEAD" {
			callmethod = "GET"
		}
		methods := entity.Methods()
		handler, method_allowed := methods[req.Method]
		if method_allowed {
			ret = handler(req)
		} else {
			ret = req.statusMethodNotAllowed(methods2string(methods))
		}
		if callmethod == "OPTIONS" {
			ret.Status = 200
			ret.Headers.Set("Allow", methods2string(methods))
		}
	} else {
		parts := strings.SplitN(upath, "/", 2)
		if len(parts) != 2 {
			panic(fmt.Sprintf("path parser logic failure: %#v", upath))
		}
		child := parts[0]
		grandchildren := parts[1]
		ret = route(entity.Subentity(child, req), req, grandchildren)
	}
	return
}
