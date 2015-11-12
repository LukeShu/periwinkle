// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package store

import (
	"crypto/rand"
	"encoding/json"
	"github.com/jinzhu/gorm"
	he "httpentity"
	"httpentity/util" // heutil
	"io"
	"jsondiff"
	"math/big"
	"strings"
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
	errHelper(&errs, (GroupAddress{}).dbSchema(db)) // must come after Group and Medium
	errHelper(&errs, (Message{}).dbSchema(db))      // must come after Group
	errHelper(&errs, (User{}).dbSchema(db))
	errHelper(&errs, (Session{}).dbSchema(db)) // must come after User
	errHelper(&errs, (ShortUrl{}).dbSchema(db))
	errHelper(&errs, (UserAddress{}).dbSchema(db))  // must come after User and Medium
	errHelper(&errs, (Subscription{}).dbSchema(db)) // must come after Group and UserAddress
	return errorList(errs)
}

func DbDrop(db *gorm.DB) error {
	// This must be in the reverse order of DbSchema()
	errs := []error{}
	errHelper(&errs, db.DropTable(&Subscription{}).Error)
	errHelper(&errs, db.DropTable(&UserAddress{}).Error)
	errHelper(&errs, db.DropTable(&ShortUrl{}).Error)
	errHelper(&errs, db.DropTable(&Session{}).Error)
	errHelper(&errs, db.DropTable(&User{}).Error)
	errHelper(&errs, db.DropTable(&Message{}).Error)
	errHelper(&errs, db.DropTable(&GroupAddress{}).Error)
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

func safeDecodeJSON(in interface{}, out interface{}) *he.Response {
	decoder, ok := in.(*json.Decoder)
	if !ok {
		ret := he.StatusUnsupportedMediaType(heutil.NetString(k("PUT and POST requests must have a document media type")))
		return &ret
	}
	var tmp interface{}
	err := decoder.Decode(&tmp)
	if err != nil {
		ret := he.StatusUnsupportedMediaType(heutil.NetPrintf(k("Couldn't parse: %v"), err))
		return &ret
	}
	str, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, out)
	if err != nil {
		ret := he.StatusUnsupportedMediaType(heutil.NetPrintf(k("Request body didn't have expected structure (field had wrong data type): %v"), err))
		return &ret
	}
	if !jsondiff.Equal(tmp, out) {
		diff, err := jsondiff.NewJSONPatch(tmp, out)
		if err != nil {
			panic(err)
		}
		entity := heutil.NetMap{
			"message": k("Request body didn't have expected structure (extra or missing fields).  The included diff would make the request acceptable."),
			"diff":    diff,
		}
		ret := he.StatusUnsupportedMediaType(entity)
		return &ret
	}
	return nil
}

// Simple dump to JSON, good for most entities
func defaultEncoders(o interface{}) map[string]func(io.Writer) error {
	return map[string]func(io.Writer) error{
		"application/json": func(w io.Writer) error {
			bytes, err := json.Marshal(o)
			if err != nil {
				return err
			}
			_, err = w.Write(bytes)
			return err
		},
	}
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var alphabet_len = big.NewInt(int64(len(alphabet)))

func randomString(size int) string {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bigint, err := rand.Int(rand.Reader, alphabet_len)
		if err != nil {
			panic(err)
		}
		bytes[i] = alphabet[bigint.Int64()]
	}
	return string(bytes[:])
}
