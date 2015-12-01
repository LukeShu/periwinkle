// Copyright 2015 Luke Shumaker

package periwinkle

import (
	"fmt"
	"locale"
	"locale/gettext"
	"os"
	"strings"

	docopt "github.com/LukeShu/go-docopt"
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

func Docopt(usage string) map[string]interface{} {
	usage = strings.TrimSpace(fmt.Sprintf(usage, os.Args[0]))
	options, _ := docopt.Parse(usage, os.Args[1:], true, "", false, true)
	return options
}
