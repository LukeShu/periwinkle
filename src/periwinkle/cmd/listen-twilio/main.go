// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/cmdutil"
	"periwinkle/twilio"
	"strings"
	"locale"
	"time"
)

const usage = `
Usage: %[1]s [-c CONFIG_FILE]
       %[1]s -h | --help
Repeatedly poll Twilio for new messages.

Options:
  -h, --help      Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`


var timeZero time.Time
var lastPoll time.Time

func main() {
	options := cmdutil.Docopt(usage)
	config := cmdutil.GetConfig(options["-c"].(string))

	for {
		time.Sleep(time.Second)
		conflict := config.DB.Do(func(tx *periwinkle.Tx) {
			numbers := backend.GetAllUsedTwilioNumbers(tx)

			for _, number := range numbers {
				checkNumber(config, tx, number)
			}
		})
		if conflict != nil {
			periwinkle.LogErr(conflict)
		}
		lastPoll = time.Now().UTC()
	}
}

func checkNumber(config *periwinkle.Cfg, tx *periwinkle.Tx, number backend.TwilioNumber) {
	url := "https://api.twilio.com/2010-04-01/Accounts/" + config.TwilioAccountID + "/Messages.json?To=" + number.Number
	if lastPoll != timeZero {
		url += "&DateSent>=" + strings.Split(lastPoll.UTC().String(), " ")[0]
	}
				
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(config.TwilioAccountID, config.TwilioAuthToken)
	resp, uerr := (&http.Client{}).Do(req)
	if uerr != nil {
		periwinkle.LogErr(locale.UntranslatedError(uerr))
	}

	defer resp.Body.Close()
	body, uerr := ioutil.ReadAll(resp.Body)
	if uerr != nil {
		periwinkle.LogErr(locale.UntranslatedError(uerr))
	}

	// converts JSON messages
	var page twilio.Paging
	json.Unmarshal([]byte(body), &page)

	for _, message := range page.Messages {
		timeSend, uerr := time.Parse(time.RFC1123Z, message.DateSent)
		if uerr != nil {
			periwinkle.LogErr(locale.UntranslatedError(uerr))
			continue
		}
		if timeSend.Unix() < lastPoll.Unix() {
			periwinkle.Logf("message %q older than our last poll; ignoring", message.Sid)
			continue
		}
		user := backend.GetUserByAddress(tx, "sms", message.From)
		if user == nil {
			periwinkle.Logf("could not figure out which user has number %q", message.From)
			continue
		}
		group := backend.GetGroupByUserAndTwilioNumber(tx, user.ID, message.To)
		if group == nil {
			periwinkle.Logf("could not figure out which group this is meant for: user: %q, number: %q", user.ID, message.To)
			continue
		}
		periwinkle.Logf("received message for group %q", group.ID)
		MessageBuilder{
			Maildir: config.Mailstore,
			Headers: map[string]string{
				"To":      group.ID + "@" + config.GroupDomain,
				"From":    backend.UserAddress{Medium: "sms", Address: message.From}.AsEmailAddress(),
				"Subject": user.ID + ": " + message.Body,
			},
			Body: "",
		}.Done()
	}
}
