// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package handlers

import (
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"periwinkle/store"
)

func HandleEmail(r io.Reader, name string, db *gorm.DB) int {

	msg, err := mail.ReadMessage(r)

	if err != nil {
		panic(err)
	}
	header := msg.Header

	group := store.GetGroupById(db, name)
	if group == nil {
		panic("No group: " + name)
	}
	// TODO: check if group == nil
	address_ids := make([]int64, len(group.Subscriptions))
	for i := range group.Subscriptions {
		address_ids[i] = group.Subscriptions[i].AddressId
	}
	var address_list []store.UserAddress
	db.Where("id in (?)", address_ids).Find(&address_list)
	// TODO: error handling
	// convert forward_ary into a set
	forward_set := make(map[string]bool, len(address_list))
	for _, addr := range address_list {
		var str string
		if addr.Medium == "email" {
			str = addr.Address
		} else {
			str = addr.Address + "@" + addr.Medium + ".gateway"
		}
		forward_set[str] = true
	}
	/////////////////////////////////////////////////////////////////////
	addresses, err := mail.ParseAddressList(header.Get("To"))
	if err != nil {
		panic(err)
	}
	for _, addr := range addresses {
		delete(forward_set, addr.Address)
	}
	/////////////////////////////////////////////////////////////////////
	addresses, err = mail.ParseAddressList(header.Get("From"))
	if err != nil {
		panic(err)
	}
	for _, addr := range addresses {
		delete(forward_set, addr.Address)
	}
	/////////////////////////////////////////////////////////////////////
	addresses, err = mail.ParseAddressList(header.Get("Cc"))
	if err != nil {
		panic(err)
	}
	for _, addr := range addresses {
		delete(forward_set, addr.Address)
	}

	forward_ary := make([]string, len(forward_set))
	i := uint(0)
	for addr := range forward_set {
		forward_ary[i] = addr
		i++
	}


	auth := smtp.PlainAuth("", "", "", "")
	from := header.Get("From")
	body, _ := ioutil.ReadAll(msg.Body)
	err = smtp.SendMail("localhost:25", auth, from, forward_ary, body)
	if err != nil {
		panic(err)
	}
	return 0
}
