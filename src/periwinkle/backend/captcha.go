// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal

package backend

import (
	"io"
	"time"

	"github.com/dchest/captcha"
	"github.com/jinzhu/gorm"
)

const (
	// Default number of digits in captcha solution.
	DefaultLen = 6
	// Expiration time of captchas used by default store.
	DefaultExpiration = 20 * time.Minute
	//	Default Captcha Image Width
	DefaultWidth = 640
	//	Default Captcha Image Height
	DefaultHeight = 480
)

type Captcha struct {
	ID         string
	Value      string
	Token      string
	Expiration time.Time
}

func (o Captcha) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

func NewCaptcha(db *gorm.DB) *Captcha {
	o := Captcha{
		ID:    captcha.New(),
		Value: string(captcha.RandomDigits(DefaultLen)),
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return &o
}

func UseCaptcha(db *gorm.DB, token string) bool {
	panic("TODO")
}

func CheckCaptcha(db *gorm.DB, userInput string, captchaID string) bool {
	o := GetCaptchaByID(db, captchaID)
	if o == nil {
		return false
	}
	return userInput == o.Value
}

func GetCaptchaByID(db *gorm.DB, id string) *Captcha {
	var o Captcha
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func (o *Captcha) MarshalPNG(w io.Writer) error {
	// TODO: generate PNG and write it to w
	return captcha.WriteImage(w, o.ID, DefaultWidth, DefaultHeight)
}

func (o *Captcha) MarshalWAV(w io.Writer) error {
	// TODO: generate WAV and write it to w
	return captcha.WriteAudio(w, o.ID, "en")
}

func (o *Captcha) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}
