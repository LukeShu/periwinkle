// Copyright 2015 Luke Shumaker

package store

import (	
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type TwilioNumber struct {
	Id     int64
	Number string
	// TODO
}

func (o TwilioNumber) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

type TwilioPool struct {
	UserId   string
	GroupId  string
	NumberId string
}

type Incoming_numbers struct {
	Phone_numbers []Incoming_number `json:"incoming_phone_numbers"`
}

type Incoming_number struct {
	Number string `json:"phone_number"`
}

func (o TwilioPool) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT").
		AddForeignKey("number_id", "twilio_numbers(id)", "RESTRICT", "RESTRICT").
		Error
}

func GetAllTwilioNumbers(db *gorm.DB) (ret []TwilioNumber) {
	panic("TODO")
}

func GetAllExistingTwilioNumbers() []string {

	// account SID for Twilio account
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")

	// gets url for the numbers we own in the Twilio Account
	incoming_num_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/IncomingPhoneNumbers.json"

	client := &http.Client{}

	req, err := http.NewRequest("GET", incoming_num_url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	req.SetBasicAuth(account_sid, auth_token)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return nil
	}	

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		return nil
	}

	numbers := Incoming_numbers{}
	if err := json.Unmarshal(body, &numbers); err != nil {
		log.Println(err)
		return nil
	}

	if len(numbers.Phone_numbers) > 0 {	

		existing_numbers := make([]string, len(numbers.Phone_numbers))

		for i, num := range numbers.Phone_numbers {
			existing_numbers[i] = num.Number
		}

		return existing_numbers

	} else {
		log.Println("You do not have a number in your Twilio account")
		return nil
	}	

}
