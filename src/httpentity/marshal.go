// Copyright 2015 Luke Shumaker

package httpentity

import (
	"httpentity/heutil"
	"io"
	"mime"
)

// If the Response has an entity, write it to the given output stream.
func (r Response) writeEntity(w io.Writer) error {
	if r.Entity == nil {
		return nil
	}
	mimetype, _, _ := mime.ParseMediaType(r.Headers.Get("Content-Type"))
	encoders := r.Entity.Encoders()
	return encoders[mimetype](w)
}

// Read an entity from the input stream, using the given content type.
func (req *Request) readEntity(router *Router, reader io.Reader, contenttype string) *Response {
	mimetype, params, err := mime.ParseMediaType(contenttype)
	if err != nil {
		res := router.responseBadRequest(heutil.NetPrintf("400 Bad Request: Could not parse Content-Type: %v", err))
		return &res
	}
	decoder, foundDecoder := router.Decoders[mimetype]
	if !foundDecoder {
		res := router.responseUnsupportedMediaType(heutil.NetString("415 Unsupported Media Type: Unsupported MIME type: " + mimetype))
		return &res
	}
	entity, err := decoder(reader, params)
	if err != nil {
		res := router.responseUnsupportedMediaType(heutil.NetPrintf("415 Unsupported Media Type: Error reading request body (%s): %v", mimetype, err))
		return &res
	}
	req.Entity = entity
	return nil
}
