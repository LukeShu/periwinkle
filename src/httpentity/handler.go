// Copyright 2015 Luke Shumaker

package httpentity

import (
	"net/http"
	"strings"
	"fmt"
)

type Middleware func(req *Request)

type netHttpHandler struct {
	prefix string
	root   Entity
	middle []Middleware
}

func NetHttpHandler(prefix string, entity Entity, middlewares ...Middleware) http.Handler {
	return netHttpHandler{prefix: prefix, root: entity, middle: middlewares}
}

func (h netHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var res Response
	// Adapt the request from `net/http` format to `httpentity` format
	req := Request{
		Headers: r.Header,
		Query:   r.URL.Query(),
		Entity:  nil,
	}
	switch r.Method {
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
		middleware(&req)
	}

	// Run the request
	res = Route(h.prefix, h.root, req, r.Method, r.URL)
end:
	// Adapt the response from `httpentity` format to `net/http` format
	for k, v := range res.Headers {
		w.Header().Set(k, strings.Join(v, ", "))
	}
	w.WriteHeader(int(res.status))
	res.WriteEntity(w)
}
