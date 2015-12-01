// Copyright 2015 Luke Shumaker

package heutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"jsonpatch"
	"mime/multipart"
	"net/http"
	"locale"
	"net/url"
	"strings"
)

func fuckitJSON(entity interface{}) (interface{}, locale.Error) {
	str, uerr := json.Marshal(entity)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	return json.NewDecoder(strings.NewReader(string(str))), nil
}

// DecoderFormURLEncoded maps application/x-www-form-urlencoded => json.Decoder # because fuckit
func DecoderFormURLEncoded(r io.Reader, params map[string]string) (interface{}, locale.Error) {
	bytes, uerr := ioutil.ReadAll(r)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	entity, uerr := url.ParseQuery(string(bytes))
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	return fuckitJSON(entity)
}

// DecoderFormData maps multipart/form-data => json.Decoder # because fuckit
func DecoderFormData(r io.Reader, params map[string]string) (interface{}, locale.Error) {
	boundary, ok := params["boundary"]
	if !ok {
		return nil, locale.UntranslatedError(http.ErrMissingBoundary)
	}
	reader := multipart.NewReader(r, boundary)
	form, uerr := reader.ReadForm(0)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
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

// DecoderJSON maps application/json => json.Decoder
func DecoderJSON(r io.Reader, params map[string]string) (interface{}, locale.Error) {
	return json.NewDecoder(r), nil
}

// DecoderJSONPatch maps application/json-patch+json => jsonpatch.Patch
func DecoderJSONPatch(r io.Reader, params map[string]string) (interface{}, locale.Error) {
	bytes, uerr := ioutil.ReadAll(r)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	var patch jsonpatch.JSONPatch
	uerr = json.Unmarshal(bytes, &patch)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	return jsonpatch.Patch(patch), nil
}

// DecoderJSONMergePatch maps application/merge-patch+json => jsonpatch.Patch
func DecoderJSONMergePatch(r io.Reader, params map[string]string) (interface{}, locale.Error) {
	bytes, uerr := ioutil.ReadAll(r)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	var patch jsonpatch.JSONMergePatch
	uerr = json.Unmarshal(bytes, &patch)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	return jsonpatch.Patch(patch), nil
}

// DecoderOctetStream maps application/octet-stream => io.Reader
func DecoderOctetStream(r io.Reader, params map[string]string) (interface{}, locale.Error) {
	return r, nil
}
