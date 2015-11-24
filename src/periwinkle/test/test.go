package test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"periwinkle/backend"
)

func Test(db *gorm.DB) {

	err := db.Create(&backend.User{
		ID:        "Alex",
		FullName:  "",
		Addresses: []backend.UserAddress{{Medium: "email", Address: "zsuleime@purdue.edu", Confirmed: true}},
	}).Error

	if err != nil {
		log.Println(err)
	}

	err = db.Create(&backend.User{
		ID:        "John",
		FullName:  "",
		Addresses: []backend.UserAddress{{Medium: "sms", Address: "+17656027006", Confirmed: true}, {Medium: "email", Address: "s.jandos91@gmail.com", Confirmed: true}},
	}).Error

	if err != nil {
		log.Println(err)
	}

	benAddr := []backend.UserAddress{
		{Medium: "email",
			Address:   "s_jandos@mail.ru",
			Confirmed: true},
	}
	ben := backend.User{
		ID:        "Ben",
		FullName:  "",
		Addresses: benAddr,
	}

	err = db.Create(&ben).Error

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
			Address:   backend.UserAddress{UserID: "Alex", Medium: "email", Address: "zsuleime@purdue.edu", Confirmed: true},
			Confirmed: true,
		},

			{Address: benAddr[0],
				Confirmed: true,
			},

			{Address: backend.UserAddress{UserID: "John", Medium: "sms", Address: "+17656027006", Confirmed: true},
				Confirmed: true,
			},

			{Address: backend.UserAddress{UserID: "John", Medium: "email", Address: "s.jandos91@gmail.com", Confirmed: true},
				Confirmed: true,
			},
		},
	}).Error

	if err != nil {
		log.Println(err)
	}

	fmt.Println("All existing twilio numbers: ", backend.GetAllExistingTwilioNumbers())
	fmt.Println("All unused numbers for John", backend.GetUnusedTwilioNumbersByUser(db, "John"))
	//backend.AssignTwilioNumber(db, "John", "Purdue", backend.GetUnusedTwilioNumbersByUser(db, "John")[0])
	//fmt.Println("All unused numbers for John", backend.GetUnusedTwilioNumbersByUser(db, "John"))

}
