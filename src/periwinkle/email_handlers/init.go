// Copyright 2015 Luke Shumaker

package handlers

import (
	"periwinkle/cfg"
)

func init() {
	cfg.DomainHandlers = map[string]cfg.DomainHandler{
		"sms.gateway":   HandleSMS,
		"mms.gateway":   HandleMMS,
		cfg.GroupDomain: HandleEmail,
	}
}
