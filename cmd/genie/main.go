package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/nlopes/slack"

	"github.com/pestophagous/brevity-genie-slack-bot/pkg/brevity"
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/util"
)

func main() {
	brevityBot := brevity.NewBrevityBot()

	token := util.MustGetEnv("SLACKTOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		var msg slack.RTMEvent
		select {
		case msg = <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.ConnectedEvent:
				log.Print("ConnectedEvent")
			case *slack.ConnectionErrorEvent:
				log.Print("ConnectionErrorEvent")
			case *slack.DisconnectedEvent:
				log.Print("DisconnectedEvent")
			case *slack.MessageTooLongEvent:
				log.Print("MessageTooLongEvent")
			case *slack.OutgoingErrorEvent:
				log.Print("OutgoingErrorEvent")
			case *slack.IncomingEventError:
				log.Print("IncomingEventError")
			case *slack.UnmarshallingErrorEvent:
				log.Print("UnmarshallingErrorEvent")
			case *slack.HelloEvent:
				log.Print("HelloEvent")
			case *slack.RateLimitEvent:
				log.Print("RateLimitEvent")
			case *slack.AckErrorEvent:
				log.Print("AckErrorEvent")
			case *slack.LatencyReport:
				log.Print("LatencyReport")

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				text := ev.Text
				log.Print(ev.User)      // US8FL2R91
				log.Print(ev.Channel)   // CS1EW078Q
				log.Print(ev.Type)      // "message"
				log.Print(ev.Timestamp) // 1578794317.001700 // t=1578794317001700 ; date -d @${t%??????}
				log.Print(ev.Text)
				log.Print(ev.BotID)
				log.Print(ev.Username)

				text = strings.TrimSpace(text)
				text = strings.ToLower(text)

				matched, _ := regexp.MatchString("test poke the bot", text)

				if ev.User != info.User.ID && matched {
					rtm.SendMessage(rtm.NewOutgoingMessage("you got it", ev.Channel))
				}

				brevityBot.Track(brevity.NewChatActivity(util.MessageTimestampToUTC(ev.Timestamp), ev.User))

			case *slack.RTMError:
				log.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop

			default:
				log.Printf("Switch missing a case for: %s", msg.Type)
				// Take no action
			}
		}
	}
}
