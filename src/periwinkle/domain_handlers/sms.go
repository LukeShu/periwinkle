// Copyright 2015 Davis Webb
// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Luke Shumaker

package domain_handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"locale"
	"net/http"
	"net/mail"
	"net/url"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/twilio"
	"postfixpipe"
	"strings"
)

func HandleSMS(r io.Reader, name string, db *periwinkle.Tx, cfg *periwinkle.Cfg) postfixpipe.ExitStatus {
	message, uerr := mail.ReadMessage(r)
	if uerr != nil {
		periwinkle.LogErr(locale.UntranslatedError(uerr))
		return postfixpipe.EX_NOINPUT
	}

	group := message.Header.Get("From")
	user := backend.GetUserByAddress(db, "sms", name)

	smsFrom := backend.GetTwilioNumberByUserAndGroup(db, user.ID, strings.Split(group, "@")[0])

	if smsFrom == "" {
		twilio_num := twilio.GetUnusedTwilioNumbersByUser(cfg, db, user.ID)
		if twilio_num == nil {
			new_num, err := twilio.NewPhoneNum(cfg)
			if err != nil {
				periwinkle.LogErr(err)
				return postfixpipe.EX_UNAVAILABLE
			}
			backend.AssignTwilioNumber(db, user.ID, strings.Split(group, "@")[0], new_num)
			smsFrom = new_num
		} else {
			backend.AssignTwilioNumber(db, user.ID, strings.Split(group, "@")[0], twilio_num[0])
			smsFrom = twilio_num[0]
		}
	}

	smsBody := message.Header.Get("Subject")
	//smsBody, err := ioutil.ReadAll(message.Body)
	//if err != nil {
	//	return "", err
	//}

	messagesURL := "https://api.twilio.com/2010-04-01/Accounts/" + cfg.TwilioAccountID + "/Messages.json"

	v := url.Values{}
	v.Set("From", smsFrom)
	v.Set("To", name)
	v.Set("Body", string(smsBody))
	v.Set("StatusCallback", cfg.WebRoot+"/callbacks/twilio-sms")
	//host,_ := os.Hostname()
	//v.Set("StatusCallback", "http://" + host + ":8080/callbacks/twilio-sms")
	client := &http.Client{}

	req, uerr := http.NewRequest("POST", messagesURL, bytes.NewBuffer([]byte(v.Encode())))
	if uerr != nil {
		periwinkle.LogErr(locale.UntranslatedError(uerr))
		return postfixpipe.EX_UNAVAILABLE
	}
	req.SetBasicAuth(cfg.TwilioAccountID, cfg.TwilioAuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, uerr := client.Do(req)
	defer resp.Body.Close()
	if uerr != nil {
		periwinkle.LogErr(locale.UntranslatedError(uerr))
		return postfixpipe.EX_UNAVAILABLE
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return postfixpipe.EX_UNAVAILABLE
	}

	body, uerr := ioutil.ReadAll(resp.Body)
	if uerr != nil {
		periwinkle.LogErr(locale.UntranslatedError(uerr))
		return postfixpipe.EX_UNAVAILABLE
	}

	tmessage := twilio.Message{}
	json.Unmarshal([]byte(body), &tmessage)
	_, err := TwilioSMSWaitForCallback(cfg, tmessage.Sid)
	if err != nil {
		periwinkle.LogErr(err)
		return postfixpipe.EX_UNAVAILABLE
	}
	return postfixpipe.EX_OK
}
