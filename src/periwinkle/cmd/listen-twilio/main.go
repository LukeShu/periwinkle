// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/cmdutil"
	"periwinkle/twilio"
	"strings"
	"time"
)

const usage = `
Usage: %[1]s [-c CONFIG_FILE]
       %[1]s -h | --help
Repeatedly poll Twilio for new messages.

Options:
  -h, --help      Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`

func main() {
	options := cmdutil.Docopt(usage)
	config := cmdutil.GetConfig(options["-c"].(string))

	var arrTemp [1000]string
	var curTimeSec int64
	curTimeSec = 0

	for {
		time.Sleep(time.Second)
		conflict := config.DB.Do(func(tx *periwinkle.Tx) {
			numbers := backend.GetAllUsedTwilioNumbers(tx)

			for _, number := range numbers {
				// clear the array
				if curTimeSec != time.Now().UTC().Unix() {
					for j := 0; j != len(arrTemp); j++ {
						arrTemp[j] = ""
					}
				}

				curTime := time.Now().UTC()
				curTimeSec = curTime.Unix()

				// gets url for received  Twilio messages for a given date
				url := "https://api.twilio.com/2010-04-01/Accounts/" + config.TwilioAccountID + "/Messages.json?To=" + number.Number + "&DateSent>=" + strings.Split(curTime.String(), " ")[0]

				client := &http.Client{}

				req, _ := http.NewRequest("GET", url, nil)
				req.SetBasicAuth(config.TwilioAccountID, config.TwilioAuthToken)

				resp, err := client.Do(req)

				if err != nil {
					log.Println(err)
				}

				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
				}

				// converts JSON messages
				message := twilio.Paging{}
				json.Unmarshal([]byte(body), &message)

				mesLen := len(message.Messages)

				if mesLen != 0 {
					for i := 0; i < mesLen; i++ {
						timeSend, _ := time.Parse(time.RFC1123Z, message.Messages[i].DateSent)

						if err != nil {
							log.Println(err)
						}

						if timeSend.Unix() >= curTime.Unix() {
							mSID := message.Messages[i].Sid

							// Since we only can get
							// recived Twilio messages for
							// a specific date, we need to
							// store messages received in
							// a second and clear them
							// once a second elapsed.
							//
							// In a second, one message
							// may appear multiple
							// times. So we want to avoid
							// duplicates.
							for j := 0; j != len(arrTemp); j++ {
								if arrTemp[j] == "" {
									arrTemp[j] = mSID

									user := backend.GetUserByAddress(tx, "sms", message.Messages[i].From)
									group := backend.GetGroupByUserAndTwilioNumber(tx, user.ID, message.Messages[i].To)
									//Not yet set: cfg.GroupDomain="periwinkle.lol"
									fmt.Println("GroupName:", group.ID)
									MessageBuilder{
										Maildir: config.Mailstore,
										Headers: map[string]string{
											"To":      group.ID + "@" + "periwinkle.lol", //config.GroupDomain,
											"From":    backend.GetAddressByUserAndMedium(tx, user.ID, "sms").AsEmailAddress(),
											"Subject": user.ID + "--> " + message.Messages[i].Body,
										},
										Body: "",
									}.Done()
									break
								} else if arrTemp[j] == mSID {
									break
								}
							}
						}
					}
				}
			}
		})
		if conflict != nil {
			periwinkle.LogErr(conflict)
		}
	}
}