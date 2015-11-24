// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package httpapi

import (
	"encoding/json"
	he "httpentity"
	"httpentity/heutil"
	"io"
	"jsondiff"
)

func safeDecodeJSON(in interface{}, out interface{}) *he.Response {
	decoder, ok := in.(*json.Decoder)
	if !ok {
		ret := he.StatusUnsupportedMediaType(heutil.NetString("PUT and POST requests must have a document media type"))
		return &ret
	}
	var tmp interface{}
	err := decoder.Decode(&tmp)
	if err != nil {
		ret := he.StatusUnsupportedMediaType(heutil.NetPrintf("Couldn't parse: %v", err))
		return &ret
	}
	str, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, out)
	if err != nil {
		ret := he.StatusUnsupportedMediaType(heutil.NetPrintf("Request body didn't have expected structure (field had wrong data type): %v", err))
		return &ret
	}
	if !jsondiff.Equal(tmp, out) {
		diff, err := jsondiff.NewJSONPatch(tmp, out)
		if err != nil {
			panic(err)
		}
		entity := heutil.NetMap{
			"message": "Request body didn't have expected structure (extra or missing fields).  The included diff would make the request acceptable.",
			"diff":    diff,
		}
		ret := he.StatusUnsupportedMediaType(entity)
		return &ret
	}
	return nil
}

// Simple dump to JSON, good for most entities
func defaultEncoders(o interface{}) map[string]func(io.Writer) error {
	return map[string]func(io.Writer) error{
		"application/json": func(w io.Writer) error {
			bytes, err := json.Marshal(o)
			if err != nil {
				return err
			}
			_, err = w.Write(bytes)
			return err
		},
	}
}
