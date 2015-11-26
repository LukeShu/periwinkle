// Copyright 2015 Luke Shumaker

package httpentity

import (
	"log"
	"net/http"
	"strings"
)

func (r *Router) Route(req Request) (res Response) {
	if r.LogRequest {
		log.Printf("Route: %s %q\n", req.Method, req.URL.String())
	}
	u, mimetype := normalizeURL(req.URL)
	req.URL = u
	if mimetype != "" {
		// the file extension overrides the Accept: header
		req.Headers.Set("Accept", mimetype)
	}

	defer r.finish(req, &res)
	res = r.outsideHandler(req)
	return
}

func (h *Router) serveHTTP(w http.ResponseWriter, r *http.Request) (res Response) {
	// adapt the request from `net/http` format to `httpentity` format
	req := Request{
		Method:  r.Method,
		URL:     r.URL,
		Headers: r.Header,
		Entity:  nil,
		Things:  map[string]interface{}{},
	}
	if r.TLS != nil {
		req.URL.Scheme = "https"
	}
	if h.LogRequest {
		log.Printf("ServeHTTP: %s %q\n", req.Method, r.URL.String())
	}
	if h.TrustForwarded {
		if scheme := req.Headers.Get("X-Forwarded-Proto"); scheme != "" {
			req.URL.Scheme = scheme
		}
		if str := req.Headers.Get("Forwarded"); str != "" {
			parts := strings.Split(str, ";")
			for i := range parts {
				ary := strings.SplitN(parts[i], "=", 2)
				if len(ary) == 2 {
					if strings.EqualFold("proto", ary[0]) {
						req.URL.Scheme = ary[1]
					}
				}
			}
		}
	}
	u, mimetype := normalizeURL(req.URL)
	req.URL = u
	if mimetype != "" {
		// the file extension overrides the Accept: header
		req.Headers.Set("Accept", mimetype)
	}

	defer h.finish(req, &res)

	// parse the submitted entity
	switch req.Method {
	case "POST", "PUT", "PATCH":
		resperr := req.readEntity(h, r.Body, r.Header.Get("Content-Type"))
		if resperr != nil {
			return *resperr
		}
	}

	// run the request
	res = h.outsideHandler(req)

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
