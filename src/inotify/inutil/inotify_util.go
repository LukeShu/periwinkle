// Copyright 2015 Luke Shumaker <lukeshu@sbcglobal.net>.
//
// This is free software; you can redistribute it and/or modify it
// under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation; either version 2.1 of
// the License, or (at your option) any later version.
//
// This software is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this manual; if not, see
// <http://www.gnu.org/licenses/>.

// Package inutil provides a channel-based interface to inotify.
package inutil

import (
	"inotify"
	"os"
	"syscall"
)

type Watcher struct {
	Events <-chan inotify.Event
	Errors <-chan error
	events chan<- inotify.Event
	errors chan<- error
	in     *inotify.Inotify
}

// Wraps inotify.InotifyInit()
func WatcherInit() (*Watcher, error) {
	in, err := inotify.InotifyInit()
	return newWatcher(in, err)
}

// Wraps inotify.InotifyInit1()
func WatcherInit1(flags int) (*Watcher, error) {
	in, err := inotify.InotifyInit1(flags &^ inotify.IN_NONBLOCK)
	return newWatcher(in, err)
}

func newWatcher(in *inotify.Inotify, err error) (*Watcher, error) {
	events := make(chan inotify.Event)
	errors := make(chan error)
	o := &Watcher{
		Events: events,
		events: events,
		Errors: errors,
		errors: errors,
		in:     in,
	}
	go o.worker()
	return o, err
}

// Wraps inotify.Inotify.AddWatch(); adds or modifies a watch.
func (o *Watcher) AddWatch(path string, mask inotify.Mask) (inotify.Wd, error) {
	return o.in.AddWatch(path, mask)
}

// Wraps inotify.Inotify.RmWatch(); removes a watch.
func (o *Watcher) RmWatch(wd inotify.Wd) error {
	return o.in.RmWatch(wd)
}

// Wraps inotify.Inotify.Close().  Unlike inotify.Inotify.Close(),
// this cannot block.  Also unlike inotify.Inotify.Close(), nothing
// may be received from the channel after this is called.
func (o *Watcher) Close() {
	func() {
		defer recover()
		close(o.events)
		close(o.errors)
	}()
	go o.in.Close()
}

func (o *Watcher) worker() {
	defer recover()
	for {
		ev, err := o.in.Read()
		if ev.Wd >= 0 {
			o.events <- ev
		}
		if err != nil {
			if err.(*os.SyscallError).Err == syscall.EBADF {
				o.Close()
			}
			o.errors <- err
		}
	}
}
