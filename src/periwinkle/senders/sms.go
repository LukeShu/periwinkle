// Copyright 2015 Davis Webb

package senders

import(

	"net/http"
	"net/url"
	"os"
	"bytes"
	"fmt"

)

func Url_handler(w http.ResponseWriter, req *http.Request) {
    
	//body, err := ioutil.ReadAll(req.Body)		

	//if err != nil {
      //  fmt.Println("Error occurred!")
    //}

	

}







// TODO: checking for message status: create an url for retrieving message status

// Returns the status of the sent message.
//If successful, returns OK

func sender(message, sender, receiver string) string {

	// account SID for Twilio account
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")

	messages_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/Messages.json"

	v := url.Values{}
	v.Set("From", sender)
	v.Set("To", receiver)
	v.Set("Body", message)
	v.Set("StatusCallback", "url_for_status_update")

	client := &http.Client{}

	req, err := http.NewRequest("POST", messages_url, bytes.NewBuffer([]byte(v.Encode())))
	if err != nil {
		fmt.Printf("%s\n", err)
		return "Error"
	}
	req.SetBasicAuth(account_sid, auth_token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%s\n", err)
		return "Error"
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {

		return "OK"

	} else {

		return resp.Status

	}
}


