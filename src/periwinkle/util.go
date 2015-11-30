// Copyright 2015 Luke Shumaker

package periwinkle

import (
	"fmt"
	"locale"
	"locale/gettext"
	"os"
)

var serverMessageLocale = gettext.GetLocale(gettext.Messages)

func Logf(format string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, locale.Sprintf(format, a...).L10NString(serverMessageLocale))
}

func LogErr(errs ...locale.Error) {
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err.L10NString(serverMessageLocale))
	}
}
