// Copyright 2015 Zhandos Suleimenov

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"periwinkle/twilio"
)

// function returns  a phone number and Status
//if successful, returns a new phone number and OK

func NewPhoneNum() (string, error) {

	// account SID for Twilio account
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")

	// gets url for available numbers

	availNum_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/AvailablePhoneNumbers/US/Local.json?SmsEnabled=true&MmsEnabled=true"

	// gets url for a new phone number

	newPhoneNum_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/IncomingPhoneNumbers.json"

	client := &http.Client{}

	req, err := http.NewRequest("GET", availNum_url, nil)

	if err != nil {
		log.Println(err)
		return "", err
	}

	req.SetBasicAuth(account_sid, auth_token)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	if resp.StatusCode == 302 {

		url, err := resp.Location()
		if err != nil {
			log.Println(err)
			return "", err
		}

		req, err = http.NewRequest("GET", url.String(), nil)
		if err != nil {
			log.Println(err)
			return "", err
		}

		req.SetBasicAuth(account_sid, auth_token)
		resp, err = client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Println(err)
			return "", err
		}

		if resp.StatusCode != 200 {

			return "", fmt.Errorf("%s", resp.Status)
		}

	} else if resp.StatusCode == 200 {

		//continue

	} else {
		return "", fmt.Errorf("%s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	avail_number := twilio.Avail_ph_num{}
	json.Unmarshal(body, &avail_number)

	if len(avail_number.PhoneNumberList) != 0 {

		number := avail_number.PhoneNumberList[0].PhoneNumber

		val := url.Values{}
		val.Set("PhoneNumber", avail_number.PhoneNumberList[0].PhoneNumber)
		val.Set("SmsUrl", "http://twimlets.com/echo?Twiml=%3CResponse%3E%3C%2FResponse%3E")

		req, err = http.NewRequest("POST", newPhoneNum_url, bytes.NewBuffer([]byte(val.Encode())))
		if err != nil {
			log.Println(err)
			return "", err
		}

		req.SetBasicAuth(account_sid, auth_token)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err = client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Println(err)
			return "", err
		}

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return "", fmt.Errorf("%s", resp.Status)
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return "", err
		}
		log.Println(string(body))
		return number, nil

	}

	log.Println("There are no available phone numbers!!!")
	return "", fmt.Errorf("There are no available phone numbers!!!")

}
