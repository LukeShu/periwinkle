// Copyright 2015 Zhandos Suleimenov

package twilio

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"locale"
	"net/http"
	"net/url"
	"periwinkle"
)

// function returns  a phone number and Status
// if successful, returns a new phone number and OK
func NewPhoneNum(cfg *periwinkle.Cfg) (string, locale.Error) {
	// gets url for available numbers
	availNumURL := "https://api.twilio.com/2010-04-01/Accounts/" + cfg.TwilioAccountID + "/AvailablePhoneNumbers/US/Local.json?SmsEnabled=true&MmsEnabled=true"

	// gets url for a new phone number
	newPhoneNumURL := "https://api.twilio.com/2010-04-01/Accounts/" + cfg.TwilioAccountID + "/IncomingPhoneNumbers.json"

	client := &http.Client{}

	req, err := http.NewRequest("GET", availNumURL, nil)
	if err != nil {
		return "", locale.UntranslatedError(err)
	}

	req.SetBasicAuth(cfg.TwilioAccountID, cfg.TwilioAuthToken)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", locale.UntranslatedError(err)
	}

	if resp.StatusCode == 302 {
		url, err := resp.Location()
		if err != nil {
			return "", locale.UntranslatedError(err)
		}

		req, err = http.NewRequest("GET", url.String(), nil)
		if err != nil {
			return "", locale.UntranslatedError(err)
		}

		req.SetBasicAuth(cfg.TwilioAccountID, cfg.TwilioAuthToken)
		resp, err = client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			return "", locale.UntranslatedError(err)
		}

		if resp.StatusCode != 200 {
			return "", locale.Errorf("%s", resp.Status)
		}

	} else if resp.StatusCode == 200 {

		//continue

	} else {
		return "", locale.Errorf("%s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", locale.UntranslatedError(err)
	}

	availNumber := AvailPhNum{}
	json.Unmarshal(body, &availNumber)

	if len(availNumber.PhoneNumberList) != 0 {

		number := availNumber.PhoneNumberList[0].PhoneNumber

		val := url.Values{}
		val.Set("PhoneNumber", availNumber.PhoneNumberList[0].PhoneNumber)
		val.Set("SmsUrl", "http://twimlets.com/echo?Twiml=%3CResponse%3E%3C%2FResponse%3E")

		req, err = http.NewRequest("POST", newPhoneNumURL, bytes.NewBuffer([]byte(val.Encode())))
		if err != nil {
			return "", locale.UntranslatedError(err)
		}

		req.SetBasicAuth(cfg.TwilioAccountID, cfg.TwilioAuthToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err = client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			return "", locale.UntranslatedError(err)
		}

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return "", locale.Errorf("%s", resp.Status)
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", locale.UntranslatedError(err)
		}
		return number, nil

	}

	return "", locale.Errorf("There are no available phone numbers!!!")
}
