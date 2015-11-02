// Copyright 2013-2015 Docker, Inc.
// Copyright 2014 CoreOS, Inc.
// Copyright 2015 Luke Shumaker
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sd

import (
	"errors"
	"net"
	"os"
)

// errNotifyNoSocket is an error returned if no socket was specified.
var errNotifyNoSocket = errors.New("No socket")

// Notify sends a message to the service manager aobout state
// changes.  It is common to ignore the error.
//
// If unsetEnv is true, then (regarless of whether the function call
// itself succeeds or not) it will unset the environmental variable
// NOTIFY_SOCKET, which will cause further calls to this function to
// fail.
//
// The state parameter should countain a newline-separated list of
// variable assignments.
//
// See the documentation for sd_notify(3) for well-known variable
// assignments.
func Notify(unsetEnv bool, state string) error {
	if unsetEnv {
		defer os.Unsetenv("NOTIFY_SOCKET")
	}

	socketAddr := &net.UnixAddr{
		Name: os.Getenv("NOTIFY_SOCKET"),
		Net:  "unixgram",
	}

	if socketAddr.Name == "" {
		return errNotifyNoSocket
	}

	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(state))
	return err
}
