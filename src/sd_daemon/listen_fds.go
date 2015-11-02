// Copyright 2015 CoreOS, Inc.
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
	"os"
	"strconv"
	"syscall"
)

//#define SD_LISTEN_FDS_START 3
import "C"

// Returns a list of file descriptors passed in by the service manager
// as part of the socket-based activation logic.
//
// If unsetEnv is true, then (regarless of whether the function call
// itself succeeds or not) it will unset the environmental variables
// LISTEN_FDS and LISTEN_PID, which will cause further calls to this
// function to fail.
//
// In the case of an error, this function returns nil.
func ListenFds(unsetEnv bool) []*os.File {
	if unsetEnv {
		defer os.Unsetenv("LISTEN_PID")
		defer os.Unsetenv("LISTEN_FDS")
	}

	pid, err := strconv.Atoi(os.Getenv("LISTEN_PID"))
	if err != nil || pid != os.Getpid() {
		return nil
	}

	nfds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil || nfds == 0 {
		return nil
	}

	files := make([]*os.File, 0, nfds)
	for fd := C.SD_LISTEN_FDS_START; fd < C.SD_LISTEN_FDS_START+nfds; fd++ {
		syscall.CloseOnExec(fd)
		files = append(files, os.NewFile(uintptr(fd), "LISTEN_FD_"+strconv.Itoa(fd)))
	}

	return files
}
