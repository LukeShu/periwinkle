// Copyright 2015 Davis Webb

package domain_handlers

import (
	"io"
	"periwinkle"
	"postfixpipe"

	"github.com/jinzhu/gorm"
)

func HandleMMS(r io.Reader, name string, db *gorm.DB, cfg *periwinkle.Cfg) postfixpipe.ExitStatus {
	panic("TODO")
}
