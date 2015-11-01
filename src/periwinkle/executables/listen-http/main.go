// Copyright 2015 Luke Shumaker

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	sd "parabola.nshd/src/sd_daemon"
	"periwinkle/listeners/web"
	"strings"
	"syscall"
)

func usage(w io.Writer) {
	fmt.Fprintf(w,
		`Usage: %s [ADDR_TYPE] [ADDR]
Do the HTTP that you do, baby.

Address types are "tcp", "tcp4", "tcp6", "unix", and "systemd".  If
only one argument is given, if it matches one of type it is taken to
be the type; otherwise it is taken as an address.  The default address
for "tcp", "tcp4", and "tcp6" is ":8080"; the default "unix" address
is "/dev/stdin"; "systemd" doesn't have an address.  If the address is
given, the type is assumed to be "unix" if it contains a slash, or
"tcp" otherwise.
`, os.Args[0])
}

func parse_args() net.Listener {
	var stype, saddr string

	switch len(os.Args) - 1 {
	case 0:
		stype = "tcp"
		saddr = ":8080"
	case 1:
		switch os.Args[1] {
		case "tcp", "tcp4", "tcp6":
			stype = os.Args[1]
			saddr = ":8080"
		case "unix":
			stype = os.Args[1]
			saddr = "/dev/stdin"
		case "systemd":
			stype = os.Args[1]
		default:
			switch {
			case strings.ContainsRune(os.Args[1], '/'):
				stype = "unix"
			default:
				stype = "tcp"
			}
			saddr = os.Args[1]
		}
	case 2:
		stype = os.Args[1]
		saddr = os.Args[2]
	default:
		usage(os.Stderr)
		os.Exit(1)
	}

	var socket net.Listener
	var err error

	if saddr == "systemd" {
		socket, err = sd_get_socket()
	} else {
		socket, err = net.Listen(stype, saddr)
		if tcpsock, ok := socket.(*net.TCPListener); ok {
			socket = tcpKeepAliveListener{tcpsock}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return socket
}

func sd_get_socket() (socket net.Listener, err error) {
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

	socket := parse_args()

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
