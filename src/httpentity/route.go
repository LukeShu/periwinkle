// Copyright 2015 Luke Shumaker

package httpentity

import (
	"bitbucket.org/ww/goautoneg"
	"mime"
	"net/url"
	"path"
	"strings"
)

func methods2string(methods map[string]Handler) string {
	set := make(map[string]bool, len(methods)+2)
	for method := range methods {
		set[method] = true
	}
	set["OPTIONS"] = true
	if _, get := set["GET"]; get {
		set["HEAD"] = true
	}
	list := make([]string, len(set))
	i := uint(0)
	for method := range set {
		list[i] = method
		i++
	}
	return strings.Join(list, ", ")
}

// Takes the normalized path without the leading slash
func route(entity Entity, req Request, method string, upath string) (ret Response) {
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
		parts := strings.SplitN(upath, "/", 2)
		if len(parts) != 2 {
			panic("path parser logic borked")
		}
		child := parts[0]
		grandchildren := parts[1]
		ret = route(entity.Subentity(child, req), req, method, grandchildren)
	}
	return
}

func Route(prefix string, entity Entity, req Request, method string, u *url.URL) (res Response) {
	// sanitize the URL
	u, _ = u.Parse("") // normalize
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
		if res.status == 201 {
			ilist := []interface{}(res.entity.(NetList))
			slist := make([]string, len(ilist))
			for i, iface := range ilist {
				slist[i] = iface.(string)
			}
			res.entity = extensions2net(u2, slist)
		}
	}
	// figure out the content type of the response
	if res.entity != nil && res.Headers.Get("Content-Type") == "" {
		encoders := res.entity.Encoders()
		mimetypes := encoders2mimetypes(encoders)
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
