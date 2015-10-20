// Copyright 2015 Luke Shumaker

package httpentity

import (
	"httpentity/util"
	"net/http"
)

// Return the cookie `name`, or nil if it isn't set.
func (req *Request) Cookie(name string) *http.Cookie {
	if req.cookies == nil {
		req.cookies = heutil.ParseCookies(req.Headers)
	}
	cookie, ok := req.cookies[name]
	if ok {
		return cookie
	}
	return nil
}
