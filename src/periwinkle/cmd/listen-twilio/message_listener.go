// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Davis Webb

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"periwinkle/cfg"
	"periwinkle/putil"
	"periwinkle/backend"
	"periwinkle/twilio"
	"strings"
	"time"
)

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s [CONFIG_FILE]\n", os.Args[0])
}

func main() {
	config_filename := "./periwinkle.yaml"
	switch len(os.Args) {
	case 1:
		// do nothing
	case 2:
		config_filename = os.Args[1]
	default:
		usage(os.Stderr)
		os.Exit(2)
	}

	file, err := os.Open(config_filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %q: %v\n", config_filename, err)
		os.Exit(1)
	}

	config, err := cfg.Parse(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %q: %v\n", config_filename, err)
		os.Exit(1)
	}

	var arr_temp [1000]string
	var cur_time_sec int64
	cur_time_sec = 0

	for {
		time.Sleep(time.Second)
		numbers := backend.GetAllUsedTwilioNumbers(config.DB)

		for _, number := range numbers {
			// clear the array
			if cur_time_sec != time.Now().UTC().Unix() {
				for j := 0; j != len(arr_temp); j++ {
					arr_temp[j] = ""
				}
			}

			cur_time := time.Now().UTC()
			cur_time_sec = cur_time.Unix()

			// gets url for received  Twilio messages for a given date
			url := "https://api.twilio.com/2010-04-01/Accounts/" + config.TwilioAccountId + "/Messages.json?To=" + number.Number + "&DateSent>=" + strings.Split(cur_time.String(), " ")[0]

			client := &http.Client{}

			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(config.TwilioAccountId, config.TwilioAuthToken)

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

			mes_len := len(message.Messages)

			if mes_len != 0 {
				for i := 0; i < mes_len; i++ {
					time_send, _ := time.Parse(time.RFC1123Z, message.Messages[i].DateSent)

					if err != nil {
						log.Println(err)
					}

					if time_send.Unix() >= cur_time.Unix() {
						m_sid := message.Messages[i].Sid

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
						for j := 0; j != len(arr_temp); j++ {
							if arr_temp[j] == "" {
								arr_temp[j] = m_sid

								user := backend.GetUserByAddress(config.DB, "sms", message.Messages[i].From)
								group := backend.GetGroupByUserAndTwilioNumber(config.DB, user.Id, message.Messages[i].To)
								//Not yet set: cfg.GroupDomain="periwinkle.lol"
								putil.MessageBuilder{
									Maildir: config.Mailstore,
									Headers: map[string]string{
										"To":      group.Id + "@" + config.GroupDomain,
										"From":    backend.GetAddressByIdAndMedium(config.DB, user.Id, "sms").AsEmailAddress(),
										"Subject": user.Id + "--> " + message.Messages[i].Body,
									},
									Body: "",
								}.Done()
								break
							} else if arr_temp[j] == m_sid {
								break
							}
						}
					}
				}
			}
		}
	}
}
