// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb
// Copyright 2014 Zhandos Suleimenov

package main

import (
	"os"
	"periwinkle"
	//"periwinkle/backend"
	"periwinkle/cmdutil"
	"periwinkle/test"

	"lukeshu.com/git/go/libsystemd.git/sd_daemon/lsb"
)

const usage = `
Usage: %[1]s [-c CONFIG_FILE]
       %[1]s -h | --help
Set up the RDBMS schema and seed data.

Options:
  -h, --help      Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`

func main() {
	options := cmdutil.Docopt(usage)
	config := cmdutil.GetConfig(options["-c"].(string))

	conflict := config.DB.Do(func(tx *periwinkle.Tx) {
		/*		
		err := backend.DbSchema(tx)
		if err != nil {
			periwinkle.Logf("Encountered an error while setting up the database schema, not attempting to seed data:")
			periwinkle.LogErr(err)
			os.Exit(int(lsb.EXIT_FAILURE))
		}

		err = backend.DbSeed(tx)
		if err != nil {
			periwinkle.Logf("Encountered an error while seeding the database:")
			periwinkle.LogErr(err)
			os.Exit(int(lsb.EXIT_FAILURE))
		}
		*/
		test.Test(config, tx)
	})
	if conflict != nil {
		periwinkle.LogErr(conflict)
		os.Exit(int(lsb.EXIT_FAILURE))
	}
}
