// Copyright 2015 Luke Shumaker

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	sd "lukeshu.com/git/go/libsystemd.git/sd_daemon"
	"periwinkle/listeners/web"
	"strconv"
	"strings"
	"syscall"
)

func usage(w io.Writer) {
	fmt.Fprintf(w,
		`Usage: %s [ADDR_TYPE] [ADDR]
Do the HTTP that you do, baby.

Address types are "tcp", "tcp4", "tcp6", "unix", and "fd".

If only one argument is given, if it matches one of type it is taken
to be the type; otherwise it is taken as an address.

  | type            | default address |
  |-----------------+-----------------|
  | tcp, tcp4, tcp6 | :8080           |
  | unix            | ./http.sock     |
  | fd              | stdin           |

If only the address is given, the type is assumed to be "unix" if it
contains a slash, "fd" if it only contains numeric digits or matches
one of the special "fd" values (below), or "tcp" otherwise.  If no
arguments are given, "tcp" is used.

The address for "fd" is numeric; however, there are several special
cases. "stdin", "stdout", and "stderr" are aliases for "0", "1", and
2", respectively. "systemd" causes it to look up the file descriptor
from systemd socket-activation.

If one argument is given, and it starts with a "-" or is "help", then
this message is displayed.
`, os.Args[0])
}

func parse_args() net.Listener {
	var stype, saddr string

	switch len(os.Args) - 1 {
	case 0:
		stype = "tcp"
		saddr = ":8080"
	case 1:
		if strings.HasPrefix(os.Args[1], "-") {
			usage(os.Stdout)
			os.Exit(0)
		}
		switch os.Args[1] {
		case "tcp", "tcp4", "tcp6":
			stype = os.Args[1]
			saddr = ":8080"
		case "unix":
			stype = os.Args[1]
			saddr = "./http.sock"
		case "fd":
			stype = os.Args[1]
			saddr = "stdin"
		case "systemd", "stdin", "stdout", "stderr":
			stype = "fd"
			saddr = os.Args[1]
		case "help":
			usage(os.Stdout)
			os.Exit(0)
		default:
			if strings.ContainsRune(os.Args[1], '/') {
				stype = "unix"
			} else if _, err := strconv.Atoi(os.Args[1]); err == nil {
				stype = "fd"
			} else {
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

	if stype == "fd" {
		switch saddr {
		case "systemd":
			socket, err = sd_get_socket()
		case "stdin":
			socket, err = listenfd(0, "/dev/stdin")
		case "stdout":
			socket, err = listenfd(1, "/dev/stdout")
		case "stderr":
			socket, err = listenfd(2, "/dev/stderr")
		default:
			var n int
			n, err = strconv.Atoi(saddr)
			if err == nil {
				socket, err = listenfd(n, "/dev/fd/"+saddr)
			}
		}
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

func listenfd(fd int, name string) (net.Listener, error) {
	return net.FileListener(os.NewFile(uintptr(fd), name))
}

func sd_get_socket() (socket net.Listener, err error) {
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

	web.Start(socket)
	go func() {
		err := web.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			done <- 1
		} else {
			done <- 0
		}
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
