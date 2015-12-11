// Copyright 2015 Luke Shumaker

package domain_handlers

import (
	"periwinkle"
	"periwinkle/backend"
)

func GetHandlers(cfg *periwinkle.Cfg) error {
	cfg.DomainHandlers = map[string]periwinkle.DomainHandler{
		"sms.gateway":   HandleSMS,
		"mms.gateway":   HandleMMS,
		cfg.GroupDomain: HandleEmail,
	}
	return nil
}

func CanPost(db *periwinkle.Tx, group *backend.Group, userID string) bool {
	subscribed := backend.IsSubscribed(db, userID, *group)
        //  I'm assuming we wont have time to implement moderating
	//moderate := false
        if !backend.IsAdmin(db, userID, *group) {
                if group.PostPublic == 1 {
                        if subscribed == 0 {
                                return false
                        }
                        if group.PostConfirmed == 1 && subscribed == 1 {
                                return false
                        }
                        if group.PostMember == 1 {
                                return false
                        } /*  Probably not going to have time to implement moderating messages
                                if group.PostConfirmed == 2 && subscribed == 1 {
                                        moderate = true
                                } else if group.PostMember == 2 && subscribed == 2 {
                                        moderate = true
                                }
                        } else if group.PostPublic == 2 {
                                if subscribed == 0 {
                                        moderate = true
                                } else if group.PostConfirmed == 2 && subscribed == 1 {
                                        moderate = true
                                } else if group.PostMember == 2 && subscribed == 2 {
                                        moderate = true
                                }*/
                }
        }
	return true
}
