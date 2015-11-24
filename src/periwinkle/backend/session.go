// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Session struct {
	Id       string    `json:"session_id"`
	UserId   string    `json:"user_id"`
	LastUsed time.Time `json:"-"`
}

func (o Session) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		Error
}

func NewSession(db *gorm.DB, user *User, password string) *Session {
	if user == nil || !user.CheckPassword(password) {
		return nil
	}
	o := Session{
		Id:       randomString(24),
		UserId:   user.Id,
		LastUsed: time.Now(),
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func GetSessionById(db *gorm.DB, id string) *Session {
	var o Session
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func (o *Session) Delete(db *gorm.DB) {
	db.Delete(o)
}

func (o *Session) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}
