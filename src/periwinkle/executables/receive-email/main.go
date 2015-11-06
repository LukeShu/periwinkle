// Copyright 2015 Luke Shumaker

package main

import (
	"log"
	"os"
	"periwinkle/cfg"
	"periwinkle/util" // putil
	"postfixpipe"
	"strings"
)

func main() {
	var ret uint8
	defer func() {
		if obj := recover(); obj != nil {
			log.Println(obj)
			if err, ok := obj.(error); ok {
				perror := putil.ErrorToError(err)
				ret = perror.PostfixCode()
			} else {
				ret = postfixpipe.EX_UNAVAILABLE
			}
		}
		os.Exit(int(ret))
	}()
	recipient := postfixpipe.OriginalRecipient()
	if recipient == "" {
		log.Println("ORIGINAL_RECIPIENT or RECIPIENT must be set")
		os.Exit(int(postfixpipe.EX_USAGE))
	}
	parts := strings.SplitN(recipient, "@", 2)
	user := parts[0]
	domain := "localhost"
	if len(parts) == 2 {
		domain = parts[1]
	}
	domain = strings.ToLower(domain)

	transaction := cfg.DB.Begin()
	defer func() {
		if err := transaction.Commit().Error; err != nil {
			panic(err)
		}
	}()

	handler, ok := cfg.DomainHandlers[domain]
	if ok {
		ret = handler(os.Stdin, user, transaction)
	} else {
		ret = cfg.DefaultDomainHandler(os.Stdin, recipient, transaction)
	}
}
