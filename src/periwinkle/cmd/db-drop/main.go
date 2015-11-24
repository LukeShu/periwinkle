// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package main

import (
	"fmt"
	"io"
	"os"
	"periwinkle/backend"
	"periwinkle/cfg"
)

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s [CONFIG_FILE]\n", os.Args[0])
}

func main() {
	configFilename := "./periwinkle.yaml"
	switch len(os.Args) {
	case 1:
		// do nothing
	case 2:
		configFilename = os.Args[1]
	default:
		usage(os.Stderr)
		os.Exit(2)
	}

	file, err := os.Open(configFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %q: %v\n", configFilename, err)
		os.Exit(1)
	}

	config, err := cfg.Parse(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %q: %v\n", configFilename, err)
		os.Exit(1)
	}

	backend.DbDrop(config.DB)
}
