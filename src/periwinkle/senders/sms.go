// Copyright 2015 Davis Webb

package senders

import(

	"net/http"
	"net/url"
	"os"
	"bytes"
	"fmt"
	"io"
	"encoding/json"
)

var status string

func Url_handler(w http.ResponseWriter, req *http.Request) {
    
	//body, err := ioutil.ReadAll(req.Body)		

	//if err != nil {
		//fmt.Printf("%v", err)
	//}

	//converts JSON messages
	//message := Message{}
	//json.Unmarshal([]byte(body), &message)
	//status = Message.MessageStatus	

}







// TODO: checking for message status: create an url for retrieving message status

// Returns the status of the sent message.
//If successful, returns OK

func sender(reader io.Reader) string {

	status = ""

	buf := make([]byte, 1000)

	b :=  bytes.NewBuffer(buf)
	_, err := b.ReadFrom(reader)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	message := make(map[string]string)

	json.Unmarshal(buf, &message)
		

	// account SID for Twilio account
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")

	messages_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/Messages.json"

	v := url.Values{}
	v.Set("From", message["From"])
	v.Set("To", message["To"])
	v.Set("Body", message["Body"])
	v.Set("StatusCallback", "https://github.com/LukeShu/periwinkle/webui/twilio/sms")

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

		//return "OK"
		return status

	} else {

		return resp.Status

	}
}


