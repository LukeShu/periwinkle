package test

import (
	//"fmt"
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
		ID:        "Alex",
		FullName:  "",
		Addresses: []backend.UserAddress{{Medium: "email", Address: "zsuleime@purdue.edu", Confirmed: true}},
	}

	err := db.Create(&user1).Error
	if err != nil {
		log.Println(err)
	}

	user2 := backend.User{
		ID:        "John",
		FullName:  "",
		Addresses: []backend.UserAddress{{Medium: "sms", Address: "+17656027006", Confirmed: true}, {Medium: "email", Address: "s.jandos91@gmail.com", Confirmed: true}},
	}

	err = db.Create(&user2).Error
	if err != nil {
		log.Println(err)
	}

	user3 := backend.User{
		ID:        "Ben",
		FullName:  "",
		Addresses: []backend.UserAddress{{Medium: "email", Address: "s_jandos@mail.ru", Confirmed: true}},
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

}
