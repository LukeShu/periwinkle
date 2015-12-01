package httpentity

import (
	"fmt"
	"locale"
	"net/http"
	"net/url"
	"runtime"
)

func (router Router) responseMultipleChoices(u *url.URL, mimetypes []string) Response {
	return Response{
		Status:                 300,
		Headers:                http.Header{},
		Entity:                 mimetypes2net(u, mimetypes),
		InhibitNotAcceptable:   true,
		InhibitMultipleChoices: true,
	}
}

func (router Router) responseNotAcceptable(u *url.URL, mimetypes []string) Response {
	return Response{
		Status:                 406,
		Headers:                http.Header{},
		Entity:                 mimetypes2net(u, mimetypes),
		InhibitNotAcceptable:   true,
		InhibitMultipleChoices: true,
	}
}

func (router Router) responseBadRequest(err locale.Error) Response {
	return Response{
		Status:                 400,
		Headers:                http.Header{},
		Entity:                 ErrorToNetEntity(400, err),
		InhibitNotAcceptable:   true,
		InhibitMultipleChoices: true,
	}
}

func (router Router) responseUnsupportedMediaType(err locale.Error) Response {
	return Response{
		Status:  415,
		Headers: http.Header{},
		Entity:  ErrorToNetEntity(415, err),
	}
}

func (router Router) responseServerError(reason interface{}) Response {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	st := fmt.Sprintf("%[1]T(%#[1]v) => %[1]v\n\n%[2]s", reason, string(buf))
	router.Log.Println(st)
	if router.Stacktrace {
		reason = st
	}
	return Response{
		Status:                 500,
		Headers:                http.Header{},
		Entity:                 NetPrintf("500 Internal Server Error: %v", reason),
		InhibitNotAcceptable:   true,
		InhibitMultipleChoices: true,
	}
}
