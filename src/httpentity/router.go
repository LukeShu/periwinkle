// Copyright 2015 Luke Shumaker

package httpentity

import (
	"httpentity/negotiate"
	"mime"
	"net/url"
	"net/http"
	"path"
	"strings"
)

func normalizeURL(u1 *url.URL) (u *url.URL, mimetype string) {
	u, _ = u1.Parse("") // normalize
	// the file extension overrides the Accept: header
	if ext := path.Ext(u.Path); ext != "" {
		mimetype = mime.TypeByExtension(ext)
		u.Path = strings.TrimSuffix(u.Path, ext)
	}
	// add a trailing slash if there isn't one (so that relative
	// child URLs don't go to the parent)
	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}
	return
}

// Init initializes the hidden fields; should be called before any other
// method.
func (r Router) Init() *Router {
	r.initHandlers()
	return &r
}

// assumes that the url has already been passed to normalizeURL()
func (r *Router) finish(req Request, res *Response) {
	// generate a 500 error if we paniced
	if err := recover(); err != nil {
		*res = r.responseServerError(err)
	}
	// figure out the content type of the response
	if res.Entity != nil && res.Headers.Get("Content-Type") == "" {
		encoders := res.Entity.Encoders()
		mimetypes := encoders2mimetypes(encoders)
		acceptStrs, haveAccept := req.Headers[http.CanonicalHeaderKey("Accept")]
		var accept *string
		if haveAccept && len(acceptStrs) > 0 {
			accept = &acceptStrs[0]
		}
		options, err := negotiate.NegotiateContentType(accept, mimetypes)
		if err != nil {
			*res = r.responseBadRequest(err)
		} else {
			switch len(options) {
			case 0:
				*res = r.responseNotAcceptable(req.URL, mimetypes)
			case 1:
				//res.Headers.Set("Content-Type", options[0]+"; charset=utf-8")
				res.Headers.Set("Content-Type", options[0])
				return
			default:
				*res = r.responseMultipleChoices(req.URL, mimetypes)
			}
		}
		// If we make it here, we're either serving a
		// BadRequest (because parsing Accept failed),
		// NotAcceptable, or MultipleChoices, but don't know
		// what content type to send the message as!  We can't
		// just recurse because we don't want the status code
		// to change.
		encoders = res.Entity.Encoders()
		mimetypes = encoders2mimetypes(encoders)
		options, err = negotiate.NegotiateContentType(accept, mimetypes)
		switch len(options) {
		case 0:
			// Just pick one
			res.Headers.Set("Content-Type", mimetypes[0])
		case 1:
			res.Headers.Set("Content-Type", options[0])
		default:
			// Just pick one
			res.Headers.Set("Content-Type", options[0])
		}
	}
}
