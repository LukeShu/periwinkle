// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package main

import (
	"fmt"
	"os"
	"periwinkle/backend"
	"periwinkle/cfg"

	docopt "github.com/LukeShu/go-docopt"
	"lukeshu.com/git/go/libsystemd.git/sd_daemon/lsb"
)

var usage = fmt.Sprintf(`
Usage: %[1]s [-c CONFIG_FILE]
       %[1]s -h | --help
Set up the RDBMS schema and seed data.

Options:
  -h --help       Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`,
	os.Args[0])

func main() {
	options, _ := docopt.Parse(usage, os.Args[1:], true, "", false, true)

	configFile, err := os.Open(options["-c"].(string))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(int(lsb.EXIT_NOTCONFIGURED))
	}

	config, err := cfg.Parse(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(int(lsb.EXIT_NOTCONFIGURED))
	}

	backend.DbSchema(config.DB)
	backend.DbSeed(config.DB)
}
