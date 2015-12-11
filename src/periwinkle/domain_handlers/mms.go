// Copyright 2015 Davis Webb

package domain_handlers

import (
	"io"
	"periwinkle"
	"postfixpipe"
)

func HandleMMS(r io.Reader, name string, db *periwinkle.Tx, cfg *periwinkle.Cfg) postfixpipe.ExitStatus {
	panic("TODO")
}
