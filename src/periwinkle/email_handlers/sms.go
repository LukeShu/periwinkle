// Copyright 2015 Davis Webb
// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Luke Shumaker

package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"periwinkle/cfg"
	"periwinkle/store"
	"strings"
	"time"
)

var message_status, error_code string

func HandleSMS(r io.Reader, name string) int {
	panic("TODO")
}

func SmsHttpCallback(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Printf("%v", err)
	}

	message_status = values.Get("MessageStatus")
	error_code = values.Get("ErrorCode")
}

// Returns the status of the message: queued, sending, sent,
// delivered, undelivered, failed.  If an error occurs, it returns
// Error.
func sender(message mail.Message, to string) (status string, err error) {
	message_status = ""
	error_code = ""

	group := message.Header.Get("From")
	user := store.GetUserByAddress(cfg.DB, "email", message.Header.Get("From"))

	sms_from := group // TODO: numberFor(group)
	sms_to := strings.Split(to, "@")[0]
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
		time.Sleep(time.Second)
		if error_code != "" {
			return message_status, fmt.Errorf("%s", error_code)
		}
		if message_status == "queued" || message_status == "sending" || message_status == "sent" {
			time.Sleep(time.Second)
		}
		status = message_status
		err = nil
		return
	} else {
		err = fmt.Errorf("%s", resp.Status)
		return
	}
}
