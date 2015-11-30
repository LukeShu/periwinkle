// Copyright 2015 Luke Shumaker

package negotiate

import (
	"locale"
)

// ParseError is an error that was encountered while parsing the
// value for Header.
type ParseError struct {
	Header  string
	Message locale.Error
}

func (e ParseError) L10NString(l locale.Spec) string {
	return locale.Sprintf("Parse Error: %s: %s", e.Header, e.Message).L10NString(l)
}

func (e ParseError) Error() string {
	return e.L10NString("C")
}

func perrorf(header string, format string, args ...interface{}) locale.Error {
	return ParseError{
		Header:  header,
		Message: locale.Errorf(format, args...),
	}
}
