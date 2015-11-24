// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package backend

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/jinzhu/gorm"
)

type errorList []error

func (errs errorList) Error() string {
	strs := make([]string, len(errs))
	for i, err := range errs {
		strs[i] = err.Error()
	}
	return " - " + strings.Join(strs, "\n - ")
}

func errHelper(errs *[]error, err error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}

func DbSchema(db *gorm.DB) error {
	errs := []error{}
	errHelper(&errs, (Captcha{}).dbSchema(db))
	errHelper(&errs, (Medium{}).dbSchema(db))
	errHelper(&errs, (Group{}).dbSchema(db))
	errHelper(&errs, (Message{}).dbSchema(db)) // must come after Group
	errHelper(&errs, (User{}).dbSchema(db))
	errHelper(&errs, (Session{}).dbSchema(db)) // must come after User
	errHelper(&errs, (ShortURL{}).dbSchema(db))
	errHelper(&errs, (UserAddress{}).dbSchema(db))  // must come after User and Medium
	errHelper(&errs, (Subscription{}).dbSchema(db)) // must come after Group and UserAddress
	errHelper(&errs, (TwilioNumber{}).dbSchema(db))
	errHelper(&errs, (TwilioPool{}).dbSchema(db)) // must come after TwilioNumber, User, and Group
	errHelper(&errs, (Admin{}).dbSchema(db))
	return errorList(errs)
}

func DbDrop(db *gorm.DB) error {
	// This must be in the reverse order of DbSchema()
	errs := []error{}
	errHelper(&errs, db.DropTable(&Admin{}).Error)
	errHelper(&errs, db.DropTable(&TwilioPool{}).Error)
	errHelper(&errs, db.DropTable(&TwilioNumber{}).Error)
	errHelper(&errs, db.DropTable(&Subscription{}).Error)
	errHelper(&errs, db.DropTable(&UserAddress{}).Error)
	errHelper(&errs, db.DropTable(&ShortURL{}).Error)
	errHelper(&errs, db.DropTable(&Session{}).Error)
	errHelper(&errs, db.DropTable(&User{}).Error)
	errHelper(&errs, db.DropTable(&Message{}).Error)
	errHelper(&errs, db.DropTable(&Group{}).Error)
	errHelper(&errs, db.DropTable(&Medium{}).Error)
	errHelper(&errs, db.DropTable(&Captcha{}).Error)
	return errorList(errs)
}

func DbSeed(db *gorm.DB) error {
	errs := []error{}
	errHelper(&errs, (Medium{}).dbSeed(db))
	errHelper(&errs, (Group{}).dbSeed(db))
	return errorList(errs)
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var alphabetLen = big.NewInt(int64(len(alphabet)))

func randomString(size int) string {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bigint, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			panic(err)
		}
		bytes[i] = alphabet[bigint.Int64()]
	}
	return string(bytes[:])
}
