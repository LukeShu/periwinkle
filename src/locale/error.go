// Copyright 2015 Luke Shumaker

package locale

import (
	"fmt"
)

type untranslatedError struct {
	error
}

var _ Error = untranslatedError{}

var internalCatalog MessageCatalog = NullMessageCatalog{} // TODO

func UntranslatedError(in error) Error {
	if in == nil {
		return nil
	}
	if le, ok := in.(Error); ok {
		return le
	}
	return untranslatedError{in}
}

func (ue untranslatedError) L10NString(spec Spec) string {
	// TODO: in the future, this should try to recognize errors
	// from the Go standard library
	return fmt.Sprintf(internalCatalog.Translate(spec, "Untranslated error: %s"), ue.Error())
}
