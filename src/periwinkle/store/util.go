// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package store

import (
	"crypto/rand"
	"encoding/json"
	"github.com/jinzhu/gorm"
	he "httpentity"
	"httpentity/util"
	"io"
	"jsondiff"
	"math/big"
	"periwinkle/util" // putil
)

func Schema(db *gorm.DB) {
	// TODO: error detection
	(Captcha{}).schema(db)
	(Medium{}).schema(db)
	(Group{}).schema(db)
	(GroupAddress{}).schema(db) // must come after Group and Medium
	(Message{}).schema(db)      // must come after Group
	(User{}).schema(db)
	(Session{}).schema(db) // must come after User
	(ShortUrl{}).schema(db)
	(UserAddress{}).schema(db)  // must come after User and Medium
	(Subscription{}).schema(db) // must come after Group and UserAddress
}

func SchemaDrop(db *gorm.DB) {
	// This must be in the reverse order of Schema()
	// TODO: error detection
	db.DropTable(&Subscription{})
	db.DropTable(&UserAddress{})
	db.DropTable(&ShortUrl{})
	db.DropTable(&Session{})
	db.DropTable(&User{})
	db.DropTable(&Message{})
	db.DropTable(&GroupAddress{})
	db.DropTable(&Group{})
	db.DropTable(&Medium{})
	db.DropTable(&Captcha{})
}

func safeDecodeJSON(in interface{}, out interface{}) putil.HTTPError {
	decoder, ok := in.(*json.Decoder)
	if !ok {
		return putil.HTTPErrorf(415, "PUT and POST requests must have a document media type")
	}
	var tmp interface{}
	err := decoder.Decode(&tmp)
	if err != nil {
		return putil.HTTPErrorf(415, "Couldn't parse: %v", err)
	}
	str, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, out)
	if err != nil {
		return putil.HTTPErrorf(415, "Request body didn't have expected structure (field had wrong data type): %v", err)
	}
	if !jsondiff.Equal(tmp, out) {
		diff, err := jsondiff.NewJSONPatch(tmp, out)
		if err != nil {
			panic(err)
		}
		entity := heutil.NetMap{
			"message": "Request body didn't have expected structure (extra or missing fields).  The included diff would make the request acceptable.",
			"diff":    diff,
		}
		return putil.RawHTTPError(he.StatusUnsupportedMediaType(entity))
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
