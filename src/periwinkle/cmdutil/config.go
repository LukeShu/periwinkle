// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package cmdutil

import (
	"locale"
	"os"
	"periwinkle"
	"periwinkle/cfg"

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
