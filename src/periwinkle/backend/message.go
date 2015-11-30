// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend

import (
	"locale"
	"maildir"

	"github.com/jinzhu/gorm"
)

type Message struct {
	ID      string
	GroupID string
	Unique  string
	// cached fields??????
}

func (o Message) dbSchema(db *gorm.DB) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).
		AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT").
		AddUniqueIndex("filename_idx", "unique").
		Error)
}

func NewMessage(db *gorm.DB, id string, group Group, unique maildir.Unique) Message {
	if id == "" {
		programmerError("Message ID can't be emtpy")
	}
	o := Message{
		ID:      id,
		GroupID: group.ID,
		Unique:  string(unique),
	}
	if err := db.Create(&o).Error; err != nil {
		dbError(err)
	}
	return o
}

func GetMessageByID(db *gorm.DB, id string) *Message {
	var o Message
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return &o
}
