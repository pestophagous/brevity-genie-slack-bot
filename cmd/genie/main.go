package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"

	"github.com/pestophagous/brevity-genie-slack-bot/pkg/brevity"
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/trace"
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/util"
)

var approvedMissingEventCases []string = []string{"ack", "user_typing"}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Lshortfile for file.go:NN

	trace.EnableTraces([]string{"incoming-event"})
	trace.EnableTraces([]string{"user-metadata"})
	trace.EnableTraces([]string{"package-brevity"})
}

func main() {
	token := util.MustGetEnv("SLACKTOKEN")
	vtoken := util.MustGetEnv("SLACKVTOKEN_DEPRECATED")
	api := slack.New(token, slack.OptionDebug(true))

	brevityBot := brevity.NewBrevityBot(&slackAdapter{client: api})

	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: vtoken}))
		if e != nil {
			log.Print("ERROR: ", e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		trace.Tracef("incoming-event", "outer event type: %s", eventsAPIEvent.Type)

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				log.Print("ERROR: ", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
			return
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			trace.Tracef("incoming-event", "inner event type: %s", innerEvent.Type)

			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				trace.Tracef("incoming-event", "app mention: %v", ev)

			case *slackevents.MessageEvent:
				trace.Trace("incoming-event", ev.User)      // US8FL2R91
				trace.Trace("incoming-event", ev.Channel)   // CS1EW078Q
				trace.Trace("incoming-event", ev.Type)      // "message"
				trace.Trace("incoming-event", ev.TimeStamp) // 1578794317.001700 // t=1578794317001700 ; date -d @${t%??????}
				trace.Trace("incoming-event", ev.Text)
				trace.Trace("incoming-event", ev.BotID)
				trace.Trace("incoming-event", ev.Username)

				// only "react" to MessageEvent(s) that came from NON-bots
				if ev.BotID == "" {
					text := ev.Text

					text = strings.TrimSpace(text)
					text = strings.ToLower(text)

					matched, _ := regexp.MatchString("test poke the bot", text)

					if matched {
						api.PostMessage(ev.Channel, slack.MsgOptionText("sure thing", false))
					}

					excusals := brevityBot.Track(brevity.NewChatActivity(util.MessageTimestampToUTC(ev.TimeStamp), ev.User))

					for _, excuse := range excusals {
						api.PostMessage(ev.Channel, slack.MsgOptionText(excuse, false))
					}
				}
			}
		}
	})

	log.Print("About to ListenAndServe")
	// In order to bind to port 83, you need not use sudo if you instead do:
	//    sudo setcap CAP_NET_BIND_SERVICE=+eip genie # last word is the NAME OF YOUR SERVER BINARY.
	err := http.ListenAndServe(":83", nil) // not ALWAYS necessary to run as 'sudo' (see above)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
}
