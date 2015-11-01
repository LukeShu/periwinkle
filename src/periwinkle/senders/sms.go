// Copyright 2015 Davis Webb
// Copyright 2015 Zhandos Suleimenov

package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"periwinkle/cfg"
	"time"
)

var message_status, error_code string

func Url_handler(w http.ResponseWriter, req *http.Request) {
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
func sender(reader io.Reader) string {
	message_status = ""
	error_code = ""
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("%v", err)
		return "Error"
	}

	message := make(map[string]string)
	json.Unmarshal(body, &message)

	// account SID for Twilio account
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")

	messages_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/Messages.json"

	host_name, err := os.Hostname()
	if err != nil {
		fmt.Printf("%v", err)
		return "Error"
	}

	v := url.Values{}
	v.Set("From", message["From"])
	v.Set("To", message["To"])
	v.Set("Body", message["Body"])
	v.Set("StatusCallback", "http://"+host_name+cfg.WebAddr+"/webui/twilio/sms")

	client := &http.Client{}

	req, err := http.NewRequest("POST", messages_url, bytes.NewBuffer([]byte(v.Encode())))
	if err != nil {
		fmt.Printf("%v\n", err)
		return "Error"
	}
	req.SetBasicAuth(account_sid, auth_token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%v\n", err)
		return "Error"
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {

		time.Sleep(time.Second)
		if error_code != "" {
			fmt.Printf("%v\n", error_code)
		}
		if message_status == "queued" || message_status == "sending" || message_status == "sent" {
			time.Sleep(time.Second)
		}

		return message_status
	} else {
		fmt.Printf("%v\n", resp.Status)
		return "Error"
	}
}
