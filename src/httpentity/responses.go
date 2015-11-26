package httpentity

import (
	"fmt"
	"httpentity/heutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

func (r Router) responseMultipleChoices(u *url.URL, mimetypes []string) Response {
	return Response{
		Status:  300,
		Headers: http.Header{},
		Entity:  mimetypes2net(u, mimetypes),
	}
}

func (r Router) responseNotAcceptable(u *url.URL, mimetypes []string) Response {
	return Response{
		Status:  406,
		Headers: http.Header{},
		Entity:  mimetypes2net(u, mimetypes),
	}
}

func (r Router) responseBadRequest(e interface{}) Response {
	var ne NetEntity
	if e == nil {
		ne = heutil.NetString("400 Bad Request")
	} else if ne2, ok := e.(NetEntity); ok {
		ne = ne2
	} else {
		ne = heutil.NetString(fmt.Sprintf("400 Bad Request: %v", e))
	}
	return Response{
		Status:  400,
		Headers: http.Header{},
		Entity:  ne,
	}
}

func (r Router) responseUnsupportedMediaType(e NetEntity) Response {
	if e == nil {
		e = heutil.NetString("415 Unsupported Media Type")
	}
	return Response{
		Status:  415,
		Headers: http.Header{},
		Entity:  e,
	}
}

func (r Router) responseServerError(reason interface{}) Response {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	st := fmt.Sprintf("%T(%#v) => %v\n\n%s\n", reason, reason, reason, string(buf))
	for _, line := range strings.Split(st, "\n") {
		log.Println(line)
	}
	if r.Stacktrace {
		reason = st
	}
	return Response{
		Status: 500,
		Headers: http.Header{
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		Entity: heutil.NetPrintf("500 Internal Server Error: %v", reason),
	}
}
