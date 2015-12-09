// Copyright 2015 Luke Shumaker

package rfc7231

import (
	he "httpentity"
	"mime"
	"net/url"
	"strings"
)

func encoders2contenttypes(encoders map[string]he.Encoder) []string {
	list := make([]string, len(encoders))
	i := uint(0)
	for mimetype := range encoders {
		list[i] = mimetype
		i++
	}
	return list
}

func mimetypes2net(u *url.URL, mimetypes []string) he.NetEntity {
	u, _ = u.Parse("") // dup
	u.Path = strings.TrimSuffix(u.Path, "/")
	locations := make([]*url.URL, len(mimetypes))
	for i, mimetype := range mimetypes {
		u2, _ := u.Parse("") // dup
		exts, _ := mime.ExtensionsByType(mimetype)
		if exts == nil || len(exts) == 0 {
			u2.Path += "rfc7231_mimetypes2net_no_extension_should_never_happen?" + mimetype
		} else {
			u2.Path += exts[0]
		}
		locations[i] = u2
	}
	return he.NetLocations(locations)
}
