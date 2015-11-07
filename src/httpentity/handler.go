// Copyright 2015 Luke Shumaker

package httpentity

import (
	"httpentity/util"
	"net/url"
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

type middlewareHolder struct {
	middleware  Middleware
	nextHandler func(request Request, u *url.URL) Response
}

func (mwh middlewareHolder) handler(request Request, u *url.URL) Response {
	return mwh.middleware(request, u, mwh.nextHandler)
}

// assumes that the url has already been passed to normalizeURL()
func (router *Router) defaultHandler(request Request, u *url.URL) Response {
	entity := findEntity(router.Root, request, strings.TrimPrefix(u.Path, router.Prefix))
	if entity == nil {
		return statusNotFound()
	}

	callmethod := request.Method
	if callmethod == "HEAD" {
		callmethod = "GET"
	}
	methods := entity.Methods()
	handler, method_allowed := methods[request.Method]

	var response Response
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

	// make sure the Location: header is absolute
	if l := response.Headers.Get("Location"); l != "" {
		u2, _ := u.Parse(l)
		response.Headers.Set("Location", u2.String())
		// XXX: this is pretty hacky, because it is tightly
		// integrated with the entity format used by
		// (Request).StatusCreated()
		if response.Status == 201 {
			ilist := []interface{}(response.Entity.(heutil.NetList))
			slist := make([]string, len(ilist))
			for i, iface := range ilist {
				slist[i] = iface.(string)
			}
			response.Entity = extensions2net(u2, slist)
		}
	}

	return response
}

func (router *Router) initHandler() {
	handler := router.defaultHandler
	for i := 0; i < len(router.Middlewares); i++ {
		handler = middlewareHolder{
			middleware:  router.Middlewares[len(router.Middlewares)-1-i],
			nextHandler: handler,
		}.handler
	}
	router.handler = handler
}
