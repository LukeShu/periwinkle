// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package maildir

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// A handle on an Acknowledge()ed (not newly-delivered) message in the
// maildir.
type CurMail struct {
	md   Maildir
	uniq Unique
	info string
}

func (m *CurMail) path() string {
	return string(m.md) + "/cur/" + string(m.uniq) + ":" + m.info
}

// Open the previously Acknowledg()ed message with the given unique
// identifier.
func (md Maildir) Open(u Unique) (mail *CurMail, err error) {
	matches, err := filepath.Glob(string(md) + "/cur/" + string(u) + ":*")
	if err != nil {
		return
	}
	if len(matches) != 1 {
		err = fmt.Errorf("Found %d files with unique %#v", len(matches), u)
		return
	}
	parts := strings.SplitN(filepath.Base(matches[0]), ":", 2)
	mail = &CurMail{md: md, uniq: Unique(parts[0]), info: parts[1]}
	return
}

// Acknowledge a newly delivered message (marking it as no longer
// newly delivered), and return a handle on it.
func (md Maildir) Acknowledge(u Unique) (mail *CurMail, err error) {
	err = os.Rename(
		string(md)+"/new/"+string(u),
		string(md)+"/cur/"+string(u)+":")
	if err != nil {
		return
	}
	mail = &CurMail{
		uniq: u,
		info: "",
	}
	return
}

// Return the unique idenfier for the message.
func (m *CurMail) GetUnique() Unique {
	return m.uniq
}

// Return the info string for the message, which "is morally
// equivalent to the Status field used by mbox readers."
//
// This package treats info as an opaque string, but obviously it
// would be good for it to have a common format between
// implementations. See <http://cr.yp.to/proto/maildir.html> for a
// recomendation for common semantics.
func (m *CurMail) GetInfo() string {
	return m.info
}

// Set the info string for the message, which "is morally equivalent
// to the Status field used by mbox readers."
//
// This package treats info as an opaque string, but obviously it
// would be good for it to have a common format between
// implementations. See <http://cr.yp.to/proto/maildir.html> for a
// recomendation for common semantics.
func (m *CurMail) SetInfo(info string) error {
	err := os.Rename(m.path(), string(m.md)+"/cur/"+string(m.uniq)+":"+info)
	if err != nil {
		m.info = info
	}
	return err
}

// Delete the message.
func (m *CurMail) Delete() error {
	return syscall.Unlink(m.path())
}

type Reader interface {
	io.Reader
	io.Closer
	io.Seeker
}

// Return an io.{Reader,Seeker,Closer} for the message, so that you
// can read its contents.
func (m *CurMail) Reader() Reader {
	file, err := os.Open(m.path())
	if err != nil {
		return nil
	}
	return file
}
