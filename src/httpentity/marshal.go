// Copyright 2015 Luke Shumaker

package httpentity

import (
	"io"
	"locale"
	"mime"
)

// If the Response has an entity, write it to the given output stream.
func (response Response) writeEntity(w io.Writer) locale.Error {
	if response.encoder == nil {
		return nil
	}
	return response.encoder.Write(w, locale.Spec(response.Headers.Get("Content-Language")))
}

// Read an entity from the input stream, using the given content type.
func (req *Request) readEntity(router *Router, reader io.Reader, contenttype string) *Response {
	mimetype, params, uerr := mime.ParseMediaType(contenttype)
	if uerr != nil {
		res := router.responseBadRequest(locale.Errorf("Could not parse Content-Type: %v", locale.UntranslatedError(uerr)))
		return &res
	}
	decoder, foundDecoder := router.Decoders[mimetype]
	if !foundDecoder {
		res := router.responseUnsupportedMediaType(locale.Errorf("Unsupported MIME type: " + mimetype))
		return &res
	}
	entity, err := decoder(reader, params)
	if err != nil {
		res := router.responseUnsupportedMediaType(locale.Errorf("Error reading request body (%s): %v", mimetype, err))
		return &res
	}
	req.Entity = entity
	return nil
}
