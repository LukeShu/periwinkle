package test

import (
	"fmt"
	"log"
	"periwinkle"
	"periwinkle/backend"
	//"periwinkle/twilio"
	//"github.com/jinzhu/gorm"
)

func Test(cfg *periwinkle.Cfg, db *periwinkle.Tx) {

	num := backend.TwilioNumber{
		Number: "+13346038139",
	}

	if err := db.Create(&num).Error; err != nil {
		panic(err)
	}

	user1 := backend.User{
		ID:        "alex",
		FullName:  "",
		Addresses: []backend.UserAddress{{Medium: "email", Address: "zsuleime@purdue.edu", Confirmed: true}},
	}

	err := db.Create(&user1).Error
	if err != nil {
		log.Println(err)
	}

	user2 := backend.User{
		ID:        "john",
		FullName:  "",
		Addresses: []backend.UserAddress{{Medium: "sms", Address: "+17656027006", Confirmed: true}, {Medium: "email", Address: "s.jandos91@gmail.com", Confirmed: true}},
	}

	err = db.Create(&user2).Error
	if err != nil {
		log.Println(err)
	}

	user3 := backend.User{
		ID:        "guntas",
		FullName:  "",
		Addresses: []backend.UserAddress{{Medium: "sms", Address: "+16166342620", Confirmed: true}},
	}

	err = db.Create(&user3).Error
	if err != nil {
		log.Println(err)
	}
	existence := [2]int{2, 2}
	read := [2]int{2, 2}
	post := [3]int{1, 1, 1}
	join := [3]int{1, 1, 1}
	err = db.Create(&backend.Group{
		ID:                 "Purdue",
		ReadPublic:         read[0],
		ReadConfirmed:      read[1],
		ExistencePublic:    existence[0],
		ExistenceConfirmed: existence[1],
		PostPublic:         post[0],
		PostConfirmed:      post[1],
		PostMember:         post[2],
		JoinPublic:         join[0],
		JoinConfirmed:      join[1],
		JoinMember:         join[2],
		 Subscriptions: []backend.Subscription{
		 	{AddressID: user1.Addresses[0].ID, Confirmed: true},
		 	{AddressID: user2.Addresses[0].ID, Confirmed: true},
		 	{AddressID: user2.Addresses[1].ID, Confirmed: true},
		 	{AddressID: user3.Addresses[0].ID, Confirmed: true},
		 },
	}).Error
	if err != nil {
		log.Println(err)
	}

	fmt.Println()
}

//SUCCESSFULL SMS TEST
// ORIGINAL_RECIPIENT=+16166342620@sms.gateway bin/receive-email < <(printf '%s\r\n' 'To: +16166342620@sms.gateway' 'From: Purdue@periwinkle.lol' "Subject: sms testing" "Message-Id: $RANDOM@bar" '' 'body')
