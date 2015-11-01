// Copyright 2015 Luke Shumaker

// Package stoppable provides a wrapper for http.Server that can be
// gracefully stopped.
//
// It is a _much_ simpler package than "graceful" or "httpdown".
//
// From the main goroutine:
//
//     srv := stoppable.HTTPServer{Server: myserver, Socket: mylistener}
//     err := srv.Serve() // blocks until done
//
// From another goroutine:
//
//     srv.Stop() // does not block
package stoppable

import (
	"net"
	"net/http"
	"sync"
)

type HTTPServer struct {
	Server http.Server
	Socket net.Listener
	wg     sync.WaitGroup
}

func (ss *HTTPServer) handleConnStateChange(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		ss.wg.Add(1)
	case http.StateHijacked, http.StateClosed:
		ss.wg.Done()
	}
}

// Blocks until the server is done.
func (ss *HTTPServer) Serve() (err error) {
	ss.Server.ConnState = ss.handleConnStateChange
	ss.wg.Add(1)
	go func() {
		defer ss.wg.Done()
		err = ss.Server.Serve(ss.Socket)
	}()
	ss.wg.Wait()
	return
}

// Does not block.
func (ss *HTTPServer) Stop() {
	ss.Server.SetKeepAlivesEnabled(false)
	ss.Socket.Close()
}
