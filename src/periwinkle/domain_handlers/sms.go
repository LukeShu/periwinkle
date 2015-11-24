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
	"os"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/twilio"
	"postfixpipe"
	"strings"
	"time"

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
func sender(message mail.Message, sms_to string, db *gorm.DB, cfg *periwinkle.Cfg) (status string, err error) {

	group := message.Header.Get("From")
	user := backend.GetUserByAddress(db, "sms", sms_to)

	sms_from := backend.GetTwilioNumberByUserAndGroup(db, user.Id, strings.Split(group, "@")[0])
	sms_body := message.Header.Get("Subject")
	//sms_body, err := ioutil.ReadAll(message.Body)
	//if err != nil {
	//	return "", err
	//}

	// account SID for Twilio account
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")

	messages_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/Messages.json"

	v := url.Values{}
	v.Set("From", sms_from)
	v.Set("To", sms_to)
	v.Set("Body", string(sms_body))
	v.Set("StatusCallback", "http://"+cfg.WebRoot+"/callbacks/twilio-sms")

	client := &http.Client{}

	req, err := http.NewRequest("POST", messages_url, bytes.NewBuffer([]byte(v.Encode())))
	if err != nil {
		return
	}
	req.SetBasicAuth(account_sid, auth_token)
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
		sms_status, err := SmsWaitForCallback(message.Sid)

		if err != nil {
			return "", err
		}

		if sms_status.MessageStatus == "undelivered" || sms_status.MessageStatus == "failed" {
			return sms_status.MessageStatus, fmt.Errorf("%s", sms_status.ErrorCode)
		}
		if sms_status.MessageStatus == "queued" || sms_status.MessageStatus == "sending" || sms_status.MessageStatus == "sent" {
			time.Sleep(time.Second)
			sms_status, err = SmsWaitForCallback(message.Sid)

			if err != nil {
				return "", err
			}
		}

		if sms_status.MessageStatus == "undelivered" || sms_status.MessageStatus == "failed" {
			return sms_status.MessageStatus, fmt.Errorf("%s", sms_status.ErrorCode)
		}

		status = sms_status.MessageStatus
		err = nil
		return status, err
	} else {
		err = fmt.Errorf("%s", resp.Status)
		return
	}
}
