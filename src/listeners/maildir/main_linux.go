// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package maildir

// +build ignore
/*
import (
	"inotify"
	"cfg"
)

func Main() {
	md := cfg.IncomingMail
	in := inotify.InotifyInit()
	in.AddWatch(string(md)+"/new", inotify.IN_ADD)
	for {
		select {
		case event := <-in.Event:
			handle(md)
		case err := <-in.Error:
			in.Close()
			return
		}
	}
}
*/
