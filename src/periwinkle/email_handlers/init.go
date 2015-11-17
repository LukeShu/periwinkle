// Copyright 2015 Luke Shumaker

package handlers

import (
	"periwinkle"
)

func GetHandlers(cfg *periwinkle.Cfg) error {
	cfg.DomainHandlers = map[string]periwinkle.DomainHandler{
		"sms.gateway":   HandleSMS,
		"mms.gateway":   HandleMMS,
		cfg.GroupDomain: HandleEmail,
	}
	return nil
}
