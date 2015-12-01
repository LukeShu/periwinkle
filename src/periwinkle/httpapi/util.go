// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package httpapi

import (
	"encoding/json"
	he "httpentity"
	"httpentity/rfc7231"
	"io"
	"jsondiff"
	"jsonpatch"
	"locale"
)

type decodeJSONError struct {
	message locale.Stringer
	diff    jsonpatch.JSONPatch
}

// Encoders fulfills the httpentity.NetEntity interface.
func (l decodeJSONError) Encoders() map[string]he.Encoder {
	return map[string]he.Encoder{
		"application/json": l,
	}
}

func (l decodeJSONError) Locales() []locale.Spec {
	return l.message.Locales()
}

func (l decodeJSONError) IsText() bool {
	return true
}

func (l decodeJSONError) Write(w io.Writer, loc locale.Spec) locale.Error {
	data := map[string]interface{}{
		"message": l.message.L10NString(loc),
		"diff":    l.diff,
	}
	bytes, uerr := json.Marshal(data)
	if uerr != nil {
		return locale.UntranslatedError(uerr)
	}
	_, uerr = w.Write(bytes)
	return locale.UntranslatedError(uerr)
}

func safeDecodeJSON(in interface{}, out interface{}) *he.Response {
	decoder, ok := in.(*json.Decoder)
	if !ok {
		ret := rfc7231.StatusUnsupportedMediaType(he.NetPrintf("PUT and POST requests must have a document media type"))
		return &ret
	}
	var tmp interface{}
	err := decoder.Decode(&tmp)
	if err != nil {
		ret := rfc7231.StatusUnsupportedMediaType(he.NetPrintf("Couldn't parse: %v", err))
		return &ret
	}
	str, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, out)
	if err != nil {
		ret := rfc7231.StatusUnsupportedMediaType(he.NetPrintf("Request body didn't have expected structure (field had wrong data type): %v", err))
		return &ret
	}
	if !jsondiff.Equal(tmp, out) {
		diff, err := jsondiff.NewJSONPatch(tmp, out)
		if err != nil {
			panic(err)
		}
		entity := decodeJSONError{
			message: locale.Sprintf("Request body didn't have expected structure (extra or missing fields).  The included diff would make the request acceptable."),
			diff:    diff,
		}
		ret := rfc7231.StatusUnsupportedMediaType(entity)
		return &ret
	}
	return nil
}

// Simple dump to JSON, good for most entities
func defaultEncoders(o interface{}) map[string]he.Encoder {
	return map[string]he.Encoder{
		"application/json": he.EncoderJSON{o},
	}
}
