// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend

import (
	"locale"
	"periwinkle"
	"time"
)

type Session struct {
	ID       string    `json:"session_id"`
	UserID   string    `json:"user_id" sql:"type:varchar(255) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	LastUsed time.Time `json:"-"`
}

func (o Session) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func NewSession(db *periwinkle.Tx, user *User, password string) *Session {
	if user == nil || !user.CheckPassword(password) {
		return nil
	}
	o := Session{
		ID:       RandomString(24),
		UserID:   user.ID,
		LastUsed: time.Now(),
	}
	if err := db.Create(&o).Error; err != nil {
		dbError(err)
	}
	return &o
}

func GetSessionByID(db *periwinkle.Tx, id string) *Session {
	var o Session
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return &o
}

func (o *Session) Delete(db *periwinkle.Tx) {
	db.Delete(o)
}

func (o *Session) Save(db *periwinkle.Tx) {
	if err := db.Save(o).Error; err != nil {
		dbError(err)
	}
}
