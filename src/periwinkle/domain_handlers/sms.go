// Copyright 2015 Davis Webb
// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Luke Shumaker

package domain_handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/twilio"
	"postfixpipe"
	"strings"

	"github.com/jinzhu/gorm"
)

func HandleSMS(r io.Reader, name string, db *gorm.DB, cfg *periwinkle.Cfg) postfixpipe.ExitStatus {
	message, err := mail.ReadMessage(r)
	if err != nil {
		log.Println(err)
		return postfixpipe.EX_NOINPUT
	}
	status, err := sender(*message, name, db, cfg)
	if err != nil {
		log.Println(err)
		return postfixpipe.EX_NOINPUT
	}
	log.Println(status)
	return postfixpipe.EX_OK
}

// Returns the status of the message: queued, sending, sent,
// delivered, undelivered, failed.  If an error occurs, it returns
// Error.
func sender(message mail.Message, smsTo string, db *gorm.DB, cfg *periwinkle.Cfg) (status string, err error) {

	group := message.Header.Get("From")
	user := backend.GetUserByAddress(db, "sms", smsTo)

	smsFrom := backend.GetTwilioNumberByUserAndGroup(db, user.ID, strings.Split(group, "@")[0])
	smsBody := message.Header.Get("Subject")
	//smsBody, err := ioutil.ReadAll(message.Body)
	//if err != nil {
	//	return "", err
	//}

	messagesURL := "https://api.twilio.com/2010-04-01/Accounts/" + cfg.TwilioAccountID + "/Messages.json"

	v := url.Values{}
	v.Set("From", smsFrom)
	v.Set("To", smsTo)
	v.Set("Body", string(smsBody))
	v.Set("StatusCallback", cfg.WebRoot+"/callbacks/twilio-sms")

	client := &http.Client{}

	req, err := http.NewRequest("POST", messagesURL, bytes.NewBuffer([]byte(v.Encode())))
	if err != nil {
		return
	}
	req.SetBasicAuth(cfg.TwilioAccountID, cfg.TwilioAuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		message := twilio.Message{}
		json.Unmarshal([]byte(body), &message)
		smsStatus, err := SmsWaitForCallback(message.Sid)

		if err != nil {
			return "", err
		}

		if smsStatus.MessageStatus == "undelivered" || smsStatus.MessageStatus == "failed" {
			return smsStatus.MessageStatus, fmt.Errorf("%s", smsStatus.ErrorCode)
		}
		if smsStatus.MessageStatus == "queued" || smsStatus.MessageStatus == "sending" || smsStatus.MessageStatus == "sent" {
			smsStatus, err = SmsWaitForCallback(message.Sid)

			if err != nil {
				return "", err
			}
		}

		if smsStatus.MessageStatus == "undelivered" || smsStatus.MessageStatus == "failed" {
			return smsStatus.MessageStatus, fmt.Errorf("%s", smsStatus.ErrorCode)
		}

		status = smsStatus.MessageStatus
		err = nil
		return status, err
	} else {
		err = fmt.Errorf("%s", resp.Status)
		return
	}
}
