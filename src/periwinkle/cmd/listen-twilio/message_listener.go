// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package main

import (
	"encoding/json"
	"io/ioutil"
	"locale"
	"log"
	"net/http"
	"os"
	"periwinkle"
	"periwinkle/backend"
	"periwinkle/cfg"
	"periwinkle/putil"
	"periwinkle/twilio"
	"strings"
	"time"

	"lukeshu.com/git/go/libsystemd.git/sd_daemon/lsb"
)

const usage = `
Usage: %[1]s [-c CONFIG_FILE]
       %[1]s -h | --help
Repeatedly poll Twilio for new messages.

Options:
  -h, --help      Display this message.
  -c CONFIG_FILE  Specify the configuration file [default: ./config.yaml].`

func main() {
	options := periwinkle.Docopt(usage)

	configFile, uerr := os.Open(options["-c"].(string))
	if uerr != nil {
		periwinkle.LogErr(locale.UntranslatedError(uerr))
		os.Exit(int(lsb.EXIT_NOTCONFIGURED))
	}

	config, err := cfg.Parse(configFile)
	if err != nil {
		periwinkle.LogErr(err)
		os.Exit(int(lsb.EXIT_NOTCONFIGURED))
	}

	var arrTemp [1000]string
	var curTimeSec int64
	curTimeSec = 0

	for {
		time.Sleep(time.Second)
		numbers := backend.GetAllUsedTwilioNumbers(config.DB)

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

								user := backend.GetUserByAddress(config.DB, "sms", message.Messages[i].From)
								group := backend.GetGroupByUserAndTwilioNumber(config.DB, user.ID, message.Messages[i].To)
								//Not yet set: cfg.GroupDomain="periwinkle.lol"
								putil.MessageBuilder{
									Maildir: config.Mailstore,
									Headers: map[string]string{
										"To":      group.ID + "@" + config.GroupDomain,
										"From":    backend.GetAddressByIDAndMedium(config.DB, user.ID, "sms").AsEmailAddress(),
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
	}
}
