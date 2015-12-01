// Copyright 2015 Luke Shumaker

// Package stoppable provides a wrapper for http.Server that can be
// gracefully stopped.
//
// It is a _much_ simpler package than "graceful" or "httpdown".
//
//     srv := stoppable.HTTPServer{Server: myserver, Socket: mylistener}
//     srv.Start() // does not block
//     srv.Stop() // does not block
//     err := srv.Wait() // blocks
package stoppable

import (
	"locale"
	"net"
	"net/http"
	"sync"
)

// HTTPServer provides an HTTP Server that can be gracefully stopped.
type HTTPServer struct {
	Server http.Server
	Socket net.Listener
	wg     sync.WaitGroup
	err    locale.Error
}

func (ss *HTTPServer) handleConnStateChange(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		ss.wg.Add(1)
	case http.StateHijacked, http.StateClosed:
		ss.wg.Done()
	}
}

// Start the server; does not block.
func (ss *HTTPServer) Start() {
	ss.Server.ConnState = ss.handleConnStateChange
	ss.wg.Add(1)
	go func() {
		defer ss.wg.Done()
		ss.err = locale.UntranslatedError(ss.Server.Serve(ss.Socket))
	}()
}

// Stop tells the server to stop; does not block.
func (ss *HTTPServer) Stop() {
	ss.Server.SetKeepAlivesEnabled(false)
	ss.Socket.Close()
}

// Wait for the server to stop; blocks.
func (ss *HTTPServer) Wait() locale.Error {
	ss.wg.Wait()
	return ss.err
}
