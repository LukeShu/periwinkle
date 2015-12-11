// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package backend

import (
	"locale"
	"periwinkle"
)

type TwilioNumber struct {
	ID     int64  `json:"number_id"`
	Number string `json:"number"`
	// TODO
}

func (o TwilioNumber) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

type TwilioPool struct {
	UserID   string `json:"user_id"   sql:"type:varchar(255) REFERENCES users(id)          ON DELETE CASCADE  ON UPDATE RESTRICT"`
	GroupID  string `json:"group_id"  sql:"type:varchar(255) REFERENCES groups(id)         ON DELETE CASCADE  ON UPDATE RESTRICT"`
	NumberID int64  `json:"number_id" sql:"type:bigint       REFERENCES twilio_numbers(id) ON DELETE RESTRICT ON UPDATE RESTRICT"`
}

func (o TwilioPool) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

type TwilioSMSMessage struct {
	MessageSID    string `json:"MessageSid" sql:"type:varchar(34)"`
	MessageStatus string
	ErrorCode     string
}

func (o TwilioSMSMessage) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func (o *TwilioSMSMessage) Save(db *periwinkle.Tx) {
	if err := db.Save(o).Error; err != nil {
		dbError(err)
	}
}

func (o *TwilioSMSMessage) Delete(db *periwinkle.Tx) {
	if err := db.Where("message_s_id = ?", o.MessageSID).Delete(TwilioSMSMessage{}).Error; err != nil {
		dbError(err)
	}
}

func GetTwilioSMSMessageBySID(db *periwinkle.Tx, sid string) *TwilioSMSMessage {
	var o TwilioSMSMessage
	if result := db.First(&o, "message_s_id = ?", sid); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return &o
}

func GetAllUsedTwilioNumbers(db *periwinkle.Tx) (ret []TwilioNumber) {
	var twilioNum []TwilioNumber
	if result := db.Find(&twilioNum); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return twilioNum

}

func GetTwilioPoolByUserID(db *periwinkle.Tx, userid string) []TwilioPool {

	var o []TwilioPool
	if result := db.Where("user_id = ?", userid).Find(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return o
}

func GetTwilioNumberByUserAndGroup(db *periwinkle.Tx, userid string, groupid string) string {

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

func AssignTwilioNumber(db *periwinkle.Tx, userid string, groupid string, twilioNum string) *TwilioPool {

	num := TwilioNumber{}
	err := db.Where(TwilioNumber{Number: twilioNum}).FirstOrCreate(&num).Error

	if err != nil {
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

func GetGroupByUserAndTwilioNumber(db *periwinkle.Tx, userid string, twilioNum string) *Group {

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

func GetTwilioNumberByID(db *periwinkle.Tx, id int64) *TwilioNumber {

	var twilio_num TwilioNumber

	if result := db.Where("id = ?", id).First(&twilio_num); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}

	return &twilio_num
}

func DeleteUnusedTwilioNumber(db *periwinkle.Tx, num string) locale.Error {
	var twilioNum TwilioNumber
	if result := db.Where("number = ?", num).First(&twilioNum); result.Error != nil {
		if result.RecordNotFound() {
			periwinkle.Logf("The number is already deleted!!!")
			return nil
		}
		dbError(result.Error)
		return locale.UntranslatedError(result.Error)
	}

	var twilioPool TwilioPool
	result := db.Where("number_id = ?", twilioNum.ID).First(&twilioPool)
	if result.Error != nil {
		if result.RecordNotFound() {

			o := db.Where("number = ?", num).Delete(&TwilioNumber{})
			if o.Error != nil {
				dbError(o.Error)
				return locale.UntranslatedError(o.Error)
			}
			periwinkle.Logf("The number is deleted")
			return nil
		}
		dbError(result.Error)
		return locale.UntranslatedError(result.Error)
	}
	periwinkle.Logf("The number is used for a twilio pool")
	return nil
}
