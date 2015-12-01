// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package backend

import (
	"encoding/json"
	"io/ioutil"
	"locale"
	"net/http"
	"periwinkle"

	"github.com/jinzhu/gorm"
)

type TwilioNumber struct {
	ID     int64  `json:"number_id"`
	Number string `json:"number"`
	// TODO
}

func (o TwilioNumber) dbSchema(db *gorm.DB) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

type TwilioPool struct {
	UserID   string `json:"user_id"   sql:"type:varchar(255) REFERENCES users(id)          ON DELETE CASCADE  ON UPDATE RESTRICT"`
	GroupID  string `json:"group_id"  sql:"type:varchar(255) REFERENCES groups(id)         ON DELETE CASCADE  ON UPDATE RESTRICT"`
	NumberID int64  `json:"number_id" sql:"type:bigint       REFERENCES twilio_numbers(id) ON DELETE RESTRICT ON UPDATE RESTRICT"`
}

type IncomingNumbers struct {
	PhoneNumbers []IncomingNumber `json:"incoming_phone_numbers"`
}

type IncomingNumber struct {
	Number string `json:"phone_number"`
}

func (o TwilioPool) dbSchema(db *gorm.DB) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func GetAllUsedTwilioNumbers(db *gorm.DB) (ret []TwilioNumber) {
	var twilioNum []TwilioNumber
	if result := db.Find(&twilioNum); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return twilioNum

}

func GetTwilioPoolByUserID(db *gorm.DB, userid string) []TwilioPool {

	var o []TwilioPool
	if result := db.Where("user_id = ?", userid).Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
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

			if result := db.Where("id = ?", twilioPools[i].NumberID).First(&usedNums); result.Error != nil {
				if result.RecordNotFound() {
					return nil
				}
				dbError(result.Error)
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
			return ""
		}
		dbError(result.Error)
	}

	var twilioNum TwilioNumber
	if result := db.Where("id = ?", o.NumberID).First(&twilioNum); result.Error != nil {
		if result.RecordNotFound() {
			return ""
		}
		dbError(result.Error)
	}

	return twilioNum.Number
}

func AssignTwilioNumber(db *gorm.DB, userid string, groupid string, twilioNum string) *TwilioPool {

	num := TwilioNumber{
		Number: twilioNum,
	}

	if err := db.FirstOrCreate(&num).Error; err != nil {
		dbError(err)
	}

	o := TwilioPool{
		UserID:   userid,
		GroupID:  groupid,
		NumberID: num.ID,
	}

	if err := db.Create(&o).Error; err != nil {
		dbError(err)
	}

	return &o

}

func GetGroupByUserAndTwilioNumber(db *gorm.DB, userid string, twilioNum string) *Group {

	var num TwilioNumber

	if result := db.Where("number = ?", twilioNum).First(&num); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}

	var o TwilioPool

	if result := db.Where(&TwilioPool{UserID: userid, NumberID: num.ID}).First(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}

	var group Group

	if result := db.Where("id = ?", o.GroupID).First(&group); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}

	return &group
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

func DeleteUnusedTwilioNumber(db *gorm.DB, num string) error {

	var twilioNum TwilioNumber
	if result := db.Where("number = ?", num).First(&twilioNum); result.Error != nil {
		if result.RecordNotFound() {
			periwinkle.Logf("RecordNotFound")
			return result.Error
		}
		dbError(result.Error)
	}

	var twilioPool TwilioPool
	result := db.Where("number_id = ?", twilioNum.ID).First(&twilioPool)
	if result.Error != nil {
		if result.RecordNotFound() {

			if result := db.Where("number = ?", num).Delete(&TwilioNumber{}); result.Error != nil {
				dbError(result.Error)
			}

			periwinkle.Logf("RecordNotFound")
			return locale.UntranslatedError(result.Error)
		}
		dbError(result.Error)
	}

	return nil
}
