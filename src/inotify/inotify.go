package inotify

import (
	"sync"
	"syscall"
	"unsafe"
)

type Inotify struct {
	fd       Fd
	fdLock   sync.RWMutex
	buffFull [4096]byte
	buff     []byte
	buffLock sync.Mutex
}

type Event struct {
	Wd     Wd      /* Watch descriptor */
	Mask   Mask    /* Mask describing event */
	Cookie uint32  /* Unique cookie associating related events (for rename(2)) */
	Name   *string /* Optional name */
}

func InotifyInit() (*Inotify, error) {
	fd, err := inotify_init()
	o := Inotify{
		fd: fd,
	}
	o.buff = o.buffFull[:0]
	return &o, err
}

func InotifyInit1(flags int) (*Inotify, error) {
	fd, err := inotify_init1(flags)
	o := Inotify{
		fd: fd,
	}
	o.buff = o.buffFull[:0]
	return &o, err
}

func (o *Inotify) AddWatch(path string, mask Mask) (Wd, error) {
	o.fdLock.RLock()
	defer o.fdLock.RUnlock()
	return inotify_add_watch(o.fd, path, mask)
}

func (o *Inotify) RmWatch(wd Wd) error {
	o.fdLock.RLock()
	defer o.fdLock.RUnlock()
	return inotify_rm_watch(o.fd, wd)
}

func (o *Inotify) Close() error {
	o.fdLock.Lock()
	defer o.fdLock.Unlock()
	defer func() { o.fd = -1 }()
	return sysclose(o.fd)
}

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
