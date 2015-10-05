// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

// +build !linux

package maildir

import (
	"periwinkle/cfg"
	"time"
)

func Main() error {
	md := cfg.IncomingMail
	for {
		time.Sleep(time.Second / 5)
		handle(md)
	}
}
