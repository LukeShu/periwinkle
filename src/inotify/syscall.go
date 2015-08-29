package inotify

import (
	"C"
	"os"
	"syscall"
)

type Cint C.int

/* Create and initialize inotify instance.  */
func inotify_init() (Cint, error) {
	fd, errno := syscall.InotifyInit()
	return Cint(fd), os.NewSyscallError("inotify_init", errno)
}

/* Create and initialize inotify instance.  */
func inotify_init1(flags Cint) (Cint, error) {
	fd, errno := syscall.InotifyInit1(int(flags))
	return Cint(fd), os.NewSyscallError("inotify_init1", errno)
}

/* Add watch of object NAME to inotify instance FD.  Notify about
   events specified by MASK.  */
func inotify_add_watch(fd Cint, name string, mask uint32) (Cint, error) {
	wd, errno := syscall.InotifyAddWatch(int(fd), name, mask)
	return Cint(wd), os.NewSyscallError("inotify_add_watch", errno)
}

/* Remove the watch specified by WD from the inotify instance FD.  */
func inotify_rm_watch(fd Cint, wd Cint) error {
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

func sysclose(fd Cint) error {
	return os.NewSyscallError("close", syscall.Close(int(fd)))
}

func sysread(fd Cint, p []byte) (Cint, error) {
	n, err := syscall.Read(int(fd), p)
	return Cint(n), os.NewSyscallError("read", err)
}
