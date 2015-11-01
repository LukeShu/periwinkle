// Copyright 2015 Luke Shumaker

package main

import (
	"periwinkle/listeners/web"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"io"
	sd "parabola.nshd/src/sd_daemon"
)

func usage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [ADDR]\n", os.Args[0])
}

func get_socket() (socket net.Listener, err error) {
	socket = nil
	err = nil
	fds := sd.ListenFds(true)
	if fds == nil {
		err = fmt.Errorf("Failed to aquire sockets from systemd")
		return
	}
	if len(fds) != 1 {
		err = fmt.Errorf("Wrong number of sockets from systemd: expected %d but got %d", 1, len(fds))
		return
	}
	socket, err = net.FileListener(fds[0])
	fds[0].Close()
	return
}

func main() {
	done := make(chan uint8)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGHUP)
	
	addr := ":http"
	
	switch len(os.Args)-1 {
	case 0:
		// do nothing
	case 1:
		addr = os.Args[1]
	default:
		usage(os.Stderr)
		os.Exit(1)
	}

	var socket net.Listener
	var err error
	if addr == "systemd" {
		socket, err = get_socket()
	} else {
		socket, err = listen(addr)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sd.Notify(false, "READY=1")

	go func() {
		err := web.Main(socket)
		fmt.Fprintln(os.Stderr, err)
		done <- 1
	}()

	for {
		select {
		case sig := <-signals:
			switch sig {
			case syscall.SIGTERM:
				sd.Notify(false, "STOPPING=1")
				web.Stop()
			case syscall.SIGHUP:
				sd.Notify(false, "RELOADING=1")
				// TODO: reload configuration file
				sd.Notify(false, "READY=1")
			}
		case status := <-done:
			os.Exit(int(status))
		}
	}
}
