// Copyright 2015 Luke Shumaker

package postfixpipe

import (
	"os"
)

// ExitStatus is an exit status that a program can pass to Postfix.
// May be a simple Sysexit code or (TODO) an RFC 3463 Enhanced Status
// Code.
type ExitStatus interface {
	ppexit()
}

// Exit and return a status to Postfix.
func Exit(es ExitStatus) {
	es.ppexit()
}

// simple exit(3) with codes from <sysexits.h>

// Sysexit is a simple numeric POSIX exit code.
type Sysexit uint8

func (n Sysexit) ppexit() {
	os.Exit(int(n))
}

// TODO: RFC 3463 Enhanced Status Codes

//func (es EnancedStatus) ppexit() {
//	// print stuff
//	os.Exit(1)
//}
