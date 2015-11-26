// Copyright 2015 Luke Shumaker

package httpentity

import (
	"httpentity/negotiate"
	"mime"
	"net/url"
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

func (r Router) Init() *Router {
	r.initHandler()
	return &r
}

// assumes that the url has already been passed to normalizeURL()
func (r *Router) finish(req Request, u *url.URL, res *Response) {
	// generate a 500 error if we paniced
	if err := recover(); err != nil {
		*res = r.responseServerError(err)
	}
	// figure out the content type of the response
	if res.Entity != nil && res.Headers.Get("Content-Type") == "" {
		encoders := res.Entity.Encoders()
		mimetypes := encoders2mimetypes(encoders)
		accept := req.Headers.Get("Accept")
		if len(encoders) > 1 && accept == "" {
			*res = r.responseMultipleChoices(u, mimetypes)
		} else {
			options, err := negotiate.NegotiateContentType(&accept, mimetypes)
			if err != nil {
				*res = r.responseBadRequest(err)
			} else {
				switch len(options) {
				case 0:
					*res = r.responseNotAcceptable(u, mimetypes)
				case 1:
					//res.Headers.Set("Content-Type", mimetype+"; charset=utf-8")
					res.Headers.Set("Content-Type", mimetypes[0])
				default:
					*res = r.responseMultipleChoices(u, mimetypes)
				}
			}
		}
	}
}
