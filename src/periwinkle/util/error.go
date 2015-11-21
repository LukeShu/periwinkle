// Copyright 2015 Luke Shumaker

package putil

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3" // sqlite3
	"httpentity"
	"httpentity/util"
	"io"
	"net/http"
	"postfixpipe"
)

type Error interface {
	HttpCode() int16
	PostfixCode() postfixpipe.ExitStatus
	httpentity.NetEntity
	error
}

type simpleError struct {
	httpCode    int16
	postfixCode postfixpipe.ExitStatus
	httpStr     heutil.NetString
	plainStr    string
}

func (e simpleError) HttpCode() int16 {
	return e.httpCode
}

func (e simpleError) PostfixCode() postfixpipe.ExitStatus {
	return e.postfixCode
}

func (e simpleError) Encoders() map[string]func(io.Writer) error {
	return e.httpStr.Encoders()
}

func (e simpleError) Error() string {
	return string(e.plainStr)
}

func PErrorf(http int16, postfix postfixpipe.ExitStatus, format string, a ...interface{}) Error {
	str := fmt.Sprintf(format, a...)
	return simpleError{
		httpCode:    http,
		postfixCode: postfix,
		httpStr:     heutil.NetPrintf("%d: %s", http, str),
		plainStr:    str,
	}
}

func ErrorToError(err error) Error {
	switch err := err.(type) {
	case Error:
		return err
	case sqlite3.Error:
		if err.Code == sqlite3.ErrConstraint {
			return PErrorf(409, postfixpipe.EX_DATAERR, "%s", err.Error())
		}
	case *mysql.MySQLError:
		// TODO: this list of error numbers might not be complete
		// see https://mariadb.com/kb/en/mariadb/mariadb-error-codes/
		switch err.Number {
		case 1022, 1062, 1169, 1216, 1217, 1451, 1452, 1557, 1761, 1762, 1834:
			return PErrorf(409, postfixpipe.EX_DATAERR, "%s", err.Error())
		}
	}
	// we didn't intercept the error, so pass it along
	return PErrorf(500, postfixpipe.EX_UNAVAILABLE, "%s", err.Error())
}

func ErrorToHTTP(err Error) httpentity.Response {
	return httpentity.Response{
		Status:  err.HttpCode(),
		Headers: http.Header{},
		Entity:  err,
	}
}
