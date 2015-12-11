// Copyright 2015 Luke Shumaker

package backend

// The list of tables, in the order that they need to be created in.
var tables = []table{
	Captcha{},
	Medium{},
	Group{},
	Message{}, // must come after Group
	User{},
	Session{}, // must come after User
	ShortURL{},
	UserAddress{},  // must come after User and Medium
	Subscription{}, // must come after Group and UserAddress
	TwilioNumber{},
	TwilioPool{}, // must come after TwilioNumber, User, and Group
	TwilioSMSMessage{},
	Admin{},
}
