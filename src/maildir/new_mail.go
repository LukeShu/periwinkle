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

var numDeliveries = big.NewInt(0)

func newUnique() Unique {
	// http://cr.yp.to/proto/maildir.html
	now := time.Now()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	hostname = strings.Replace(strings.Replace(hostname, ":", "\\072", -1), "/", "\\057", -1)

	deliveryID := ""
	// deliveryID += fmt.Sprintf("#%x", unix_sequencenumber())
	// deliveryID += fmt.Sprintf("X%x", unix_bootnumber())
	randomData := make([]byte, 16)
	_, err = rand.Read(randomData)
	if err != nil {
		deliveryID += fmt.Sprintf("R%x", randomData)
	}
	// deliveryID += fmt.Sprintf("I%x", inode) // inodes are hard in go
	deliveryID += fmt.Sprintf("M%d", (now.UnixNano()/1000)-(now.Unix()*1000000))
	deliveryID += fmt.Sprintf("P%d", os.Getpid())
	deliveryID += fmt.Sprintf("Q%v", numDeliveries)
	numDeliveries.Add(big.NewInt(1), numDeliveries)

	return Unique(fmt.Sprintf("%v.%v.%v", now.Unix(), deliveryID, hostname))
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
		fmt.Println(err)
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
