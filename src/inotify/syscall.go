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

/* Create and initialize inotify instance.  */
func inotify_init() (Fd, error) {
	fd, errno := syscall.InotifyInit()
	return Fd(fd), os.NewSyscallError("inotify_init", errno)
}

/* Create and initialize inotify instance.  */
func inotify_init1(flags int) (Fd, error) {
	fd, errno := syscall.InotifyInit1(flags)
	return Fd(fd), os.NewSyscallError("inotify_init1", errno)
}

/* Add watch of object NAME to inotify instance FD.  Notify about
   events specified by MASK.  */
func inotify_add_watch(fd Fd, name string, mask Mask) (Wd, error) {
	wd, errno := syscall.InotifyAddWatch(int(fd), name, uint32(mask))
	return Wd(wd), newPathError("inotify_add_watch", name, errno)
}

/* Remove the watch specified by WD from the inotify instance FD.  */
func inotify_rm_watch(fd Fd, wd Wd) error {
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

func sysclose(fd Fd) error {
	return os.NewSyscallError("close", syscall.Close(int(fd)))
}

func sysread(fd Fd, p []byte) (int, error) {
	n, err := syscall.Read(int(fd), p)
	return n, os.NewSyscallError("read", err)
}
