
// Copyright 2015 Zhandos Suleimenov


package twilio



import(

	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"strings"
	"periwinkle/listeners/twilio/twilio_json"

)




func Main() error {


	// account SID for Twilio account 
	account_sid := "AC3140ba4160934bc0e4f874e2d99871ee"

	// Authorization token for Twilio account
	auth_token := "b7ac9496408440c4755087dc4a934a4d"


	var arr_temp [1000]string
	var cur_time_sec int64
	cur_time_sec = 0;




	for {


	// clear the array

	if cur_time_sec != time.Now().UTC().Unix() {

	for j:= 0; j != len(arr_temp); j++ {

	arr_temp[j] = ""

	}
	}
 	

	cur_time := time.Now().UTC()
	cur_time_sec = cur_time.Unix()


	// gets url for received  Twilio messages for a given date 

	url := "https://api.twilio.com/2010-04-01/Accounts/" + account_sid + "/Messages.json?To=+17653569541&DateSent>=" + strings.Split(cur_time.String() , " ")[0]

	

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil )
	req.SetBasicAuth(account_sid, auth_token)
	
	resp, err := client.Do(req)
	
	if err != nil {
		fmt.Print("%s", "error");
	}
	

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("%s", err);
	}

	
	//converts JSON messages
	
	message := twilio_json.Paging{}
	json.Unmarshal([]byte(body), &message)
	
	mes_len := len(message.Messages)

	if mes_len != 0 {
	
	for i := 0; i<mes_len; i++ {


	time_send,_ := time.Parse(time.RFC1123Z, message.Messages[i].DateSent)
	
	if err != nil {
	fmt.Print("%s", err)
	}

	

	if time_send.Unix() >= cur_time.Unix() {

	m_sid := message.Messages[i].Sid

	// Since we only can get recived Twilio messages for a specific date, we need to store messages received in a second and clear them once a second elapsed. 
	//In a second, one message may appear multiple times. So we want to avoid duplicates.  	

	for j:= 0; j != len(arr_temp); j++ {
	
	if(arr_temp[j] == "" ) {
	arr_temp[j] = m_sid

	//TODO adding the message to sqrl

	fmt.Println(message.Messages[i].Body)
	break

	}else if(arr_temp[j] == m_sid) {
	break}	


	}
	}

	
	}
	}
	}
	


}
