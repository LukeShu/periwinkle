package test

import (
	"fmt"
	"log"
	"periwinkle"
	"periwinkle/backend"
	"github.com/jinzhu/gorm"
)

func Test(cfg *periwinkle.Cfg, db *gorm.DB) {

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

	err = db.Create(&backend.Group{
		ID:        "Purdue",
		Existence: 1,
		Read:      1,
		Post:      1,
		Join:      1,
		Subscriptions: []backend.Subscription{{
			Address:   user1.Addresses[0],
			Confirmed: true,
		},

			{Address: user2.Addresses[0],
				Confirmed: true,
			},

			{Address: user2.Addresses[1],
				Confirmed: true,
			},

			{Address: user3.Addresses[0],
				Confirmed: true,
			},
		},
	}).Error

	if err != nil {
		log.Println(err)
	}

	backend.AssignTwilioNumber(db, "guntas", "Purdue", "+13346038139")	
	gr := backend.GetGroupByUserAndTwilioNumber(db, "guntas", "+13346038139") 
	fmt.Println(gr.ID)
}
//SUCCESSFULL SMS TEST
// ORIGINAL_RECIPIENT=+16166342620@sms.gateway bin/receive-email < <(printf '%s\r\n' 'To: +16166342620@sms.gateway' 'From: Purdue@periwinkle.lol' "Subject: email testing" "Message-Id: $RANDOM@bar" '' 'body')