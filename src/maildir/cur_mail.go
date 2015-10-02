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

type CurMail struct {
	md   Maildir
	uniq Unique
	info string
}

func (m *CurMail) path() string {
	return string(m.md) + "/cur/" + string(m.uniq) + ":" + m.info
}

func (md Maildir) Open(u Unique) (mail CurMail, err error) {
	matches, err := filepath.Glob(string(md) + "/cur/" + string(u) + ":*")
	if err != nil {
		return
	}
	if len(matches) != 1 {
		err = fmt.Errorf("Found %d files with unique %#v", u)
		return
	}
	parts := strings.SplitN(filepath.Base(matches[0]), ":", 2)
	mail = CurMail{md: md, uniq: Unique(parts[0]), info: parts[1]}
	return
}

func (md Maildir) Acknowledge(u Unique) (mail CurMail, err error) {
	err = os.Rename(string(md)+"/new/"+string(u),
	                string(md)+"/cur/"+string(u)+":")
	if err != nil {
		return
	}
	mail = CurMail{
		uniq: u,
		info: "",
	}
	return
}

func (m *CurMail) GetUnique() Unique {
	return m.uniq
}

func (m *CurMail) GetInfo() string {
	return m.info
}

func (m *CurMail) SetInfo(info string) error {
	err := os.Rename(m.path(), string(m.md)+"/cur/"+string(m.uniq)+":"+info)
	if err != nil {
		m.info = info
	}
	return err
}

func (m *CurMail) Delete() error {
	return syscall.Unlink(m.path())
}

type Reader interface {
	io.Reader
	io.Closer
	io.Seeker
}

func (m *CurMail) Reader() Reader {
	file, err := os.Open(m.path())
	if err != nil {
		return nil
	}
	return file
}
