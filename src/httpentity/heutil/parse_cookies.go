// Copyright 2015 Luke Shumaker

package heutil

import (
	"net/http"
)

// Parses the cookies from an HTTP request.
func ParseCookies(h http.Header) map[string]*http.Cookie {
	req := http.Request{
		Header: h,
	}
	ary := req.Cookies()
	cookies := map[string]*http.Cookie{}
	for _, cookie := range ary {
		cookies[cookie.Name] = cookie
	}
	return cookies
}
