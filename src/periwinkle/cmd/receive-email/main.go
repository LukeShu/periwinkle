// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package main

import (
	"fmt"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/cmdutil"
	pp "postfixpipe"
	"runtime"
	"strings"

	"github.com/jinzhu/gorm"
)

const usage = `
Usage: %[1]s [-c CONFIG_FILE]
       %[1]s -h | --help
Install this in your Postfix ~/.forward or aliases file.

Options:
  -h, --help      Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`

func main() {
	options := cmdutil.Docopt(usage)
	config := cmdutil.GetConfig(options["-c"].(string))

	var ret pp.ExitStatus = pp.EX_OK
	defer func() {
		if reason := recover(); reason != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			st := fmt.Sprintf("%[1]T(%#[1]v) => %[1]v\n\n%[2]s", reason, string(buf))
			periwinkle.Logf("%s", st)
			ret = pp.EX_UNAVAILABLE
		}
		pp.Exit(ret)
	}()

	conflict := backend.WithTransaction(config.DB, func(transaction *gorm.DB) {
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
	})
	if conflict != nil {
		ret = pp.EX_DATAERR
	}
}
