// Copyright 2015 Luke Shumaker

package main

import (
	"fmt"
	"os"
	"periwinkle/cfg"
	"postfixpipe"
	"strings"
)

func main() {
	recipient := postfixpipe.OriginalRecipient()
	if recipient == "" {
		fmt.Fprintln(os.Stderr, "ORIGINAL_RECIPIENT or RECIPIENT must be set")
		os.Exit(postfixpipe.EX_USAGE)
	}
	parts := strings.SplitN(recipient, "@", 2)
	//user := parts[0]
	domain := "localhost"
	if len(parts) == 2 {
		domain = parts[1]
	}

	domain = strings.ToLower(domain)

	handler, ok := cfg.DomainHandlers[domain]
	if ok {
		os.Exit(handler())
	} else {
		os.Exit(cfg.DefaultDomainHandler())
	}
}
