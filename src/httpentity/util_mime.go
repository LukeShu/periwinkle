// Copyright 2015 Luke Shumaker

package httpentity

import (
	"httpentity/util"
	"io"
	"mime"
	"net/url"
	"strings"
)

func encoders2mimetypes(encoders map[string]func(out io.Writer) error) []string {
	list := make([]string, len(encoders))
	i := uint(0)
	for mimetype := range encoders {
		list[i] = mimetype
		i++
	}
	return list
}

func mimetypes2net(u *url.URL, mimetypes []string) NetEntity {
	u, _ = u.Parse("") // dup
	u.Path = strings.TrimSuffix(u.Path, "/")
	locations := make([]string, len(mimetypes))
	for i, mimetype := range mimetypes {
		u2, _ := u.Parse("")
		exts, _ := mime.ExtensionsByType(mimetype)
		u2.Path += exts[0]
		locations[i] = u2.String()
	}
	return heutil.NetList{locations}
}

func extensions2net(u *url.URL, extensions []string) NetEntity {
	u, _ = u.Parse("") // dup
	u.Path = strings.TrimSuffix(u.Path, "/")
	locations := make([]string, len(extensions))
	for i, extension := range extensions {
		u2, _ := u.Parse("")
		u2.Path += extension
		locations[i] = u2.String()
	}
	return heutil.NetList{locations}
}
