// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

// This file is named backend.go because it was pulled out of
// backend/twilio.go.

package twilio

import (
	"encoding/json"
	"io/ioutil"
	"locale"
	"net/http"
	"periwinkle"
	"periwinkle/backend"
)

type IncomingNumbers struct {
	PhoneNumbers []IncomingNumber `json:"incoming_phone_numbers"`
}

type IncomingNumber struct {
	Number string `json:"phone_number"`
}

func GetAllExistingTwilioNumbers(cfg *periwinkle.Cfg) []string {
	// gets url for the numbers we own in the Twilio Account
	incomingNumURL := "https://api.twilio.com/2010-04-01/Accounts/" + cfg.TwilioAccountID + "/IncomingPhoneNumbers.json"

	client := &http.Client{}

	req, err := http.NewRequest("GET", incomingNumURL, nil)
	if err != nil {
		periwinkle.LogErr(locale.UntranslatedError(err))
		return nil
	}

	req.SetBasicAuth(cfg.TwilioAccountID, cfg.TwilioAuthToken)

	resp, err := client.Do(req)

	if err != nil {
		periwinkle.LogErr(locale.UntranslatedError(err))
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		periwinkle.LogErr(locale.UntranslatedError(err))
		return nil
	}

	if resp.StatusCode != 200 {
		periwinkle.Logf("Response code %s", resp.Status)
		return nil
	}

	numbers := IncomingNumbers{}
	if err := json.Unmarshal(body, &numbers); err != nil {
		periwinkle.LogErr(locale.UntranslatedError(err))
		return nil
	}

	if len(numbers.PhoneNumbers) > 0 {

		existingNumbers := make([]string, len(numbers.PhoneNumbers))

		for i, num := range numbers.PhoneNumbers {
			existingNumbers[i] = num.Number
		}

		return existingNumbers

	} else {
		return nil
	}
}

func DeleteUnusedTwilioNumbers(cfg *periwinkle.Cfg) {
	conflict := cfg.DB.Do(func(db *periwinkle.Tx) {
		twilio_num := GetAllExistingTwilioNumbers(cfg)
		if twilio_num == nil {
			return
		}
		for _, v := range twilio_num {
			backend.DeleteUnusedTwilioNumber(db, v)
		}
	})
	if conflict != nil {
		panic(conflict) // FIXME
	}
}

func GetUnusedTwilioNumbersByUser(cfg *periwinkle.Cfg, db *periwinkle.Tx, userid string) []string {

	str := []string{}
	allTwilioNum := GetAllExistingTwilioNumbers(cfg)
	twilioPools := backend.GetTwilioPoolByUserID(db, userid)

	var isNumberUsed bool

	for _, allNum := range allTwilioNum {
		isNumberUsed = false
		for i := range twilioPools {

			usedNum := backend.GetTwilioNumberByID(db, twilioPools[i].NumberID)

			if allNum == usedNum.Number {
				isNumberUsed = true
				break
			}
		}

		if isNumberUsed == false {
			str = append(str, allNum)
		}
	}

	return str

}
