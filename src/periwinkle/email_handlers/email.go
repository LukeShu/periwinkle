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
	"postfixpipe"
)

func HandleEmail(r io.Reader, name string, db *gorm.DB) (ret uint8) {
	mdWriter := cfg.Mailstore.NewMail()
	r = io.TeeReader(r, mdWriter)
	msg, err := mail.ReadMessage(r)
	if err != nil {
		mdWriter.Cancel()
		ret = postfixpipe.EX_NOINPUT
		return
	}

	group := store.GetGroupById(db, name)
	if group == nil {
		mdWriter.Cancel()
		ret = postfixpipe.EX_NOUSER
		return
	}

	// collect IDs of addresses subscribed to the group
	address_ids := make([]int64, len(group.Subscriptions))
	for i := range group.Subscriptions {
		address_ids[i] = group.Subscriptions[i].AddressId
	}

	// fetch all of those addresses
	var address_list []store.UserAddress
	db.Where("id in (?)", address_ids).Find(&address_list)

	// convert that list into a set
	forward_set := make(map[string]bool, len(address_list))
	for _, addr := range address_list {
		forward_set[addr.AsEmailAddress()] = true
	}

	// prune addresses that (should) already have the message
	for _, header := range []string{"To", "From", "Cc"} {
		addresses, err := msg.Header.AddressList(header)
		if err != nil {
			log.Fatalf("Parsing %q Header: %v\n", header, err)
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
		panic(err)
	}
	ret = postfixpipe.EX_OK
	return
}
