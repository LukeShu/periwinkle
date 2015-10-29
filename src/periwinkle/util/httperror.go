// Copyright 2015 Luke Shumaker

package putil

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

type RawHTTPError he.Response

func (e RawHTTPError) Error() string {
	return fmt.Sprintf("%v", he.Response(e).Entity)
}

func (e RawHTTPError) Response() he.Response {
	return he.Response(e)
}

func HTTPErrorf(code int16, format string, a ...interface{}) HTTPError {
	// TODO: long-term: This printf is gross
	return RawHTTPError(he.Response{
		Status:  code,
		Headers: http.Header{},
		Entity:  heutil.NetString(fmt.Sprintf("%d: %s", code, fmt.Sprintf(format, a...))),
	})
}
