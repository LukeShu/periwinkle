// Copyright 2015 Luke Shumaker

// Package null provides a null message catalog backend.
package locale

var _ MessageCatalog = NullMessageCatalog{}

type NullMessageCatalog struct{}

func (NullMessageCatalog) Locales() []Spec {
	return []Spec{"C", "en_US"}
}

func (NullMessageCatalog) Translate(locale Spec, str string) string {
	return str
}

func (NullMessageCatalog) TranslateN(locale Spec, singular, plural string, n int) string {
	if n == 1 {
		return singular
	} else {
		return plural
	}
}

func (NullMessageCatalog) TranslateP(locale Spec, p, str string) string {
	return str
}

func (NullMessageCatalog) TranslateNP(locale Spec, p, singular, plural string, n int) string {
	if n == 1 {
		return singular
	} else {
		return plural
	}
}
