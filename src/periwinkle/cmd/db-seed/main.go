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
	config_filename := "./periwinkle.yaml"
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

	backend.DbSchema(config.DB)
	backend.DbSeed(config.DB)
}
