// Copyright 2015 Luke Shumaker

package gettext

import (
	"locale"
)

// Translate(N)?(P)?
//
//     N = Number (pluralizable)
//     P = Particular (context)
//
// GNU gettext uses "\004" as ContextSep, glib uses "|".
type TextDomain struct {
	// TODO: more stuff
	ContextSep string
}

// The "d" is for domain, as in libintl/gettext convention.

func (d MessageCatalog) Translate(locale locale.Spec, str string) string {
	panic("TODO")
}

func (d MessageCatalog) TranslateN(locale locale.Spec, singular, plural string, n int) string {
	panic("TODO")
}

func (d MessageCatalog) TranslateP(locale locale.Spec, p, str string) string {
	return d.Translate(locale, p+d.ContextSep+str)
}

func (d MessageCatalog) TranslateNP(locale locale.Spec, p, singular, plural string, n int) string {
	return d.TranslateN(locale, p+d.ContextSep+singular, p+d.ContextSep+plural, n)
}
