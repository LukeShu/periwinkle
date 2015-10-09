// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package store

import (
	"crypto/rand"
	"encoding/json"
	he "httpentity"
	"io"
	"math/big"
	"periwinkle/cfg"
)

func init() {
	cfg.DB.AddTable(Captcha{})
	cfg.DB.AddTable(GroupAddress{})
	cfg.DB.AddTable(Group{})
	cfg.DB.AddTable(Medium{})
	cfg.DB.AddTable(Message{})
	cfg.DB.AddTable(Session{})
	cfg.DB.AddTable(ShortUrl{})
	cfg.DB.AddTable(Subscription{})
	cfg.DB.AddTable(UserAddress{})
	cfg.DB.AddTable(User{})
	if err := cfg.DB.CreateTablesIfNotExists(); err != nil {
		panic(err)
	}
}

// Simple dump to JSON, good for most entities
func defaultEncoders(o interface{}) map[string]he.Encoder {
	return map[string]he.Encoder{
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
