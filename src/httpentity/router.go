// Copyright 2015 Luke Shumaker

package httpentity

import (
	"httpentity/negotiate"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type nilLog struct{}

func (l nilLog) Printf(format string, v ...interface{}) {}

func (l nilLog) Println(v ...interface{}) {}

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
	if r.Log == nil {
		r.Log = nilLog{}
	}
	r.initHandlers()
	return &r
}

func (router *Router) finish(req Request, res *Response) {
	// generate a 500 error if we paniced
	if err := recover(); err != nil {
		*res = router.responseServerError(err)
	}
	// negotiate all the things
	router.negotiate(req, res)
}

func (router *Router) negotiate(req Request, res *Response) {
	if res.Entity == nil {
		return
	}
	// Negotiate the content type
	{
		encoders := res.Entity.Encoders()
		contenttypes := encoders2contenttypes(encoders)

		var accept *string
		{
			acceptStrs, haveAccept := req.Headers[http.CanonicalHeaderKey("Accept")]
			if haveAccept && len(acceptStrs) > 0 {
				accept = &acceptStrs[0]
			}
		}

		options, err := negotiate.NegotiateContentType(accept, contenttypes)
		if err != nil {
			old := res.Status
			*res = router.responseBadRequest(err)
			router.Log.Printf("Could not parse Accept header; rewriting response code from %d to %d", old, res.Status)
			router.negotiate(req, res)
			return
		}
		switch len(options) {
		case 0:
			if !res.InhibitNotAcceptable {
				old := res.Status
				*res = router.responseNotAcceptable(req.URL, contenttypes)
				router.Log.Printf("No acceptable content type; rewriting response code from %d to %d", old, res.Status)
				router.negotiate(req, res)
				return
			}
			// Just pick one
			res.Headers.Set("Content-Type", contenttypes[0])
		case 1:
			res.Headers.Set("Content-Type", options[0])
		default:
			if !res.InhibitMultipleChoices {
				old := res.Status
				*res = router.responseMultipleChoices(req.URL, contenttypes)
				router.Log.Printf("Multiple choices; rewriting response code from %d to %d", old, res.Status)
				router.negotiate(req, res)
				return
			}
			// Just pick one
			res.Headers.Set("Content-Type", options[0])
		}
		res.encoder = encoders[res.Headers.Get("Content-Type")]
	}
	// Negotiate the charset (if applicable)
	if res.encoder.IsText() {
		// TODO
		res.Headers.Set("Content-Type", res.Headers.Get("Content-Type")+"; charset=utf-8")
	}
	// Negotiate the language
	{
		// TODO
		res.Headers.Set("Content-Language", "en-US")
	}
	// Negotiate the encoding
	{
		// TODO
	}
}
