// Copyright 2015 Luke Shumaker

package main

import (
	"encoding/json"
	he "httpentity"
	"strings"
)

type postHack struct{}

func (p postHack) Before(req *he.Request) {
	if req.Method != "POST" {
		return
	}

	decoder, ok := req.Entity.(*json.Decoder)
	if !ok {
		return //putil.HTTPErrorf(415, "PUT and POST requests must have a document media type")
	}
	var entity interface{}
	err := decoder.Decode(&entity)
	if err != nil {
		return //putil.HTTPErrorf(415, "Couldn't parse: %v", err)
	}

	hash, ok := entity.(map[string]interface{})
	if !ok {
		return
	}

	method, ok := hash["_method"].(string)
	delete(hash, "_method")
	if ok {
		req.Method = method
	}

	xsrf_token, ok := hash["_xsrf_token"].(string)
	delete(hash, "_xsrf_token")
	if ok {
		req.Headers.Set("X-XSRF-TOKEN", xsrf_token)
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
}

func (p postHack) After(req he.Request, res *he.Response) {}
