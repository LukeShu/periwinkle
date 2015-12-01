// Copyright 2015 Luke Shumaker

package heutil

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type wLog struct {
	Writer io.Writer
}

func (l wLog) Printf(format string, v ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Fprintf(l.Writer, format, v...)
}

func (l wLog) Println(v ...interface{}) {
	fmt.Fprintln(l.Writer, v...)
}

var StderrLog = wLog{os.Stderr}
