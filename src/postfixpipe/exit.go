// Copyright 2015 Luke Shumaker

package postfixpipe

import (
	"os"
)

type ExitStatus interface {
	ppexit()
}

func Exit(es ExitStatus) {
	es.ppexit()
}

// simple exit(3) with codes from <sysexits.h>

type Sysexit uint8

func (n Sysexit) ppexit() {
	os.Exit(int(n))
}

// TODO: RFC 3463 Enhanced Status Codes

//func (es EnancedStatus) ppexit() {
//	// print stuff
//	os.Exit(1)
//}
