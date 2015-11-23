// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

// Package maildir implements access to DJB's maildir format.
//
// The type Maildir is a string of the path to the maildir.  The type
// Unique is a string of a unique maildir message identifier.  The
// type CurMail is a handle on a message that has been delivered.
package maildir

import (
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/djherbis/times"
)

// A Maildir is simply a string pathname to where the maildir is.
type Maildir string

// A Unique is a string that uniquely identifies a message in the
// maildir.  The format of the string is opaque.
//
//    "Unless you're writing messages to a maildir, the format of a
//    unique name is none of your business. A unique name can be
//    anything that doesn't contain a colon (or slash) and doesn't
//    start with a dot. Do not try to extract information from unique
//    names." -- http://cr.yp.to/proto/maildir.html
//
// Fortunatley for you, even if you are writing messages to a maildir,
// this package takes care of it, so the format of the unique name is
// still none of your business.
type Unique string

//    "It is a good idea for readers to skip all filenames in new and
//    cur starting with a dot.  Other than this, readers should not
//    attempt to parse filenames." --
//    http://www.qmail.org/qmail-manual-html/man5/maildir.html

// Remove files in `md+"/tmp/"` not accessed in the last 36 hours.
func (md Maildir) Clean() error {
	dir, err := os.Open(string(md) + "/tmp")
	if err != nil {
		return err
	}
	fileinfos, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fileinfo := range fileinfos {
		if strings.HasPrefix(fileinfo.Name(), ".") {
			continue
		}
		atime := times.Get(fileinfo).AccessTime()
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

// List the identifiers of newly delivered messages in the maildir.
func (md Maildir) ListNew() (mails []Unique, err error) {
	dir, err := os.Open(string(md) + "/new")
	if err != nil {
		return
	}
	fileinfos, err := dir.Readdir(0)
	if err != nil {
		return
	}
	mails = make([]Unique, 0, len(fileinfos))
	for _, fileinfo := range fileinfos {
		if strings.HasPrefix(fileinfo.Name(), ".") {
			continue
		}
		mails = append(mails, Unique(fileinfo.Name()))
	}
	return
}

// List old messages in the maildir.
func (md Maildir) ListCur() (mails []*CurMail, err error) {
	dir, err := os.Open(string(md) + "/cur")
	if err != nil {
		return
	}
	fileinfos, err := dir.Readdir(0)
	if err != nil {
		return
	}
	mails = make([]*CurMail, 0, len(fileinfos))
	for _, fileinfo := range fileinfos {
		if strings.HasPrefix(fileinfo.Name(), ".") {
			continue
		}
		parts := strings.SplitN(fileinfo.Name(), ":", 2)
		if len(parts) == 2 {
			mails = append(mails, &CurMail{md: md, uniq: Unique(parts[0]), info: parts[1]})
		}
	}
	return
}
