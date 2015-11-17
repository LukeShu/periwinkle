// Copyright 2015 Luke Shumaker

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"periwinkle/cfg"
	_ "periwinkle/email_handlers" // handlers
	"periwinkle/util"             // putil
	"postfixpipe"
	"runtime"
	"strings"
)

func main() {
	var ret uint8
	defer func() {
		if obj := recover(); obj != nil {
			if err, ok := obj.(error); ok {
				perror := putil.ErrorToError(err)
				ret = perror.PostfixCode()
			} else {
				ret = postfixpipe.EX_UNAVAILABLE
			}
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			text := fmt.Sprintf("%T(%#v) => %v\n\n%s\n", obj, obj, obj, string(buf))
			for _, line := range strings.Split(text, "\n") {
				log.Println(line)
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

	reader := bufio.NewReader(os.Stdin)
	_, _, err := reader.ReadLine() // This is done to ignore the first line because it does not fit the RF822 format

	if err != nil {
		log.Println(err)
		ret = postfixpipe.EX_NOINPUT
		return
	}
	handler, ok := cfg.DomainHandlers[domain]
	if ok {
		ret = handler(reader, user, transaction)
	} else {
		ret = cfg.DefaultDomainHandler(reader, recipient, transaction)
	}
}
