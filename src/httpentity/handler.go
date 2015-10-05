// Copyright 2015 Luke Shumaker

package httpentity

import (
	"net/http"
	"strings"
)

type netHttpHandler struct {
	prefix string
	root   Entity
}

func NetHttpHandler(prefix string, entity Entity) http.Handler {
	return netHttpHandler{prefix: prefix, root: entity}
}

func (h netHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Adapt the request from `net/http` format to `httpentity` format
	req := Request{
		Headers: r.Header,
		Query:   r.URL.Query(),
		Entity:  nil,
	}
	ReadEntity(r.Body, r.Header.Get("Content-Type"), &req.Entity)

	// Run the request
	res := Route(h.prefix, h.root, req, r.Method, r.URL)

	// Adapt the response from `httpentity` format to `net/http` format
	for k, v := range res.Headers {
		w.Header().Set(k, strings.Join(v, ", "))
	}
	w.WriteHeader(int(res.Status))
	res.WriteEntity(w)
}
