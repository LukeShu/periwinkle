// Copyright 2015 Zhandos Suleimenov

package twilio

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"	
	"bytes"
	"os"
)


	// function returns  a phone number and Status
	//if successful, returns a new phone number and OK	

func newPhoneNum() (string, string) {


	// account SID for Twilio account 
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")



	// gets url for available numbers 

	availNum_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/AvailablePhoneNumbers/US/Local.json?SmsEnabled=true&MmsEnabled=true"  

	// gets url for a new phone number

	newPhoneNum_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/IncomingPhoneNumbers.json"  

	client := &http.Client{}

	req, err := http.NewRequest("GET", availNum_url, nil )

	if(err != nil) {
		fmt.Printf("%s\n", err)
		return "", "Error"
	}

	req.SetBasicAuth(account_sid, auth_token)
	
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Print("%s", "error");
		return "", "Error"
	}

	if resp.StatusCode == 302 {

		url, err := resp.Location()
		if(err != nil) {
		fmt.Println(err)
		return "", "Error"
		}

		req, err = http.NewRequest("GET", url.String(), nil )
			if(err != nil) {
			fmt.Printf("%s\n", err)
			return "", "Error"
			}

		req.SetBasicAuth(account_sid, auth_token)
		resp, err = client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			fmt.Print("%s", "error");
			return "", "Error"
			}

		if resp.StatusCode != 200 {

		return "", resp.Status
		}

	}else if resp.StatusCode == 200 {

	//continue

	}else{
		return "", resp.Status 	
	}

	
	



	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("%s", err);
		return "", "Error"
	}


	avail_number := Avail_ph_num{}
	json.Unmarshal(body, &avail_number )	

	



	if len(avail_number.PhoneNumberList) != 0 {

		number := avail_number .PhoneNumberList[0].PhoneNumber;

    	val := url.Values{}
		val.Set("PhoneNumber", avail_number .PhoneNumberList[0].PhoneNumber)
		val.Set("SmsUrl", "http://twimlets.com/echo?Twiml=%3CResponse%3E%3C%2FResponse%3E")

		req, err = http.NewRequest("POST", newPhoneNum_url, bytes.NewBuffer([]byte(val.Encode())) )
		if(err != nil) {
		fmt.Printf("%s\n", err)
		return "", "Error"
		}

		req.SetBasicAuth(account_sid, auth_token)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err = client.Do(req)	
		defer resp.Body.Close()
		if(err != nil) {
			fmt.Printf("%s\n", err)
			return "", "Error"
			}

		if resp.StatusCode != 200 && resp.StatusCode != 201  {
			return "", resp.Status;
		}

		

		body, err = ioutil.ReadAll(resp.Body)			
		if err != nil {
			fmt.Print("%s", err);
			return "", ""
			}
 		fmt.Println(string(body))
		return number, "OK"

	}

	fmt.Println("There are no available phone numbers!!!")
	return "", "Error"
	
}





	// TODO: checking for message status: create an url for retrieving message status

	// Returns the status of the sent message. 
	//If successful, returns OK

func sender(message, sender, receiver string ) string {

	// account SID for Twilio account 
	account_sid := os.Getenv("TWILIO_ACCOUNTID")

	// Authorization token for Twilio account
	auth_token := os.Getenv("TWILIO_TOKEN")



	messages_url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/Messages.json"

	
	v := url.Values{}
	v.Set("From", sender )
	v.Set("To", receiver)
	v.Set("Body", message)	
	v.Set("StatusCallback", "url_for_status_update" )

	client := &http.Client{}

	req, err := http.NewRequest("POST", messages_url, bytes.NewBuffer([]byte(v.Encode())) )
	if(err != nil) {
		fmt.Printf("%s\n", err)
		return "Error"
	}
	req.SetBasicAuth(account_sid, auth_token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

   
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if(err != nil) {
		fmt.Printf("%s\n", err)
		return "Error"
	}


	if resp.StatusCode == 200 || resp.StatusCode == 201 {

	return "OK"


	}else {

	return resp.Status

	}
}




