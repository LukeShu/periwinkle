// Copyright 2015 Luke Shumaker

// Package locale provides tools for internationalization (I18N) and
// localization (L10N).
package locale

// Spec is a locale specification; what to localize to the
// interpretation of it is backend-catalog dependent.
type Spec string

type MessageCatalog interface {
	Locales() []Spec
	Translate(locale Spec, str string) string
	TranslateN(locale Spec, singular, plural string, n int) string
	TranslateP(locale Spec, ctxt, str string) string
	TranslateNP(locale Spec, ctxt, singular, plural string, n int) string
}

// Stringer is a localizable fmt.Stringer.
type Stringer interface {
	Locales() []Spec
	L10NString(Spec) string
}

// Error is a localizable builtin.error
type Error interface {
	error
	Stringer
}

var DefaultCatalog MessageCatalog = NullMessageCatalog{}

func SetDefaultCatalog(c MessageCatalog) {
	DefaultCatalog = c
}
