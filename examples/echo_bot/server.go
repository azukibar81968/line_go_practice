// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"regexp"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	bot, err := linebot.New(
		"8cffe00334fcd57e073f7691cb773af3",
		"61967pgiDfx/Kl7j/pwigA+s9TCvC7NyDrDhg6M3l5WNKpwgTJAHZCPKIC5nYa38Bv1zI577HZYTUnV8SA1XBVAdSVR9DMzeRFusMOvTEOBizl3WOvuARK8rSEyd08gYqxJWzfJ5E+AwrMh4n+z5+AdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}
	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
				r := regexp.MustCompile(message.Text)
				if r.MatchString(`[0-9]`) == true{
					var num = ParseLeadingInt(message.Text)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("waiting "+strconv.Itoa(num)+"sec."), linebot.NewTextMessage("(^^)")).Do(); err != nil {
						log.Print(err)
					}
//					time.Sleep(time.Second * num)
					if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("finish!!")).Do(); err != nil {
						log.Print(err)
					}
				}

				if r.MatchString("aaa"){
					if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage(message_check(message.Text))).Do(); err != nil {
						log.Print(err)
					}
				} else {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("waiting '5'sec."), linebot.NewTextMessage("(^^)")).Do(); err != nil {
						log.Print(err)
					}
					time.Sleep(time.Second * 5)
					if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("It's time")).Do(); err != nil {
						log.Print(err)
					}
				}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}


func message_check(s string) string{
	r := regexp.MustCompile(s)
	if r.MatchString("??????"){
		return "??????"
	} else {
		return s
	}
}

var rexLeadingDigits = regexp.MustCompile(`\d+`)

func ParseLeadingInt(s string) int {
    rex := rexLeadingDigits.Copy()
    value, _ := strconv.Atoi(rex.FindString(s))
    return value
}