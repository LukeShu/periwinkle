// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal

package backend

import (
	"io"
	"locale"
	"periwinkle"
	"strings"
	"time"

	"github.com/dchest/captcha"
)

const (
	// Default number of digits in captcha solution.
	defaultLen = 6
	// Expiration time of captchas used by default backend.
	defaultExpiration = 20 * time.Minute
	// Default Captcha Image Width
	defaultWidth = 640
	// Default Captcha Image Height
	defaultHeight = 480
)

type Captcha struct {
	ID         string    `json:"-"`
	Value      string    `json:"value"`
	Token      string    `json:"-"`
	Expiration time.Time `json:"expiration_time"`
}

func (o Captcha) dbSchema(db *periwinkle.Tx) locale.Error {
	return locale.UntranslatedError(db.CreateTable(&o).Error)
}

func NewCaptcha(db *periwinkle.Tx) *Captcha {
	o := Captcha{
		ID:    captcha.New(),
		Value: string(captcha.RandomDigits(defaultLen)),
	}
	if err := db.Create(&o).Error; err != nil {
		dbError(err)
	}
	return &o
}

func UseCaptcha(db *periwinkle.Tx, id, token string) bool {
	o := GetCaptchaByID(db, id)
	if o == nil {
		panic("Captcha " + id + " does not exist.")
	}
	if strings.Compare(token, "true") == 0 {
		// destroy captcha
		db.Delete(&o)
		return true
	}
	// destroy captcha
	db.Delete(&o)
	return false
}

func CheckCaptcha(db *periwinkle.Tx, userInput string, captchaID string) bool {
	o := GetCaptchaByID(db, captchaID)
	if o == nil {
		return false
	}
	return userInput == o.Value
}

func GetCaptchaByID(db *periwinkle.Tx, id string) *Captcha {
	var o Captcha
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		dbError(result.Error)
	}
	return &o
}

func (o *Captcha) MarshalPNG(w io.Writer) locale.Error {
	// TODO: generate PNG and write it to w
	return locale.UntranslatedError(captcha.WriteImage(w, o.ID, defaultWidth, defaultHeight))
}

func (o *Captcha) MarshalWAV(w io.Writer) locale.Error {
	// TODO: generate WAV and write it to w
	return locale.UntranslatedError(captcha.WriteAudio(w, o.ID, "en"))
}

func (o *Captcha) Save(db *periwinkle.Tx) {
	if err := db.Save(o).Error; err != nil {
		dbError(err)
	}
}
