// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package maildir

import (
	"os"
	"strings"
	"syscall"
	"time"
)

type Maildir string

type Unique string

func (md Maildir) Clean() error {
	// TODO: remove files in `md+"/tmp/"` not accessed in 36 hours
	dir, err := os.Open(string(md) + "/tmp")
	if err != nil {
		return err
	}
	fileinfos, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fileinfo := range fileinfos {
		// TODO: check the access time (UNIX atime).
		// Unfortunately, Go os.FileInfo only provides the
		// modification time (UNIX mtime).  Actually, we
		// should probably take the newest of the
		// atime/mtime/ctime.
		atime := time.Now()
		if time.Since(atime) > (36 * time.Hour) {
			path := string(md) + "/tmp/" + fileinfo.Name()
			err1 := syscall.Unlink(path)
			if err1 != nil {
				err = &os.PathError{Op: "unlink", Path: path, Err: err1}
			}
		}
	}
	return err
}

func (md Maildir) ListNew() (mails []Unique, err error) {
	dir, err := os.Open(string(md) + "/new")
	if err != nil {
		return
	}
	fileinfos, err := dir.Readdir(0)
	if err != nil {
		return
	}
	mails = make([]Unique, len(fileinfos))
	for i, fileinfo := range fileinfos {
		mails[i] = Unique(fileinfo.Name())
	}
	return
}

func (md Maildir) ListCur() (mails []CurMail, err error) {
	dir, err := os.Open(string(md) + "/cur")
	if err != nil {
		return
	}
	fileinfos, err := dir.Readdir(0)
	if err != nil {
		return
	}
	mails = make([]CurMail, len(fileinfos))
	for i, fileinfo := range fileinfos {
		parts := strings.SplitN(fileinfo.Name(), ":", 2)
		if len(parts) == 2 {
			mails[i] = CurMail{md: md, uniq: Unique(parts[0]), info: parts[1]}
		}
	}
	return
}
