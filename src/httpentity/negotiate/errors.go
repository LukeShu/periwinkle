// Copyright 2015 Luke Shumaker

package negotiate

import (
	"fmt"
)

// ParseError is an error that was encountered while parsing the
// value for Header.
type ParseError struct {
	Header  string
	Message error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("Parse Error: %s: %s", e.Header, e.Message)
}

func perrorf(header string, format string, args ...interface{}) error {
	return ParseError{
		Header:  header,
		Message: fmt.Errorf(format, args...),
	}
}
