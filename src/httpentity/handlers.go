// Copyright 2015 Luke Shumaker

package httpentity

import (
	"strings"
)

type middlewareHolder struct {
	middleware  Middleware
	nextOutside func(request Request) Response
	nextInside  func(request Request, entity Entity) Response
}

func (mwh middlewareHolder) outsideHandler(request Request) Response {
	return mwh.middleware.Outside(request, mwh.nextOutside)
}

func (mwh middlewareHolder) insideHandler(request Request, entity Entity) Response {
	return mwh.middleware.Inside(request, entity, mwh.nextInside)
}

// assumes that the request.URL has already been passed to normalizeURL()
func (router *Router) defaultOutsideHandler(request Request) Response {
	entity, notFound := router.findEntity(strings.TrimPrefix(request.URL.Path, router.Prefix), request)
	if entity == nil {
		return notFound
	}
	return router.insideHandler(request, entity)
}

func (router *Router) defaultInsideHandler(request Request, entity Entity) Response {
	methods := entity.Methods()
	handler, methodAllowed := methods[request.Method]

	if !methodAllowed {
		var response Response
		if extra, ok := entity.(EntityExtra); ok {
			response = extra.MethodNotAllowed(request)
		} else {
			response = router.MethodNotAllowed(request, request.URL)
		}
		return response
	}
	return handler(request)
}

func (router *Router) initHandlers() {
	outsideHandler := router.defaultOutsideHandler
	insideHandler := router.defaultInsideHandler
	for i := 0; i < len(router.Middlewares); i++ {
		holder := middlewareHolder{
			middleware:  router.Middlewares[len(router.Middlewares)-1-i],
			nextOutside: outsideHandler,
			nextInside:  insideHandler,
		}
		if holder.middleware.Outside != nil {
			outsideHandler = holder.outsideHandler
		}
		if holder.middleware.Inside != nil {
			insideHandler = holder.insideHandler
		}
	}
	router.outsideHandler = outsideHandler
	router.insideHandler = insideHandler
}
