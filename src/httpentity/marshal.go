// Copyright 2015 Luke Shumaker

package httpentity

import (
	"io"
	"mime"
	"fmt"
)

// If the Response has an entity, write it to the given output stream.
func (r Response) WriteEntity(w io.Writer) error {
	if r.Entity == nil {
		return nil
	}
	mimetype, _, _ := mime.ParseMediaType(r.Headers.Get("Content-Type"))
	encoders := r.Entity.Encoders()
	return encoders[mimetype](w)
}

// Read an entity from the input stream, using the given content type.
//
// TODO: how this works will probably change in the future to allow
// supporting other media types.
func (router *Router) ReadEntity(r io.Reader, contenttype string) (string, interface{}, error) {
	mimetype, params, err := mime.ParseMediaType(contenttype)
	if err != nil {
		return mimetype, nil, err
	}
	decoder, found_decoder := router.Decoders[mimetype]
	if !found_decoder {
		return mimetype, nil, fmt.Errorf("No decoder found: %s", mimetype)
	}
	entity, err := decoder(r, params)
	return mimetype, entity, err
}
