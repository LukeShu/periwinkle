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

// Package logger implements a simple stderr-based logger with systemd
// log levels.
package logger

import (
	"fmt"
	"os"
)

///*#include <systemd/sd-daemon.h>*/
//#define SD_EMERG   "<0>"
//#define SD_ALERT   "<1>"
//#define SD_CRIT    "<2>"
//#define SD_ERR     "<3>"
//#define SD_WARNING "<4>"
//#define SD_NOTICE  "<5>"
//#define SD_INFO    "<6>"
//#define SD_DEBUG   "<7>"
import "C"

func log(level string, format string, a ...interface{}) {
	f := level + format + "\n"
	fmt.Fprintf(os.Stderr, f, a...)
}

// system is unusable
func Emerg(  format string, a ...interface{}) { log(C.SD_EMERG  , format, a...); }
// action must be taken immediately
func Alert(  format string, a ...interface{}) { log(C.SD_ALERT  , format, a...); }
// critical conditions
func Crit(   format string, a ...interface{}) { log(C.SD_CRIT   , format, a...); }
// error conditions
func Err(    format string, a ...interface{}) { log(C.SD_ERR    , format, a...); }
// warning conditions
func Warning(format string, a ...interface{}) { log(C.SD_WARNING, format, a...); }
// normal but significant condition
func Notice( format string, a ...interface{}) { log(C.SD_NOTICE , format, a...); }
// informational
func Info(   format string, a ...interface{}) { log(C.SD_INFO   , format, a...); }
// debug-level messages
func Debug(  format string, a ...interface{}) { log(C.SD_DEBUG  , format, a...); }
