// Copyright 2015 Luke Shumaker

package heutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"jsonpatch"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

func fuckitJSON(entity interface{}) (interface{}, error) {
	str, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	return json.NewDecoder(strings.NewReader(string(str))), nil
}

// application/x-www-form-urlencoded => json.Decoder # because fuckit
func DecoderFormURLEncoded(r io.Reader, params map[string]string) (interface{}, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	entity, err := url.ParseQuery(string(bytes))
	if err != nil {
		return nil, err
	}
	return fuckitJSON(entity)
}

// multipart/form-data => json.Decoder # because fuckit
func DecoderFormData(r io.Reader, params map[string]string) (interface{}, error) {
	boundary, ok := params["boundary"]
	if !ok {
		return nil, http.ErrMissingBoundary
	}
	reader := multipart.NewReader(r, boundary)
	form, err := reader.ReadForm(0)
	if err != nil {
		return nil, err
	}
	entity := make(map[string]interface{}, len(form.Value)+len(form.File))
	for k, v := range form.Value {
		entity[k] = v
	}
	for k, v := range form.File {
		if _, exists := entity[k]; exists {
			values := entity[k].([]string)
			list := make([]interface{}, len(values)+len(v))
			i := uint(0)
			for _, value := range values {
				list[i] = value
				i++
			}
			for _, value := range v {
				list[i] = value
				i++
			}
			entity[k] = list
		} else {
			entity[k] = v
		}
	}
	return fuckitJSON(entity)
}

// application/json => json.Decoder
func DecoderJSON(r io.Reader, params map[string]string) (interface{}, error) {
	return json.NewDecoder(r), nil
}

// application/json-patch+json => jsonpatch.Patch
func DecoderJSONPatch(r io.Reader, params map[string]string) (interface{}, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var patch jsonpatch.JSONPatch
	err = json.Unmarshal(bytes, &patch)
	if err != nil {
		return nil, err
	}
	return jsonpatch.Patch(patch), err
}

// application/merge-patch+json => jsonpatch.Patch
func DecoderJSONMergePatch(r io.Reader, params map[string]string) (interface{}, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var patch jsonpatch.JSONMergePatch
	err = json.Unmarshal(bytes, &patch)
	if err != nil {
		return nil, err
	}
	return jsonpatch.Patch(patch), err
}

// application/octet-stream => io.Reader
func DecoderOctetStream(r io.Reader, params map[string]string) (interface{}, error) {
	return r, nil
}
