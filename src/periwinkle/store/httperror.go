// Copyright 2015 Luke Shumaker

package store

import (
	"fmt"
	he "httpentity"
	"httpentity/util"
	"net/http"
)

type HTTPError interface {
	error
	Response() he.Response
}

type httpError he.Response

func (e httpError) Error() string {
	return fmt.Sprintf("%v", he.Response(e).Entity)
}

func (e httpError) Response() he.Response {
	return he.Response(e)
}

func httpErrorf(code int16, format string, a ...interface{}) HTTPError {
	// TODO: long-term: This printf is gross
	return httpError(he.Response{
		Status:  code,
		Headers: http.Header{},
		Entity:  heutil.NetString(fmt.Sprintf("%d: %s", code, fmt.Sprintf(format, a...))),
	})
}
