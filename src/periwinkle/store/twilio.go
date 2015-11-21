// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov
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
	Id     int64  `json:"number_id"`
	Number string `json:"number"`
	// TODO
}

func (o TwilioNumber) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

type TwilioPool struct {
	UserId   string `json:"user_id"`
	GroupId  string `json:"group_id"`
	NumberId int64  `json:"number_id"`
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
	var twilio_num []TwilioNumber
	if result := db.Find(&twilio_num); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return twilio_num

}

func GetTwilioPoolByUserId(db *gorm.DB, userid string) []TwilioPool {

	var o []TwilioPool
	if result := db.Where("user_id = ?", userid).Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return nil
		}
		panic(result.Error)
	}
	return o
}

func GetUnusedTwilioNumbersByUser(db *gorm.DB, userid string) []string {

	str := []string{}
	all_twilio_num := GetAllExistingTwilioNumbers()
	twilio_pools := GetTwilioPoolByUserId(db, userid)

	var used_nums TwilioNumber
	var isNumberUsed bool

	for _, all_num := range all_twilio_num {
		isNumberUsed = false
		for i, _ := range twilio_pools {

			if result := db.Where("number_id = ?", twilio_pools[i].NumberId).First(&used_nums); result.Error != nil {
				if result.RecordNotFound() {
					log.Println("RecordNotFound")
					return nil
				}
				panic(result.Error)
			}

			if all_num == used_nums.Number {
				isNumberUsed = true
				break
			}
		}

		if isNumberUsed == false {
			str = append(str, all_num)
		}
	}

	return str

}

func GetTwilioNumberByUserAndGroup(db *gorm.DB, userid string, groupid string) string {

	var o TwilioPool
	if result := db.Where(&TwilioPool{UserId: userid, GroupId: groupid}).First(&o); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return ""
		}
		panic(result.Error)
	}

	var twilio_num TwilioNumber
	if result := db.Where("number_id = ?", o.NumberId).First(&twilio_num); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return ""
		}
		panic(result.Error)
	}

	return twilio_num.Number
}

func AssignTwilioNumber(db *gorm.DB, userid string, groupid string, twilio_num string) *TwilioPool {

	num := TwilioNumber{
		Number: twilio_num,
	}

	if err := db.Create(&num).Error; err != nil {
		panic(err)
	}

	o := TwilioPool{
		UserId:   userid,
		GroupId:  groupid,
		NumberId: num.Id,
	}

	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}

	return &o

}

func GetGroupByUserAndTwilioNumber(db *gorm.DB, userid string, twilio_num string) *Group {

	var num TwilioNumber

	if result := db.Where("number = ?", twilio_num).First(&num); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return nil
		}
		panic(result.Error)
	}

	var o TwilioPool

	if result := db.Where(&TwilioPool{UserId: userid, NumberId: num.Id}).First(&o); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return nil
		}
		panic(result.Error)
	}

	var group Group

	if result := db.Where("group_id = ?", o.GroupId).First(&group); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return nil
		}
		panic(result.Error)
	}

	return &group
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
