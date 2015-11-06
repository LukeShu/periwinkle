// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package maildir

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"syscall"
	"time"
)

var num_deliveries = big.NewInt(0)

func newUnique() Unique {
	// http://cr.yp.to/proto/maildir.html
	now := time.Now()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	hostname = strings.Replace(strings.Replace(hostname, ":", "\\072", -1), "/", "\\057", -1)

	delivery_id := ""
	// delivery_id += fmt.Sprintf("#%x", unix_sequencenumber())
	// delivery_id += fmt.Sprintf("X%x", unix_bootnumber())
	random_data := make([]byte, 16)
	_, err = rand.Read(random_data)
	if err != nil {
		delivery_id += fmt.Sprintf("R%x", random_data)
	}
	// delivery_id += fmt.Sprintf("I%x", inode) // inodes are hard in go
	delivery_id += fmt.Sprintf("M%d", (now.UnixNano()/1000)-(now.Unix()*1000000))
	delivery_id += fmt.Sprintf("P%d", os.Getpid())
	delivery_id += fmt.Sprintf("Q%v", num_deliveries)
	num_deliveries.Add(big.NewInt(1), num_deliveries)

	return Unique(fmt.Sprintf("%v.%v.%v", now.Unix(), delivery_id, hostname))
}

type Writer interface {
	Cancel() error
	Close() error
	Write([]byte) (int, error)
	Unique() Unique
}

type mailWriter struct {
	md     Maildir
	unique Unique
	file   *os.File
}

func (w *mailWriter) Cancel() (err error) {
	defer syscall.Unlink(string(w.md) + "/tmp/" + string(w.unique))
	err = w.file.Close()
	return
}

func (w *mailWriter) Close() (err error) {
	defer syscall.Unlink(string(w.md) + "/tmp/" + string(w.unique))
	err = w.file.Close()
	if err != nil {
		return
	}
	err = os.Link(
		string(w.md)+"/tmp/"+string(w.unique),
		string(w.md)+"/new/"+string(w.unique))
	return
}

func (w *mailWriter) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func (w *mailWriter) Unique() Unique {
	return w.unique
}

// Start the delivery of a new message to the maildir.  This function
// returns an io.WriteCloser; when .Close() is called on it, the
// message is delivered.
func (md Maildir) NewMail() Writer {
	unique := newUnique()
	file, err := os.OpenFile(string(md)+"/tmp/"+string(unique), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil
	}
	ret := &mailWriter{
		md:     md,
		unique: newUnique(),
		file:   file,
	}
	go func() {
		time.Sleep(24 * time.Hour)
		ret.Close()
	}()
	return ret
}
