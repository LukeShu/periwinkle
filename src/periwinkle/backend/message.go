// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package backend

import (
	"maildir"

	"github.com/jinzhu/gorm"
)

type Message struct {
	Id      string
	GroupId string
	Unique  string
	// cached fields??????
}

func (o Message) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("group_id", "groups(id)", "CASCADE", "RESTRICT").
		AddUniqueIndex("filename_idx", "unique").
		Error
}

func NewMessage(db *gorm.DB, id string, group Group, unique maildir.Unique) Message {
	if id == "" {
		panic("Message Id can't be emtpy")
	}
	o := Message{
		Id:      id,
		GroupId: group.Id,
		Unique:  string(unique),
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return o
}

func GetMessageById(db *gorm.DB, id string) *Message {
	var o Message
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}
