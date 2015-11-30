// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package domain_handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/putil"
	"postfixpipe"

	"github.com/jinzhu/gorm"
)

func HandleEmail(r io.Reader, name string, db *gorm.DB, cfg *periwinkle.Cfg) postfixpipe.ExitStatus {
	mdWriter := cfg.Mailstore.NewMail()
	if mdWriter == nil {
		periwinkle.Logf("Could not open maildir for writing: %q\n", cfg.Mailstore)
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

	group := backend.GetGroupByID(db, name)
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
					if perror.HTTPCode() == 409 {
						silentbounce = true
					}
				}
				panic(obj)
			}
		}()
		backend.NewMessage(
			db,
			msg.Header.Get("Message-Id"),
			*group,
			mdWriter.Unique())
	}()
	if silentbounce {
		return postfixpipe.EX_OK
	}
	mdWriter.Close()
	mdWriter = nil

	// collect IDs of addresses subscribed to the group
	addressIDs := make([]int64, len(group.Subscriptions))
	for i := range group.Subscriptions {
		addressIDs[i] = group.Subscriptions[i].AddressID
	}

	// fetch all of those addresses
	var addressList []backend.UserAddress
	if len(addressIDs) > 0 {
		db.Where("id IN (?)", addressIDs).Find(&addressList)
	} else {
		addressList = make([]backend.UserAddress, 0)
	}

	// convert that list into a set
	forwardSet := make(map[string]bool, len(addressList))
	for _, addr := range addressList {
		forwardSet[addr.AsEmailAddress()] = true
	}

	// prune addresses that (should) already have the message
	for _, header := range []string{"To", "From", "Cc"} {
		addresses, err := msg.Header.AddressList(header)
		if err != nil {
			periwinkle.Logf("Parsing %q Header: %v\n", header, err)
		}
		for _, addr := range addresses {
			delete(forwardSet, addr.Address)
		}
	}

	// convert the set into an array
	forwardAry := make([]string, len(forwardSet))
	i := uint(0)
	for addr := range forwardSet {
		forwardAry[i] = addr
		i++
	}

	// format the message
	msg822 := []byte{}
	for k := range msg.Header {
		msg822 = append(msg822, []byte(fmt.Sprintf("%s: %s\r\n", k, msg.Header.Get(k)))...)
	}
	msg822 = append(msg822, []byte("\r\n")...)
	body, _ := ioutil.ReadAll(msg.Body) // TODO: error handling
	msg822 = append(msg822, body...)

	if len(forwardAry) > 0 {
		// send the message out
		err = smtp.SendMail("localhost:25",
			smtp.PlainAuth("", "", "", ""),
			msg.Header.Get("From"),
			forwardAry,
			msg822)
		if err != nil {
			periwinkle.Logf("Error sending: %v", err)
			return postfixpipe.EX_UNAVAILABLE
		}
	}
	return postfixpipe.EX_OK
}
