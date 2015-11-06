// Copyright 2015 Luke Shumaker

package httpentity

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (r *Router) Route(req Request, u *url.URL) (res Response) {
	if r.LogRequest {
		log.Printf("Route: %s %s %q\n", req.Scheme, req.Method, u.String())
	}
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
		Scheme:  "http",
		Method:  r.Method,
		Headers: r.Header,
		Query:   r.URL.Query(),
		Entity:  nil,
		Things:  map[string]interface{}{},
	}
	if r.TLS != nil {
		req.Scheme = "https"
	}
	if h.LogRequest {
		log.Printf("ServeHTTP: %s %s %q\n", req.Scheme, req.Method, r.URL.String())
	}
	if h.TrustForwarded {
		if scheme := req.Headers.Get("X-Forwarded-Proto"); scheme != "" {
			req.Scheme = scheme
		}
		if str := req.Headers.Get("Forwarded"); str != "" {
			parts := strings.Split(str, ";")
			for i := range parts {
				ary := strings.SplitN(parts[i], "=", 2)
				if len(ary) == 2 {
					if strings.EqualFold("proto", ary[0]) {
						req.Scheme = ary[1]
					}
				}
			}
		}
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
		resperr := req.readEntity(h, r.Body, r.Header.Get("Content-Type"))
		if resperr != nil {
			return *resperr
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
	w.WriteHeader(int(res.Status))
	res.writeEntity(w)
}
