// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package maildir

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"syscall"
	"time"
)

var num_deliveries = big.NewInt(0)

func newUnique() Unique {
	now := time.Now()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	hostname = strings.Replace("/", "\\057", strings.Replace(":", "\\072", hostname, -1), -1)

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

type mailWriter struct {
	md     Maildir
	unique Unique
	file   *os.File
}

func (w *mailWriter) Close() error {
	err := w.file.Close()
	if err != nil {
		goto end
	}
	err = os.Link(string(w.md)+"/tmp/"+string(w.unique),
	              string(w.md)+"/new/"+string(w.unique))
end:
	syscall.Unlink(string(w.md) + "/tmp/" + string(w.unique))
	return err
}

func (w *mailWriter) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func (md Maildir) NewMail() io.WriteCloser {
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
