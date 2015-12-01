// Copyright 2015 Luke Shumaker

package locale

import (
	"fmt"
)

// This is mildly magic

type l10nStringer struct {
	catalog MessageCatalog
	format  string
	args    []interface{}
}

func Sprintf(format string, args ...interface{}) Stringer {
	return l10nStringer{
		catalog: DefaultCatalog,
		format:  format,
		args:    args,
	}
}

type l10nStringerError struct {
	data   Error
	locale Spec
}

func (s l10nStringerError) Error() string {
	return s.data.L10NString(s.locale)
}

type l10nStringerStringer struct {
	data   Stringer
	locale Spec
}

func (s l10nStringerStringer) String() string {
	return s.data.L10NString(s.locale)
}

func (s l10nStringer) L10NString(locale Spec) string {
	args := make([]interface{}, len(s.args))
	for i, arg := range s.args {
		switch v := arg.(type) {
		case Error:
			arg = error(l10nStringerError{v, locale})
		case Stringer:
			arg = fmt.Stringer(l10nStringerStringer{v, locale})
		}
		args[i] = arg
	}
	return fmt.Sprintf(s.catalog.Translate(locale, s.format), args...)
}

func (s l10nStringer) Locales() []Spec {
	// TODO: this needs to take the intersection of it and
	// localizable arguments
	return s.catalog.Locales()
}
