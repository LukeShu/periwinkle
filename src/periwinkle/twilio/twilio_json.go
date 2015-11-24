// Copyright 2015 Zhandos Suleimenov

package twilio

type Paging struct {
	FirstPageURI    string    `json:"first_page_uri"`
	End             int       `json:"end"`
	PreviousPageURI string    `json:"previous_page_uri"`
	Messages        []Message `json:"messages"`
	URI             string    `json:"uri"`
	PageSize        int       `json:"page_size"`
	Start           int       `json:"start"`
	NextPageURI     string    `json:"next_page_uri"`
	Page            int       `json:"page"`
}

type Message struct {
	Sid            string  `json:"sid"`
	DateCreated    string  `json:"date_created"`
	DateUpdated    string  `json:"date_updated"`
	DateSent       string  `json:"date_sent"`
	AccountSid     string  `json:"account_sid"`
	To             string  `json:"to"`
	From           string  `json:"from"`
	Body           string  `json:"body"`
	Status         string  `json:"status"`
	NumSegments    int     `json:"num_segments,string"`
	NumMedia       int     `json:"num_media,string"`
	Direction      string  `json:"direction"`
	ApiVersion     string  `json:"api_version"`
	Price          string  `json:"price"`
	PriceUnit      string  `json:"price_unit"`
	ErrorCode      string  `json:"error_code"`
	ErrorMessage   string  `json:"error_message"`
	URI            string  `json:"uri"`
	SubresourceURI []Media `json:"subresource_uris"`
}

type Media struct {
	Media string `json:"media"`
}

type AvailPhNum struct {
	PhoneNumberList []AvailPhoneNumbers `json:"available_phone_numbers"`
	URI             string              `json:"uri"`
}

type AvailPhoneNumbers struct {
	FriendlyName        string     `json:"friendly_name"`
	PhoneNumber         string     `json:"phone_number"`
	Lata                string     `json:"lata"`
	RateCenter          string     `json:"rate_center"`
	Latitude            string     `json:"latitude"`
	Longitude           string     `json:"longitude"`
	Region              string     `json:"region"`
	PostalCode          string     `json:"postal_code"`
	IsoCountry          string     `json:"iso_country"`
	AddressRequirements string     `json:"address_requirements"`
	Beta                bool       `json:"beta"`
	Capabilities        Capability `json:"capabilities"`
}

type Capability struct {
	Voice bool `json:"voice"`
	SMS   bool `json:"SMS"`
	MMS   bool `json:"MMS"`
}
