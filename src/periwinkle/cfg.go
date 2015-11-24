// Copyright 2015 Mark Pundman
// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package periwinkle

import (
	"io"
	"maildir"
	"net/http"
	"postfixpipe"

	"github.com/jinzhu/gorm"
)

type DomainHandler func(io.Reader, string, *gorm.DB, *Cfg) postfixpipe.ExitStatus

type Cfg struct {
	Mailstore            maildir.Maildir
	WebUIDir             http.Dir
	Debug                bool
	TrustForwarded       bool // whether to trust X-Forwarded: or Forwarded: HTTP headers
	TwilioAccountID      string
	TwilioAuthToken      string
	GroupDomain          string
	WebRoot              string
	DB                   *gorm.DB
	DomainHandlers       map[string]DomainHandler
	DefaultDomainHandler DomainHandler
}
