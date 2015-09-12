package inutil

import (
	"inotify"
	"os"
	"syscall"
)

const (
	// Flags for the parameter of InotifyInit1().
	// These, oddly, appear to be 24-bit numbers.
	IN_CLOEXEC = inotify.IN_CLOEXEC
)

type Watcher struct {
	Events <-chan inotify.Event
	events chan<- inotify.Event
	Errors <-chan error
	errors chan<- error
	in     *inotify.Inotify
}

func WatcherInit() (*Watcher, error) {
	in, err := inotify.InotifyInit()
	return newWatcher(in, err)
}

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

func (o *Watcher) AddWatch(path string, mask inotify.Mask) (inotify.Wd, error) {
	return o.in.AddWatch(path, mask)
}

func (o *Watcher) RmWatch(wd inotify.Wd) error {
	return o.in.RmWatch(wd)
}

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
