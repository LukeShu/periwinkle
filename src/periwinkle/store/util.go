// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package store

import (
	"crypto/rand"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io"
	"math/big"
)

func Schema(db *gorm.DB) {
	(Captcha{}).schema(db)
	(GroupAddress{}).schema(db)
	(Group{}).schema(db)
	(Medium{}).schema(db)
	(Message{}).schema(db)
	(Session{}).schema(db)
	(ShortUrl{}).schema(db)
	(Subscription{}).schema(db)
	(UserAddress{}).schema(db)
	(User{}).schema(db)
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

var alphabetLen = big.NewInt(int64(len(alphabet)))

func randomString(size int) string {
	var randStr []byte
	for i := 0; i < size; i++ {
		bigint, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			panic(err)
		}
		randStr[i] = alphabet[bigint.Int64()]
	}
	return string(randStr[:])
}
