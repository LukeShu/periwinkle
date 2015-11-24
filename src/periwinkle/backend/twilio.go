// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package backend

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"periwinkle"

	"github.com/jinzhu/gorm"
)

type TwilioNumber struct {
	ID     int64  `json:"number_id"`
	Number string `json:"number"`
	// TODO
}

func (o TwilioNumber) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

type TwilioPool struct {
	UserID   string `json:"user_id"`
	GroupID  string `json:"group_id"`
	NumberID int64  `json:"number_id"`
}

type IncomingNumbers struct {
	PhoneNumbers []IncomingNumber `json:"incoming_phone_numbers"`
}

type IncomingNumber struct {
	Number string `json:"phone_number"`
}

func (o TwilioPool) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT").
		AddForeignKey("number_id", "twilio_numbers(id)", "RESTRICT", "RESTRICT").
		Error
}

func GetAllUsedTwilioNumbers(db *gorm.DB) (ret []TwilioNumber) {
	var twilioNum []TwilioNumber
	if result := db.Find(&twilioNum); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return twilioNum

}

func GetTwilioPoolByUserID(db *gorm.DB, userid string) []TwilioPool {

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

func GetUnusedTwilioNumbersByUser(cfg *periwinkle.Cfg, db *gorm.DB, userid string) []string {

	str := []string{}
	allTwilioNum := GetAllExistingTwilioNumbers(cfg)
	twilioPools := GetTwilioPoolByUserID(db, userid)

	var usedNums TwilioNumber
	var isNumberUsed bool

	// TODO: no queries inside of loops!
	for _, allNum := range allTwilioNum {
		isNumberUsed = false
		for i := range twilioPools {

			if result := db.Where("number_id = ?", twilioPools[i].NumberID).First(&usedNums); result.Error != nil {
				if result.RecordNotFound() {
					log.Println("RecordNotFound")
					return nil
				}
				panic(result.Error)
			}

			if allNum == usedNums.Number {
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

func GetTwilioNumberByUserAndGroup(db *gorm.DB, userid string, groupid string) string {

	var o TwilioPool
	if result := db.Where(&TwilioPool{UserID: userid, GroupID: groupid}).First(&o); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return ""
		}
		panic(result.Error)
	}

	var twilioNum TwilioNumber
	if result := db.Where("number_id = ?", o.NumberID).First(&twilioNum); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return ""
		}
		panic(result.Error)
	}

	return twilioNum.Number
}

func AssignTwilioNumber(db *gorm.DB, userid string, groupid string, twilioNum string) *TwilioPool {

	num := TwilioNumber{
		Number: twilioNum,
	}

	if err := db.Create(&num).Error; err != nil {
		panic(err)
	}

	o := TwilioPool{
		UserID:   userid,
		GroupID:  groupid,
		NumberID: num.ID,
	}

	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}

	return &o

}

func GetGroupByUserAndTwilioNumber(db *gorm.DB, userid string, twilioNum string) *Group {

	var num TwilioNumber

	if result := db.Where("number = ?", twilioNum).First(&num); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return nil
		}
		panic(result.Error)
	}

	var o TwilioPool

	if result := db.Where(&TwilioPool{UserID: userid, NumberID: num.ID}).First(&o); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return nil
		}
		panic(result.Error)
	}

	var group Group

	if result := db.Where("group_id = ?", o.GroupID).First(&group); result.Error != nil {
		if result.RecordNotFound() {
			log.Println("RecordNotFound")
			return nil
		}
		panic(result.Error)
	}

	return &group
}

func GetAllExistingTwilioNumbers(cfg *periwinkle.Cfg) []string {
	// gets url for the numbers we own in the Twilio Account
	incomingNumURL := "https://api.twilio.com/2010-04-01/Accounts/" + cfg.TwilioAccountID + "/IncomingPhoneNumbers.json"

	client := &http.Client{}

	req, err := http.NewRequest("GET", incomingNumURL, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	req.SetBasicAuth(cfg.TwilioAccountID, cfg.TwilioAuthToken)

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

	numbers := IncomingNumbers{}
	if err := json.Unmarshal(body, &numbers); err != nil {
		log.Println(err)
		return nil
	}

	if len(numbers.PhoneNumbers) > 0 {

		existingNumbers := make([]string, len(numbers.PhoneNumbers))

		for i, num := range numbers.PhoneNumbers {
			existingNumbers[i] = num.Number
		}

		return existingNumbers

	} else {
		log.Println("You do not have a number in your Twilio account")
		return nil
	}

}
