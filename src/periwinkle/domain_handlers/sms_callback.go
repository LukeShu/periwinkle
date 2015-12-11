// Copyright 2015 Davis Webb
// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Luke Shumaker

package domain_handlers

import (
	"fmt"
	"io/ioutil"
	"locale"
	"net/http"
	"net/url"
	"periwinkle"
	"periwinkle/backend"
	"time"
)

type TwilioSMSCallbackServer struct {
	DB *periwinkle.DB
}

func (server TwilioSMSCallbackServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	periwinkle.Logf("TwilioCallback")
	fmt.Fprintf(w, "Hi there, I love %s!", req.URL.String())

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		periwinkle.LogErr(locale.UntranslatedError(err))
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		periwinkle.LogErr(locale.UntranslatedError(err))
		return
	}

	status := backend.TwilioSMSMessage{
		MessageStatus: values.Get("MessageStatus"),
		ErrorCode:     values.Get("ErrorCode"),
		MessageSID:    values.Get("MessageSid"),
	}
	server.DB.Do(func(db *periwinkle.Tx) {
		status.Save(db)
	})
}

func TwilioSMSWaitForCallback(conf *periwinkle.Cfg, messageSID string) (backend.TwilioSMSMessage, locale.Error) {
	var status backend.TwilioSMSMessage
	var err locale.Error
	done := false
	for !done {
		time.Sleep(time.Second)
		conf.DB.Do(func(db *periwinkle.Tx) {
			statusptr := backend.GetTwilioSMSMessageBySID(db, messageSID)
			if statusptr == nil {
				return
			}
			status = *statusptr
			switch status.MessageStatus {
			case "delivered":
				status.Delete(db)
				done = true
				return
			case "undelivered", "failed":
				status.Delete(db)
				err = locale.Errorf("TODO")
				done = true
				return
			}
		})
	}
	return status, err
}
