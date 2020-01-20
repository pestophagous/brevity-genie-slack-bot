package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/nlopes/slack"

	"github.com/pestophagous/brevity-genie-slack-bot/pkg/brevity"
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/trace"
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/util"
)

var approvedMissingEventCases []string = []string{"ack", "user_typing"}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Lshortfile for file.go:NN

	// trace.EnableTraces([]string{"incoming-event"})
	trace.EnableTraces([]string{"user-metadata"})
	trace.EnableTraces([]string{"package-brevity"})
}

func main() {
	token := util.MustGetEnv("SLACKTOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()

	brevityBot := brevity.NewBrevityBot(&slackAdapter{client: api})

	go rtm.ManageConnection()

Loop:
	for {
		var msg slack.RTMEvent
		select {
		case msg = <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.ConnectedEvent:
				trace.Trace("incoming-event", "ConnectedEvent")
				log.Print("ConnectedEvent")
			case *slack.ConnectionErrorEvent:
				trace.Trace("incoming-event", "ConnectionErrorEvent")
			case *slack.DisconnectedEvent:
				trace.Trace("incoming-event", "DisconnectedEvent")
			case *slack.MessageTooLongEvent:
				trace.Trace("incoming-event", "MessageTooLongEvent")
			case *slack.OutgoingErrorEvent:
				trace.Trace("incoming-event", "OutgoingErrorEvent")
			case *slack.IncomingEventError:
				trace.Trace("incoming-event", "IncomingEventError")
			case *slack.UnmarshallingErrorEvent:
				trace.Trace("incoming-event", "UnmarshallingErrorEvent")
			case *slack.HelloEvent:
				trace.Trace("incoming-event", "HelloEvent")
				log.Print("HelloEvent")
			case *slack.RateLimitEvent:
				trace.Trace("incoming-event", "RateLimitEvent")
			case *slack.AckErrorEvent:
				trace.Trace("incoming-event", "AckErrorEvent")
			case *slack.LatencyReport:
				trace.Trace("incoming-event", "LatencyReport")

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				text := ev.Text
				trace.Trace("incoming-event", ev.User)      // US8FL2R91
				trace.Trace("incoming-event", ev.Channel)   // CS1EW078Q
				trace.Trace("incoming-event", ev.Type)      // "message"
				trace.Trace("incoming-event", ev.Timestamp) // 1578794317.001700 // t=1578794317001700 ; date -d @${t%??????}
				trace.Trace("incoming-event", ev.Text)
				trace.Trace("incoming-event", ev.BotID)
				trace.Trace("incoming-event", ev.Username)

				text = strings.TrimSpace(text)
				text = strings.ToLower(text)

				matched, _ := regexp.MatchString("test poke the bot", text)

				if ev.User != info.User.ID && matched {
					rtm.SendMessage(rtm.NewOutgoingMessage("you got it", ev.Channel))
				}

				excusals := brevityBot.Track(brevity.NewChatActivity(util.MessageTimestampToUTC(ev.Timestamp), ev.User))

				for _, excuse := range excusals {
					rtm.SendMessage(rtm.NewOutgoingMessage(excuse, ev.Channel))
				}

			case *slack.RTMError:
				log.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop

			default:
				if !util.MatchStringInSlice(msg.Type, approvedMissingEventCases) {
					log.Printf("Switch missing a case for: %s", msg.Type)
				}
			}
		}
	}
}
