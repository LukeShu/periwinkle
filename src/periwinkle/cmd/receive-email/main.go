// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package main

import (
	"periwinkle"
	"periwinkle/cmdutil"
	"periwinkle/putil"
	pp "postfixpipe"
	"runtime"
	"strings"
)

const usage = `
Usage: %[1]s [-c CONFIG_FILE]
       %[1]s -h | --help
Install this in your Postfix ~/.forward or aliases file.

Options:
  -h, --help      Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`

func main() {
	options := periwinkle.Docopt(usage)
	config := cmdutil.GetConfig(options["-c"].(string))

	var ret pp.ExitStatus = pp.EX_OK
	defer func() {
		if obj := recover(); obj != nil {
			if err, ok := obj.(error); ok {
				perror := putil.ErrorToError(err)
				ret = perror.PostfixCode()
			} else {
				ret = pp.EX_UNAVAILABLE
			}
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			periwinkle.Logf("%T(%#v) => %v\n\n%s", obj, obj, obj, string(buf))
		}
		pp.Exit(ret)
	}()

	msg := pp.Get()

	recipient := msg.ORIGINAL_RECIPIENT()
	if recipient == "" {
		periwinkle.Logf("ORIGINAL_RECIPIENT must be set")
		ret = pp.EX_USAGE
		return
	}
	parts := strings.SplitN(recipient, "@", 2)
	user := parts[0]
	domain := "localhost"
	if len(parts) == 2 {
		domain = parts[1]
	}
	domain = strings.ToLower(domain)

	transaction := config.DB.Begin()
	defer func() {
		if err := transaction.Commit().Error; err != nil {
			panic(err)
		}
	}()

	reader, err := msg.Reader()
	if err != nil {
		periwinkle.LogErr(err)
		ret = pp.EX_NOINPUT
		return
	}

	if handler, ok := config.DomainHandlers[domain]; ok {
		ret = handler(reader, user, transaction, config)
	} else {
		ret = config.DefaultDomainHandler(reader, recipient, transaction, config)
	}
}
