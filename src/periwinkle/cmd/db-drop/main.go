// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package main

import (
	"os"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/cmdutil"

	"lukeshu.com/git/go/libsystemd.git/sd_daemon/lsb"
)

const usage = `
Usage: %[1]s [-c CONFIG_FILE]
       %[1]s -h | --help
Drop the tables in the RDBMS, in the correct order.

Options:
  -h, --help      Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`

func main() {
	options := periwinkle.Docopt(usage)
	config := cmdutil.GetConfig(options["-c"].(string))

	err := backend.DbDrop(config.DB)
	if err != nil {
		periwinkle.LogErr(err)
		os.Exit(int(lsb.EXIT_FAILURE))
	}
}
