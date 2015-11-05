// Copyright 2015 Luke Shumaker

package postfixpipe

import (
	"os"
)

// Environment variables we get from Postfix (see `local(8postfix)`)
//
// - SHELL
// - HOME
// - USER
// - EXTENSION
// - DOMAIN
// - LOGNAME
// - LOCAL
// - ORIGINAL_RECIPIENT
// - RECIPIENT
// - SENDER
// - CLIENT_ADDRESS
// - CLIENT_HELO
// - CLIENT_HOSTNAME
// - CLIENT_PROTOCOL
// - SASL_METHOD
// - SASL_SENDER
// - SASL_USERNAME

func OriginalRecipient() string {
	recipient := os.Getenv("ORIGINAL_RECIPIENT")
	if recipient == "" {
		recipient = os.Getenv("RECIPIENT")
	}
	return recipient
}
