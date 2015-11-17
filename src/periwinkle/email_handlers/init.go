// Copyright 2015 Luke Shumaker

package handlers

import (
	"os"
	"io"
	"fmt"
	"periwinkle/cfg"
)

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s [CONFIG_FILE]\n", os.Args[0])
}

func init() {
	config_filename := "./periwinkle.conf"
	switch len(os.Args) {
	case 1:
		// do nothing
	case 2:
		config_filename = os.Args[1]
	default:
		usage(os.Stderr)
		os.Exit(2)
	}

	file, err := os.Open(config_filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %q: %v\n", config_filename, err)
		os.Exit(1)
	}

	config, err := cfg.Parse(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %q: %v\n", config_filename, err)
		os.Exit(1)
	}

	config.DomainHandlers = map[string]cfg.DomainHandler{
		"sms.gateway":   HandleSMS,
		"mms.gateway":   HandleMMS,
		config.GroupDomain: HandleEmail,
	}
}
