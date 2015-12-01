// Copyright 2015 Luke Shumaker

package rfc7231

import (
	he "httpentity"
	"strings"
)

func methods2string(methods map[string]func(request he.Request) he.Response) string {
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

// Middleware is a good default middleware for httpentity that
// provides much of the boilerplate for RFC 7231 semantics.
var Middleware = he.Middleware{
	Inside: func(request he.Request, entity he.Entity, handle func(he.Request, he.Entity) he.Response) (response he.Response) {
		switch request.Method {
		case "OPTIONS":
			response := handle(request, entity)
			if response.Status == 405 {
				methods := entity.Methods()
				response.Headers.Set("Allow", methods2string(methods))
			}
			response.Status = 204 // change to 200 when Entity is populated
			response.Entity = nil // TODO
			// TODO: this should give a lot more info
			// TODO: Accept-Patch, X-Accept-Put, X-Accept-Post headers
			// TODO: json-schema body
		case "HEAD":
			request.Method = "GET"
			response := handle(request, entity)
			response.Entity = nil
		default:
			response = handle(request, entity)
		}
		if response.Status == 405 {
			methods := entity.Methods()
			response.Headers.Set("Allow", methods2string(methods))
		}

		// Make sure the `Location:` header is absolute.  RFC
		// 7231 says they can be relative, but RFC 2616 said
		// they had to be absolute.  Plus, because of internal
		// URL rewriting, relative is possibly a bad idea.
		if l := response.Headers.Get("Location"); l != "" {
			u2, _ := request.URL.Parse(l)
			response.Headers.Set("Location", u2.String())
		}

		return
	},
}
