// Copyright 2015 Luke Shumaker

package httpentity

import (
	"bitbucket.org/ww/goautoneg"
	"io"
	"mime"
	"net/url"
	"path"
	"strings"
)

func methods2string(methods map[string]Handler) string {
	set := make(map[string]bool, len(methods)+2)
	for k, _ := range methods {
		set[k] = true
	}
	set["OPTIONS"] = true
	if _, get := set["GET"]; get {
		set["HEAD"] = true
	}
	list := make([]string, len(set))
	i := uint(0)
	for m, _ := range set {
		list[i] = m
		i++
	}
	return strings.Join(list, ", ")
}

// Takes the normalized path without the leading slash
func route(entity Entity, req Request, method string, upath string) Response {
	var ret Response
	if entity == nil {
		ret = req.statusNotFound()
	} else if upath == "" {
		callmethod := method
		if callmethod == "HEAD" {
			callmethod = "GET"
		}
		methods := entity.Methods()
		handler, method_allowed := methods[method]
		if method_allowed {
			ret = handler(req)
		} else {
			ret = req.statusMethodNotAllowed(methods2string(methods))
		}
		if callmethod == "OPTIONS" {
			ret.status = 200
			ret.Headers.Set("Allow", methods2string(methods))
		}
	} else {
		child := ""
		grandchildren := ""
		parts := strings.SplitN(upath, "/", 2)
		if len(parts) >= 1 {
			child = parts[0]
		}
		if len(parts) >= 2 {
			grandchildren = parts[1]
		}
		ret = route(entity.Subentity(child, req), req, method, grandchildren)
	}
	return ret
}

func encoders2mimelist(encoders map[string]Encoder) []string {
	list := make([]string, len(encoders))
	i := uint(0)
	for m, _ := range encoders {
		list[i] = m
		i++
	}
	return list
}

func Route(prefix string, entity Entity, req Request, method string, u *url.URL) Response {
	var res Response
	// just in case anything goes wrong, don't bring down the
	// process.
	defer func() {
		if r := recover(); r != nil {
			res = req.statusInternalServerError()
		}
	}()

	// sanitize the URL
	u, _ = url.Parse("") // normalize
	// the file extension overrides the Accept: header
	if ext := path.Ext(u.Path); ext != "" {
		req.Headers.Set("Accept", mime.TypeByExtension(ext))
		u.Path = strings.TrimSuffix(u.Path, ext)
	}
	// add a trailing slash if there isn't one (so that relative
	// child URLs don't go to the parent)
	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}

	// do the routing
	res = route(entity, req, method, strings.TrimPrefix(u.Path, prefix))

	// make sure the Location: header is absolute
	if l := res.Headers.Get("Location"); l != "" {
		u2, _ := u.Parse(l)
		res.Headers.Set("Location", u2.String())
	}
	// figure out the content type of the response
	if res.entity != nil && res.Headers.Get("Content-Type") == "" {
		encoders := res.entity.Encoders()
		mimetypes := encoders2mimelist(encoders)
		accept := req.Headers.Get("Accept")
		if len(encoders) > 1 && accept == "" {
			res = req.statusMultipleChoices(u, mimetypes)
		} else {
			mimetype := goautoneg.Negotiate(req.Headers.Get("Accept"), mimetypes)
			if mimetype == "" {
				res = req.statusNotAcceptable(u, mimetypes)
			} else {
				res.Headers.Set("Content-Type", mimetype+"; charset=utf-8")
			}
		}
	}

	// return the response
	return res
}

func (r Response) WriteEntity(w io.Writer) error {
	if r.entity == nil {
		return nil
	}
	mimetype := strings.SplitN(r.Headers.Get("Content-Type"), ";", 2)[0]
	return r.entity.Encoders()[mimetype](w)
}

func ReadEntity(w io.Reader, mimetype string, entity interface{}) {
	panic("not implemented")
}
