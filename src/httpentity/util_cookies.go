// Copyright 2015 Luke Shumaker

package httpentity

import (
	"net/http"
)

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

func (req *Request) Cookie(name string) *http.Cookie {
	if req.cookies == nil {
		stack := ParseCookies(req.Headers)
		req.cookies = &stack
	}
	cookie, ok := (*req.cookies)[name]
	if ok {
		return cookie
	}
	return nil
}
