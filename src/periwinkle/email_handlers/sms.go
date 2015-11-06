// Copyright 2015 Davis Webb
// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Luke Shumaker

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"periwinkle/cfg"
	"periwinkle/store"
	"periwinkle/twilio"
	"strings"
	"time"
)

func HandleSMS(r io.Reader, name string, db *gorm.DB) uint8 {
	panic("TODO")
}

// Returns the status of the message: queued, sending, sent,
// delivered, undelivered, failed.  If an error occurs, it returns
// Error.
func sender(message mail.Message, to string) (status string, err error) {
	group := message.Header.Get("From")
	user := store.GetUserByAddress(cfg.DB, "email", message.Header.Get("From"))

	sms_from := group                   // TODO: numberFor(group)
	sms_to := strings.Split(to, "@")[1] //test 0 or 1
	sms_body := user.FullName + ":" + message.Header.Get("Subject")

	// account SID for Twilio account
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")

	messages_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/Messages.json"

	v := url.Values{}
	v.Set("From", sms_from)
	v.Set("To", sms_to)
	v.Set("Body", sms_body)
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
		return
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		message := twilio.Message{}
		json.Unmarshal([]byte(body), &message)
		sms_status, err := SmsWaitForCallback(message.Sid)

		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second)
		if sms_status.ErrorCode != "" {
			return sms_status.MessageStatus, fmt.Errorf("%s", sms_status.ErrorCode)
		}
		if sms_status.MessageStatus == "queued" || sms_status.MessageStatus == "sending" || sms_status.MessageStatus == "sent" {
			time.Sleep(time.Second)
		}
		status = sms_status.MessageStatus
		err = fmt.Errorf("%s", sms_status.ErrorCode)
		return status, err
	} else {
		err = fmt.Errorf("%s", resp.Status)
		return
	}
}
