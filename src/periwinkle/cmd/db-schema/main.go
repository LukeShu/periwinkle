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
Set up the RDBMS schema, but don't seed it.

Options:
  -h, --help      Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`

func main() {
	options := periwinkle.Docopt(usage)
	config := cmdutil.GetConfig(options["-c"].(string))

	err := backend.DbSchema(config.DB)
	if err != nil {
		periwinkle.Logf("Encountered an error while setting up the database schema:")
		periwinkle.LogErr(err)
		os.Exit(int(lsb.EXIT_FAILURE))
	}
}
