// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package cmdutil

import (
	"fmt"
	"locale"
	"os"
	"periwinkle"
	"periwinkle/cfg"
	"strings"

	docopt "github.com/LukeShu/go-docopt"
	"lukeshu.com/git/go/libsystemd.git/sd_daemon/lsb"
)

func GetConfig(filename string) *periwinkle.Cfg {
	configFile, uerr := os.Open(filename)
	if uerr != nil {
		periwinkle.Logf("Could not open config file: %v", locale.UntranslatedError(uerr))
		os.Exit(int(lsb.EXIT_NOTCONFIGURED))
	}

	config, err := cfg.Parse(configFile)
	if err != nil {
		periwinkle.Logf("Could not parse config file: %v", err)
		os.Exit(int(lsb.EXIT_NOTCONFIGURED))
	}

	return config
}

func Docopt(usage string) map[string]interface{} {
	usage = strings.TrimSpace(fmt.Sprintf(usage, os.Args[0]))
	options, _ := docopt.Parse(usage, os.Args[1:], true, "", false, true)
	return options
}
