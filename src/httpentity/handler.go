// Copyright 2015 Luke Shumaker

package httpentity

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type Middleware interface {
	Before(req *Request)
	After(req Request, res *Response)
}

type netHttpHandler struct {
	prefix string
	root   Entity
	middle []Middleware
	debug  bool
}

func NetHttpHandler(debug bool, prefix string, entity Entity, middlewares ...Middleware) http.Handler {
	return netHttpHandler{prefix: prefix, root: entity, middle: middlewares, debug: debug}
}

func (h netHttpHandler) serveHTTP(w http.ResponseWriter, r *http.Request) (res Response) {
	// Adapt the request from `net/http` format to `httpentity` format
	req := Request{
		Scheme:  r.URL.Scheme,
		Method:  r.Method,
		Headers: r.Header,
		Query:   r.URL.Query(),
		Entity:  nil,
		Things:  map[string]interface{}{},
	}
	// just in case anything goes wrong, don't bring down the
	// process.
	defer func() {
		if err := recover(); err != nil {
			reason := err
			if h.debug {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				reason = fmt.Sprintf("%v\n\n%s", string(buf))
			}
			res = req.statusInternalServerError(reason)
		}
	}()
	// parse the submitted entity
	switch req.Method {
	case "POST", "PUT", "PATCH":
		entity, err := ReadEntity(r.Body, r.Header.Get("Content-Type"))
		if entity == nil {
			if err == nil {
				res = req.statusUnsupportedMediaType()
			} else {
				res = req.StatusBadRequest(fmt.Sprintf("reading request body: %s", err))
			}
			return
		} else {
			req.Entity = entity
		}
	}
	for _, middleware := range h.middle {
		middleware.Before(&req)
		defer middleware.After(req, &res)
	}

	// Run the request
	res = Route(h.prefix, h.root, req, req.Method, r.URL)
	return
}

func (h netHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := h.serveHTTP(w, r)
	// Adapt the response from `httpentity` format to `net/http` format
	for k, v := range res.Headers {
		w.Header().Set(k, strings.Join(v, ", "))
	}
	w.WriteHeader(int(res.status))
	res.WriteEntity(w)
}
