// Copyright 2015 Luke Shumaker

// Package null provides a null message catalog backend.
package null

import (
	"locale"
)

var _ locale.MessageCatalog = NullMessageCatalog{}

type NullMessageCatalog struct{}

func (NullMessageCatalog) Translate(locale locale.Spec, str string) string {
	return str
}

func (NullMessageCatalog) TranslateN(locale locale.Spec, singular, plural string, n int) string {
	if n == 1 {
		return singular
	} else {
		return plural
	}
}

func (NullMessageCatalog) TranslateP(locale locale.Spec, p, str string) string {
	return str
}

func (NullMessageCatalog) TranslateNP(locale locale.Spec, p, singular, plural string, n int) string {
	if n == 1 {
		return singular
	} else {
		return plural
	}
}
