// Copyright 2015 Luke Shumaker

package heutil

import (
	"os"
	"io"
	"fmt"
)

type wLog struct {
	Writer io.Writer
}

func (l wLog) Printf(format string, v ...interface{}) {
	fmt.Fprintf(l.Writer, format, v...)
}

func (l wLog) Println(v ...interface{}) {
	fmt.Fprintln(l.Writer, v...)
}

var StderrLog = wLog{os.Stderr}
