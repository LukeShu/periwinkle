// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package backend

import (
	"crypto/rand"
	"locale"
	"math/big"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type errorList []locale.Error

var _ locale.Error = errorList{}

func (errs errorList) L10NString(l locale.Spec) string {
	strs := make([]string, len(errs))
	for i, err := range errs {
		strs[i] = " - " + strings.Replace(err.L10NString(l), "\n", "\n   ", -1)
	}
	return strings.Join(strs, "\n")
}

func (errs errorList) Error() string {
	return errs.L10NString("C")
}

func errHelper(errs *errorList, err locale.Error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}

func DbSchema(db *gorm.DB) error {
	errs := errorList{}
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
	return errs
}

func DbDrop(db *gorm.DB) error {
	// This must be in the reverse order of DbSchema()
	errs := errorList{}
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&Admin{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&TwilioPool{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&TwilioNumber{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&Subscription{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&UserAddress{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&ShortURL{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&Session{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&User{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&Message{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&Group{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&Medium{}).Error))
	errHelper(&errs, locale.UntranslatedError(db.DropTable(&Captcha{}).Error))
	return errs
}

func DbSeed(db *gorm.DB) error {
	errs := errorList{}
	errHelper(&errs, locale.UntranslatedError((Medium{}).dbSeed(db)))
	errHelper(&errs, locale.UntranslatedError((Group{}).dbSeed(db)))
	return errs
}

// Panic, but a little nicer
func dbError(err error) {
	switch e := err.(type) {
	case sqlite3.Error, *mysql.MySQLError:
		panic(locale.UntranslatedError(e))
	default:
		panic(locale.Errorf("Programmer Error: the programmer said this is a database error, but it's not: %s", e))
	}
}

// Panic, but a little nicer
func programmerError(str string) {
	panic(locale.Errorf("Programmer Error: %s", locale.Sprintf(str)))
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var alphabetLen = big.NewInt(int64(len(alphabet)))

func randomString(size int) string {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bigint, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			panic(err) // Luke says this is OK
		}
		bytes[i] = alphabet[bigint.Int64()]
	}
	return string(bytes[:])
}
