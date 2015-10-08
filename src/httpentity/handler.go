// Copyright 2015 Luke Shumaker

package httpentity

import (
	"net/http"
	"strings"
	"fmt"
	"runtime"
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
		Scheme: r.URL.Scheme,
		Method: r.Method,
		Headers: r.Header,
		Query:   r.URL.Query(),
		Entity:  nil,
	}
	// just in case anything goes wrong, don't bring down the
	// process.
	defer func() {
		if r := recover(); r != nil {
			if h.debug {
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)
				stack = stack[0:n]
				r = fmt.Sprintf("%v\n\n%s", string(stack))
			}
			res = req.statusInternalServerError(r)
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
			goto end
		} else {
			req.Entity = entity
		}
	}
	for _, middleware := range h.middle {
		middleware.Before(&req)
	}

	// Run the request
	res = Route(h.prefix, h.root, req, req.Method, r.URL)
end:
	// this loop has go go backwards over h.middle
	var i int
	for i = len(h.middle)-1; i >= 0; i-- {
		middleware := h.middle[i]
		middleware.After(req, &res)
	}
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
