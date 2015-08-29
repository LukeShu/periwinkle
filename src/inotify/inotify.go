package inotify

import (
	"errors"
	"syscall"
	"unsafe"
)

var InotifyAlreadyClosedError error = errors.New("inotify instance already closed")

type Inotify struct {
	fd       Cint
	isClosed bool

	fullbuff [4096]byte
	buff     []byte
}

type Event struct {
	Wd     Cint    /* Watch descriptor */
	Mask   Mask    /* Mask describing event */
	Cookie uint32  /* Unique cookie associating related events (for rename(2)) */
	Name   *string /* Optional name */
}

func InotifyInit() (*Inotify, error) {
	fd, err := inotify_init()
	o := Inotify{
		fd:       Cint(fd),
		isClosed: false,
	}
	o.buff = o.fullbuff[:]
	return &o, err
}

func InotifyInit1(flags Cint) (*Inotify, error) {
	fd, err := inotify_init1(flags)
	o := Inotify{
		fd:       Cint(fd),
		isClosed: false,
	}
	o.buff = o.fullbuff[:]
	return &o, err
}

func (o *Inotify) AddWatch(path string, mask Mask) (Cint, error) {
	if o.isClosed {
		return -1, InotifyAlreadyClosedError
	}
	return inotify_add_watch(o.fd, path, uint32(mask))
}

func (o *Inotify) RmWatch(wd Cint) error {
	if o.isClosed {
		return InotifyAlreadyClosedError
	}
	return inotify_rm_watch(o.fd, wd)
}

func (o *Inotify) Close() error {
	if o.isClosed {
		return InotifyAlreadyClosedError
	}
	o.isClosed = true
	return sysclose(o.fd)
}

func (o *Inotify) Read() (*Event, error) {
	if len(o.buff) == 0 {
		if o.isClosed {
			return nil, InotifyAlreadyClosedError
		}
		len, err := sysread(o.fd, o.buff)
		if len == 0 {
			return nil, o.Close()
		} else if len <= 0 {
			return nil, err
		}
		o.buff = o.fullbuff[0:len]
	}
	raw := (*syscall.InotifyEvent)(unsafe.Pointer(&o.buff[0]))
	var ret Event
	ret.Wd = Cint(raw.Wd)
	ret.Mask = Mask(raw.Mask)
	ret.Cookie = raw.Cookie
	ret.Name = nil
	if raw.Len > 0 {
		bytes := (*[syscall.NAME_MAX]byte)(unsafe.Pointer(&o.buff[syscall.SizeofInotifyEvent]))
		name := string(bytes[:raw.Len-1])
		ret.Name = &name
	}
	o.buff = o.buff[0 : syscall.SizeofInotifyEvent+raw.Len]
	return &ret, nil
}
