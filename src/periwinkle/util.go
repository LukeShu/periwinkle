// Copyright 2015 Luke Shumaker

package periwinkle

import (
	"fmt"
	"locale"
	"locale/gettext"
	"os"
	"strings"
)

var serverMessageLocale = gettext.GetLocale(gettext.Messages)

func Logf(format string, a ...interface{}) {
	str := locale.Sprintf(format, a...).L10NString(serverMessageLocale)
	if strings.HasSuffix(format, "\n") {
		fmt.Fprint(os.Stderr, str)
	} else {
		fmt.Fprintln(os.Stderr, str)
	}
}

func LogErr(errs ...locale.Error) {
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err.L10NString(serverMessageLocale))
	}
}
