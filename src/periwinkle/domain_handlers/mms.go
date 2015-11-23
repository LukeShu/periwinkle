// Copyright 2015 Davis Webb

package domain_handlers

import (
	"github.com/jinzhu/gorm"
	"io"
	"periwinkle"
	"postfixpipe"
)

func HandleMMS(r io.Reader, name string, db *gorm.DB, cfg *periwinkle.Cfg) postfixpipe.ExitStatus {
	panic("TODO")
}
