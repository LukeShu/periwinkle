// Copyright 2015 Luke Shumaker

package locale

type l10nError l10nStringer

func Errorf(format string, args ...interface{}) Error {
	return l10nError{
		catalog: DefaultCatalog,
		format:  format,
		args:    args,
	}
}

func (s l10nError) L10NString(locale Spec) string {
	return l10nStringer(s).L10NString(locale)
}

func (s l10nError) Error() string {
	return s.L10NString("C")
}
