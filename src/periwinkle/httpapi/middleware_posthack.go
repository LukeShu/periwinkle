// Copyright 2015 Luke Shumaker

package httpapi

import (
	"encoding/json"
	he "httpentity"
	"httpentity/heutil"
	"httpentity/rfc7231"
	"net/url"
	"strings"
)

func MiddlewarePostHack(req he.Request, u *url.URL, handle func(he.Request, *url.URL) he.Response) he.Response {
	if req.Method != "POST" {
		return handle(req, u)
	}

	decoder, ok := req.Entity.(*json.Decoder)
	if !ok {
		return handle(req, u)
	}
	var entity interface{}
	err := decoder.Decode(&entity)
	if err != nil {
		return rfc7231.StatusUnsupportedMediaType(heutil.NetPrintf("Couldn't parse: %v", err))
	}

	hash, ok := entity.(map[string]interface{})
	if !ok {
		return handle(req, u)
	}

	method, ok := hash["_method"].(string)
	delete(hash, "_method")
	if ok {
		req.Method = method
	}

	xsrfToken, ok := hash["_xsrf_token"].(string)
	delete(hash, "_xsrf_token")
	if ok {
		req.Headers.Set("X-XSRF-TOKEN", xsrfToken)
	}

	body, ok := hash["_body"]
	delete(hash, "_body")
	if ok {
		entity = body
	} else {
		entity = hash
	}

	str, err := json.Marshal(entity)
	if err != nil {
		panic(err)
	}
	req.Entity = json.NewDecoder(strings.NewReader(string(str)))
	return handle(req, u)
}
