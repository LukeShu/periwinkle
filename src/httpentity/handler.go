// Copyright 2015 Luke Shumaker

package httpentity

import (
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

func middlewareToHandler(mw Middleware, h func(Request) Response) func(Request) Response {
	return func(r Request) Response {
		return mw(r, h)
	}
}

func entityToHandler(entity Entity) func(Request) Response {
	return func(request Request) (response Response) {
		callmethod := request.Method
		if callmethod == "HEAD" {
			callmethod = "GET"
		}
		methods := entity.Methods()
		handler, method_allowed := methods[request.Method]
		if method_allowed {
			response = handler(request)
		} else {
			if callmethod == "OPTIONS" {
				response = StatusOK(nil)
			} else {
				response = statusMethodNotAllowed(methods2string(methods))
			}
		}
		if callmethod == "OPTIONS" {
			response.Headers.Set("Allow", methods2string(methods))
		}
		return
	}
}

func (router *Router) handleEntity(entity Entity, request Request) Response {
	handler := entityToHandler(entity)
	for i := 0; i < len(router.Middlewares); i++ {
		handler = middlewareToHandler(
			router.Middlewares[len(router.Middlewares)-1-i],
			handler)
	}
	return handler(request)
}
