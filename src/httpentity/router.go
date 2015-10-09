// Copyright 2015 Luke Shumaker

package httpentity

import (
	"bitbucket.org/ww/goautoneg"
	"fmt"
	"mime"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
)

type Router struct {
	Prefix      string
	Root        Entity
	Middlewares []Middleware
	Stacktrace  bool
}

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

// assumes that the url has already been passed to normalizeURL()
func (r *Router) route(req Request, u *url.URL) (res Response) {
	// do the routing
	res = route(r.Root, req, strings.TrimPrefix(u.Path, r.Prefix))

	// make sure the Location: header is absolute
	if l := res.Headers.Get("Location"); l != "" {
		u2, _ := u.Parse(l)
		res.Headers.Set("Location", u2.String())
		// XXX: this is pretty hacky, because it is tightly
		// integrated with the entity format used by
		// (Request).StatusCreated()
		if res.status == 201 {
			ilist := []interface{}(res.entity.(NetList))
			slist := make([]string, len(ilist))
			for i, iface := range ilist {
				slist[i] = iface.(string)
			}
			res.entity = extensions2net(u2, slist)
		}
	}

	return
}

// assumes that the url has already been passed to normalizeURL()
func (r *Router) finish(req Request, u *url.URL, res *Response) {
	// generate a 500 error if we paniced
	if err := recover(); err != nil {
		reason := err
		if r.Stacktrace {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			reason = fmt.Sprintf("%v\n\n%s", string(buf))
			fmt.Fprintf(os.Stderr, "%v\n\n%s", string(buf))
		}
		*res = req.statusInternalServerError(reason)
	}
	// figure out the content type of the response
	if res.entity != nil && res.Headers.Get("Content-Type") == "" {
		encoders := res.entity.Encoders()
		mimetypes := encoders2mimetypes(encoders)
		accept := req.Headers.Get("Accept")
		if len(encoders) > 1 && accept == "" {
			*res = req.statusMultipleChoices(u, mimetypes)
		} else {
			mimetype := goautoneg.Negotiate(req.Headers.Get("Accept"), mimetypes)
			if mimetype == "" {
				*res = req.statusNotAcceptable(u, mimetypes)
			} else {
				res.Headers.Set("Content-Type", mimetype+"; charset=utf-8")
			}
		}
	}
}
