// Copyright 2015 Luke Shumaker <lukeshu@sbcglobal.net>.
//
// This is free software; you can redistribute it and/or modify it
// under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation; either version 2.1 of the
// License, or (at your option) any later version.
//
// This software is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this manual; if not, see
// <http://www.gnu.org/licenses/>.

// Package inotify provides an interface to the Linux inotify system.
// The inotify system is a mechanism for monitoring filesystem events.
package inotify

import (
	"sync"
	"syscall"
	"unsafe"
)

type Inotify struct {
	fd       file
	fdLock   sync.RWMutex
	buffFull [4096]byte
	buff     []byte
	buffLock sync.Mutex
}

type Event struct {
	Wd     Wd      // Watch descriptor
	Mask   Mask    // Mask describing event
	Cookie uint32  // Unique cookie associating related events (for rename(2))
	Name   *string // Optional name
}

// Create an inotify instance.  The variant InotifyInit1() allows
// flags to access extra functionality.
func InotifyInit() (*Inotify, error) {
	fd, err := inotify_init()
	o := Inotify{
		fd: fd,
	}
	o.buff = o.buffFull[:0]
	return &o, err
}

// Create an inotify instance, with flags specifying extra
// functionality.
func InotifyInit1(flags int) (*Inotify, error) {
	fd, err := inotify_init1(flags)
	o := Inotify{
		fd: fd,
	}
	o.buff = o.buffFull[:0]
	return &o, err
}

// Add a watch to the inotify instance, or modifies an existing watch
// item.
func (o *Inotify) AddWatch(path string, mask Mask) (Wd, error) {
	o.fdLock.RLock()
	defer o.fdLock.RUnlock()
	return inotify_add_watch(o.fd, path, mask)
}

// Remove a watch from the inotify instance.
func (o *Inotify) RmWatch(wd Wd) error {
	o.fdLock.RLock()
	defer o.fdLock.RUnlock()
	return inotify_rm_watch(o.fd, wd)
}

// Close the inotify instance; further calls to this object will
// error.
//
// Events recieved before Close() is called may still be Read() after
// the call to Close().
//
// Beware that if Close() is called while waiting on Read(), it will
// block until events are read.
func (o *Inotify) Close() error {
	o.fdLock.Lock()
	defer o.fdLock.Unlock()
	defer func() { o.fd = -1 }()
	return sysclose(o.fd)
}

// Read an event from the inotify instance.
//
// Events recieved before Close() is called may still be Read() after
// the call to Close().
func (o *Inotify) Read() (Event, error) {
	o.buffLock.Lock()
	defer o.buffLock.Unlock()

	if len(o.buff) == 0 {
		o.fdLock.RLock()
		len, err := sysread(o.fd, o.buffFull[:])
		o.fdLock.RUnlock()
		if len == 0 {
			return Event{Wd: -1}, o.Close()
		} else if len < 0 {
			return Event{Wd: -1}, err
		}
		o.buff = o.buffFull[0:len]
	}

	raw := (*syscall.InotifyEvent)(unsafe.Pointer(&o.buff[0]))
	ret := Event{
		Wd:     Wd(raw.Wd),
		Mask:   Mask(raw.Mask),
		Cookie: raw.Cookie,
		Name:   nil,
	}
	if raw.Len > 0 {
		bytes := (*[syscall.NAME_MAX]byte)(unsafe.Pointer(&o.buff[syscall.SizeofInotifyEvent]))
		name := string(bytes[:raw.Len-1])
		ret.Name = &name
	}
	o.buff = o.buff[0 : syscall.SizeofInotifyEvent+raw.Len]
	return ret, nil
}
