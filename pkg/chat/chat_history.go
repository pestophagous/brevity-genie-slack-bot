package chat

import (
	"log"
	"sort"
	"time"
)

type ChatHistory struct {
	chats           []ReadonlyChat
	user            string
	historyDuration time.Duration
}

func NewChatHistory(userID string, lookback time.Duration) *ChatHistory {
	return &ChatHistory{
		user:            userID,
		historyDuration: lookback,
	}
}

func (ch *ChatHistory) cull() {
	sort.Sort(chatSortByLastActivity(ch.chats))

	var earliestTimeToKeep time.Time = time.Now().Add(-ch.historyDuration)
	var culledChats []ReadonlyChat

	for _, chat := range ch.chats {
		// As of Jan 2020, I can't decide whether StartTime or EndTime is 'best' here. Seems arbitrary.
		if chat.EndTime().After(earliestTimeToKeep) {
			culledChats = append(culledChats, chat)
		}
	}

	ch.chats = culledChats
}

func (ch *ChatHistory) AddChat(chat ReadonlyChat) {
	if chat.User() != ch.user {
		log.Printf("Invalid request to add chat by %s to history for %s.", chat.User(), ch.user)
		return
	}

	ch.chats = append(ch.chats, chat)

	sort.Sort(chatSortByLastActivity(ch.chats))
	ch.cull()
}

func (ch *ChatHistory) GetTotalHistoryDuration() time.Duration {
	ch.cull()
	var totalMinutes float64

	for _, chat := range ch.chats {
		totalMinutes += chat.Duration().Minutes()
	}

	return time.Duration(totalMinutes) * time.Minute
}

func (ch *ChatHistory) Empty() bool {
	ch.cull()
	return len(ch.chats) == 0
}

func (ch *ChatHistory) LastActivityTime() time.Time {
	ch.cull()

	if ch.Empty() {
		log.Print("Request for LastActivityTime when history is empty.")
		return time.Time{}
	} else {
		sort.Sort(chatSortByLastActivity(ch.chats))
		return ch.chats[len(ch.chats)-1].EndTime()
	}
}
