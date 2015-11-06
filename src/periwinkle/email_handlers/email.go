// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package handlers

import (
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"periwinkle/cfg"
	"periwinkle/store"
	"periwinkle/util" // putil
	"postfixpipe"
)

func HandleEmail(r io.Reader, name string, db *gorm.DB) uint8 {
	mdWriter := cfg.Mailstore.NewMail()
	if mdWriter == nil {
		log.Printf("Could not open maildir for writing: %s\n", cfg.Mailstore)
		return postfixpipe.EX_IOERR
	}
	defer func() {
		if mdWriter != nil {
			mdWriter.Cancel()
		}
	}()
	r = io.TeeReader(r, mdWriter)
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return postfixpipe.EX_NOINPUT
	}

	group := store.GetGroupById(db, name)
	if group == nil {
		return postfixpipe.EX_NOUSER
	}

	// try/catch looks awefully funny in Go
	silentbounce := false
	func() {
		defer func() {
			if obj := recover(); obj != nil {
				if err, ok := obj.(error); ok {
					perror := putil.ErrorToError(err)
					if perror.HttpCode() == 409 {
						silentbounce = true
					}
				}
				panic(obj)
			}
		}()
		store.NewMessage(
			db,
			msg.Header.Get("Message-Id"),
			*group,
			mdWriter.Unique())
	}()
	if silentbounce {
		return 0
	}
	mdWriter.Close()
	mdWriter = nil

	// collect IDs of addresses subscribed to the group
	address_ids := make([]int64, len(group.Subscriptions))
	for i := range group.Subscriptions {
		address_ids[i] = group.Subscriptions[i].AddressId
	}

	// fetch all of those addresses
	var address_list []store.UserAddress
	if len(address_ids) > 0 {
		db.Where("id IN (?)", address_ids).Find(&address_list)
	} else {
		address_list = make([]store.UserAddress, 0)
	}

	// convert that list into a set
	forward_set := make(map[string]bool, len(address_list))
	for _, addr := range address_list {
		forward_set[addr.AsEmailAddress()] = true
	}

	// prune addresses that (should) already have the message
	for _, header := range []string{"To", "From", "Cc"} {
		addresses, err := msg.Header.AddressList(header)
		if err != nil {
			log.Printf("Parsing %q Header: %v\n", header, err)
		}
		for _, addr := range addresses {
			delete(forward_set, addr.Address)
		}
	}

	// convert the set into an array
	forward_ary := make([]string, len(forward_set))
	i := uint(0)
	for addr := range forward_set {
		forward_ary[i] = addr
		i++
	}

	// send the message out
	body, _ := ioutil.ReadAll(msg.Body)
	err = smtp.SendMail("localhost:25",
		smtp.PlainAuth("", "", "", ""),
		msg.Header.Get("From"),
		forward_ary,
		body)
	if err != nil {
		log.Println("Error sending:", err)
		return postfixpipe.EX_UNAVAILABLE
	}
	return postfixpipe.EX_OK
}
