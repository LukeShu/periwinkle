// Copyright 2015 Luke Shumaker

package store

import (
	"github.com/jinzhu/gorm"
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
