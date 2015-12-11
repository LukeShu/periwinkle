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
	options := cmdutil.Docopt(usage)
	config := cmdutil.GetConfig(options["-c"].(string))

	conflict := config.DB.Do(func(tx *periwinkle.Tx) {
		err := backend.DbDrop(tx)
		if err != nil {
			periwinkle.LogErr(err)
			os.Exit(int(lsb.EXIT_FAILURE))
		}
	})
	if conflict != nil {
		periwinkle.LogErr(conflict)
		os.Exit(int(lsb.EXIT_FAILURE))
	}
}
