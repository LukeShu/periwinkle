
// Copyright 2015 Zhandos Suleimenov

package twilio_json



type Paging struct {


	First_page_uri string `json:"first_page_uri"`

	End int `json:"end"`

	Previous_page_uri string `json:"previous_page_uri"`

	Messages []Message `json:"messages"`

	Uri string `json:"uri"`

	Page_size int `json:"page_size"`
  
	Start int `json:"start"`

	Next_page_uri string `json:"next_page_uri"`
	
	Page int `json:"page"`


}




type Message struct {

	Sid string `json:"sid"`

	DateCreated string `json:"date_created"`

	DateUpdated string `json:"date_updated"`

	DateSent string `json:"date_sent"`

	AccountSid string `json:"account_sid"`	

	To string `json:"to"`

	From string `json:"from"`

	Body string `json:"body"`

	Status string `json:"status"`

	NumSegments int `json:"num_segments,string"`

	NumMedia int `json:"num_media,string"`

	Direction string  `json:"direction"`

	ApiVersion string `json:"api_version"`	

	Price string `json:"price"`

	PriceUnit string `json:"price_unit"`

	ErrorCode string `json:"error_code"`

	ErrorMessage string `json:"error_message"`

	Uri string `json:"uri"`	

	SubresourceUri []Media  `json:"subresource_uris"`



}




type Media struct {

	media string `json:"media"`

}


