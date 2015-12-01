// Copyright 2015 Luke Shumaker

package gettext

import (
	"locale"
	"os"
)

var _ locale.MessageCatalog = TextDomain{}

type Category int

func GetLocale(c Category) locale.Spec {
	for _, varname := range []string{"LC_ALL", c.String(), "LANG"} {
		if val := os.Getenv(varname); val != "" {
			return locale.Spec(val)
		}
	}
	return locale.Spec("C")
}

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

func (d TextDomain) Locales() []locale.Spec {
	panic("TODO")
}

func (d TextDomain) Translate(locale locale.Spec, str string) string {
	panic("TODO")
}

func (d TextDomain) TranslateN(locale locale.Spec, singular, plural string, n int) string {
	panic("TODO")
}

func (d TextDomain) TranslateP(locale locale.Spec, p, str string) string {
	return d.Translate(locale, p+d.ContextSep+str)
}

func (d TextDomain) TranslateNP(locale locale.Spec, p, singular, plural string, n int) string {
	return d.TranslateN(locale, p+d.ContextSep+singular, p+d.ContextSep+plural, n)
}
