// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

// +build linux

package maildir

import (
	"inotify"
	"inotify/inutil"
	"cfg"
)

func Main() error {
	md := cfg.IncomingMail
	in, err := inutil.WatcherInit()
	if err != nil {
		return err
	}
	defer in.Close();
	in.AddWatch(string(md)+"/new", inotify.IN_CREATE | inotify.IN_MOVED_TO)
	for {
		select {
		case _ = <-in.Events:
			handle(md)
		case err := <-in.Errors:
			return err
		}
	}
}
