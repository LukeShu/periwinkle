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

package inotify

import (
	"os"
	"syscall"
)

func newPathError(op string, path string, err error) error {
	if err == nil {
		return nil
	}
	return &os.PathError{Op: op, Path: path, Err: err}
}

// Create and initialize inotify instance.
func inotify_init() (file, error) {
	fd, errno := syscall.InotifyInit()
	return file(fd), os.NewSyscallError("inotify_init", errno)
}

// Create and initialize inotify instance.
func inotify_init1(flags int) (file, error) {
	fd, errno := syscall.InotifyInit1(flags)
	return file(fd), os.NewSyscallError("inotify_init1", errno)
}

// Add watch of object NAME to inotify instance FD.  Notify about
// events specified by MASK.
func inotify_add_watch(fd file, name string, mask Mask) (Wd, error) {
	wd, errno := syscall.InotifyAddWatch(int(fd), name, uint32(mask))
	return Wd(wd), newPathError("inotify_add_watch", name, errno)
}

// Remove the watch specified by WD from the inotify instance FD.
func inotify_rm_watch(fd file, wd Wd) error {
	success, errno := syscall.InotifyRmWatch(int(fd), uint32(wd))
	switch success {
	case -1:
		if errno == nil {
			panic("should never happen")
		}
		os.NewSyscallError("inotify_rm_watch", errno)
	case 0:
		if errno != nil {
			panic("should never happen")
		}
		return nil
	}
	panic("should never happen")
}

func sysclose(fd file) error {
	return os.NewSyscallError("close", syscall.Close(int(fd)))
}

func sysread(fd file, p []byte) (int, error) {
	n, err := syscall.Read(int(fd), p)
	return n, os.NewSyscallError("read", err)
}
