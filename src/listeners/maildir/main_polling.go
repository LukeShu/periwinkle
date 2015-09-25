// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package maildir

// + b u i l d !linux

import (
	"time"
	"cfg"
)

func Main() {
	md := cfg.IncomingMail
	for {
		time.Sleep(time.Second/5)
		handle(md)
	}
}
