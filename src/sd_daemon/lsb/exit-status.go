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

// Package lsb provides constant exit codes specified by the Linux
// Standard Base.
package lsb

import (
	"os"
	"sd_daemon/logger"
)

// systemd daemon(7) recommends using the exit codes defined in the
// "LSB recomendations for SysV init scripts"[1].
//
// [1]: http://refspecs.linuxbase.org/LSB_3.1.1/LSB-Core-generic/LSB-Core-generic/iniscrptact.html
const (
	EXIT_SUCCESS         uint8 = 0
	EXIT_FAILURE         uint8 = 1
	EXIT_INVALIDARGUMENT uint8 = 2
	EXIT_NOTIMPLEMENTED  uint8 = 3
	EXIT_NOPERMISSION    uint8 = 4
	EXIT_NOTINSTALLED    uint8 = 5
	EXIT_NOTCONFIGURED   uint8 = 6
	EXIT_NOTRUNNING      uint8 = 7
	//   8- 99 are reserved for future LSB use
	// 100-149 are reserved for distribution use
	// 150-199 are reserved for application use
	// 200-254 are reserved for init system use

	// Therefore, the following are taken from systemd's
	// `src/basic/exit-status.h`
	EXIT_CHDIR               uint8 = 200
	EXIT_NICE                uint8 = 201
	EXIT_FDS                 uint8 = 202
	EXIT_EXEC                uint8 = 203
	EXIT_MEMORY              uint8 = 204
	EXIT_LIMITS              uint8 = 205
	EXIT_OOM_ADJUST          uint8 = 206
	EXIT_SIGNAL_MASK         uint8 = 207
	EXIT_STDIN               uint8 = 208
	EXIT_STDOUT              uint8 = 209
	EXIT_CHROOT              uint8 = 210
	EXIT_IOPRIO              uint8 = 211
	EXIT_TIMERSLACK          uint8 = 212
	EXIT_SECUREBITS          uint8 = 213
	EXIT_SETSCHEDULER        uint8 = 214
	EXIT_CPUAFFINITY         uint8 = 215
	EXIT_GROUP               uint8 = 216
	EXIT_USER                uint8 = 217
	EXIT_CAPABILITIES        uint8 = 218
	EXIT_CGROUP              uint8 = 219
	EXIT_SETSID              uint8 = 220
	EXIT_CONFIRM             uint8 = 221
	EXIT_STDERR              uint8 = 222
	_EXIT_RESERVED           uint8 = 223 // used to be tcpwrap don't reuse!
	EXIT_PAM                 uint8 = 224
	EXIT_NETWORK             uint8 = 225
	EXIT_NAMESPACE           uint8 = 226
	EXIT_NO_NEW_PRIVILEGES   uint8 = 227
	EXIT_SECCOMP             uint8 = 228
	EXIT_SELINUX_CONTEXT     uint8 = 229
	EXIT_PERSONALITY         uint8 = 230
	EXIT_APPARMOR_PROFILE    uint8 = 231
	EXIT_ADDRESS_FAMILIES    uint8 = 232
	EXIT_RUNTIME_DIRECTORY   uint8 = 233
	EXIT_MAKE_STARTER        uint8 = 234
	EXIT_CHOWN               uint8 = 235
	EXIT_BUS_ENDPOINT        uint8 = 236
	EXIT_SMACK_PROCESS_LABEL uint8 = 237
)

// This is a utility function to defer at the beginning of a goroutine
// in order to have the correct exit code in the case of a panic.
func Recover() {
	if r := recover(); r != nil {
		logger.Err("panic: %v", r)
		os.Exit(int(EXIT_FAILURE))
	}
}
