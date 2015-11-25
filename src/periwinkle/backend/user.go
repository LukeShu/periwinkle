// Copyright 2015 Davis Webb
// Copyright 2015 Guntas Grewal
// Copyright 2015 Luke Shumaker

package backend

import (
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string        `json:"user_id"`
	FullName  string        `json:"fullname"`
	PwHash    []byte        `json:"-"`
	Addresses []UserAddress `json:"addresses"`
}

func (o User) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).Error
}

type UserAddress struct {
	// TODO: add a "verified" boolean
	ID            int64          `json:"-"`
	UserID        string         `json:"-"`
	Medium        string         `json:"medium"`
	Address       string         `json:"address"`
	SortOrder     uint64         `json:"sort_order"`
	Confirmed     bool           `json:"confirmed"`
	Subscriptions []Subscription `json:"subscriptions"`
}

func (o UserAddress) dbSchema(db *gorm.DB) error {
	return db.CreateTable(&o).
		AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT").
		AddForeignKey("medium", "media(id)", "RESTRICT", "RESTRICT").
		AddUniqueIndex("address_idx", "medium", "address").
		AddUniqueIndex("user_idx", "user_id", "sort_order").
		Error
}

func (addr UserAddress) AsEmailAddress() string {
	if addr.Medium == "email" {
		return addr.Address
	} else {
		return addr.Address + "@" + addr.Medium + ".gateway"
	}
}

func (u *User) populate(db *gorm.DB) {
	db.Model(u).Related(&u.Addresses)
	addressIDs := make([]int64, len(u.Addresses))
	for i, address := range u.Addresses {
		addressIDs[i] = address.ID
	}
	var subscriptions []Subscription
	if len(addressIDs) > 0 {
		if result := db.Where("address_id IN (?)", addressIDs).Find(&subscriptions); result.Error != nil {
			if !result.RecordNotFound() {
				panic(result.Error)
			}
		}
	} else {
		subscriptions = make([]Subscription, 0)
	}
	for i := range u.Addresses {
		u.Addresses[i].Subscriptions = []Subscription{}
		for _, subscription := range subscriptions {
			if u.Addresses[i].ID == subscription.AddressID {
				u.Addresses[i].Subscriptions = append(u.Addresses[i].Subscriptions, subscription)
			}
		}
	}
	var addresses []UserAddress

	for i, address := range u.Addresses {
		if address.Medium == "noop" {
			addresses = append(u.Addresses[:i], u.Addresses[i+1:]...)
			break
		}
	}
	for i, address := range addresses {
		if address.Medium == "admin" {
			addresses = append(addresses[:i], addresses[i+1:]...)
			break
		}
	}
	u.Addresses = addresses
}

func (u *User) GetUserSubscriptions(db *gorm.DB) []Subscription {
	db.Model(u).Related(&u.Addresses)
	addressIDs := make([]int64, len(u.Addresses))
	for i, address := range u.Addresses {
		addressIDs[i] = address.ID
	}
	var subscriptions []Subscription
	if len(addressIDs) > 0 {
		if result := db.Where("address_id IN (?)", addressIDs).Find(&subscriptions); result.Error != nil {
			if !result.RecordNotFound() {
				panic(result.Error)
			}
		}
	} else {
		subscriptions = make([]Subscription, 0)
	}
	return subscriptions
}

func GetAddressByIDAndMedium(db *gorm.DB, id string, medium string) *UserAddress {
	id = strings.ToLower(id)
	var o UserAddress
	if result := db.Where(&UserAddress{UserID: id, Medium: medium}).First(&o); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	return &o
}

func GetUserByID(db *gorm.DB, id string) *User {
	id = strings.ToLower(id)
	var o User
	if result := db.First(&o, "id = ?", id); result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	o.populate(db)
	return &o
}

func GetUserByAddress(db *gorm.DB, medium string, address string) *User {
	var o User
	result := db.Joins("inner join user_addresses on user_addresses.user_id=users.id").Where("user_addresses.medium=? and user_addresses.address=?", medium, address).Find(&o)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil
		}
		panic(result.Error)
	}
	o.populate(db)
	return &o
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	u.PwHash = hash
	return err
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PwHash, []byte(password))
	return err == nil
}

func NewUser(db *gorm.DB, name string, password string, email string) User {
	if name == "" {
		panic("name can't be empty")
	}
	o := User{
		ID:        name,
		FullName:  "",
		Addresses: []UserAddress{{Medium: "email", Address: email, Confirmed: false}},
	}
	if err := o.SetPassword(password); err != nil {
		panic(err)
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return o
}

func NewUserAddress(db *gorm.DB, userID string, medium string, address string, confirmed bool) UserAddress {
	o := UserAddress{
		UserID:        userID,
		Medium:        medium,
		Address:       address,
		Subscriptions: make([]Subscription, 0),
		Confirmed:     confirmed,
	}
	if err := db.Create(&o).Error; err != nil {
		panic(err)
	}
	return o
}

func (o *User) Save(db *gorm.DB) {
	if err := db.Save(o).Error; err != nil {
		panic(err)
	}
}
