// Copyright 2015 Luke Shumaker

package httpentity

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (r *Router) Route(req Request, u *url.URL) (res Response) {
	u, mimetype := normalizeURL(u)
	if mimetype != "" {
		// the file extension overrides the Accept: header
		req.Headers.Set("Accept", mimetype)
	}

	defer r.finish(req, u, &res)

	for _, middleware := range r.Middlewares {
		middleware.Before(&req)
		defer middleware.After(req, &res)
	}
	res = r.route(req, u)
	return
}

func (h *Router) serveHTTP(w http.ResponseWriter, r *http.Request) (res Response) {
	// adapt the request from `net/http` format to `httpentity` format
	req := Request{
		Scheme:  r.URL.Scheme,
		Method:  r.Method,
		Headers: r.Header,
		Query:   r.URL.Query(),
		Entity:  nil,
		Things:  map[string]interface{}{},
	}
	u, mimetype := normalizeURL(r.URL)
	if mimetype != "" {
		// the file extension overrides the Accept: header
		req.Headers.Set("Accept", mimetype)
	}

	defer h.finish(req, u, &res)

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

	// run the request
	for _, middleware := range h.Middlewares {
		middleware.Before(&req)
		defer middleware.After(req, &res)
	}
	res = h.route(req, u)

	return
}

func (h *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	res := h.serveHTTP(w, req)
	// Adapt the response from `httpentity` format to `net/http` format
	for k, v := range res.Headers {
		w.Header().Set(k, strings.Join(v, ", "))
	}
	w.WriteHeader(int(res.status))
	res.WriteEntity(w)
}
