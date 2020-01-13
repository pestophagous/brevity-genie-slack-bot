package user

import (
	"fmt"
	"time"

	"github.com/pestophagous/brevity-genie-slack-bot/pkg/chat"
)

var arbitraryOldDate time.Time = time.Date(2000, time.June, 1, 23, 0, 0, 0, time.UTC)

type User struct {
	id            string
	ongoingChat   *chat.Chat
	brevityLimits UserSettings
	history       *chat.ChatHistory
}

func NewUser(userId string) *User {
	return &User{
		id:            userId,
		ongoingChat:   nil,
		brevityLimits: SettingsForUser(userId),
		history:       chat.NewChatHistory(userId, time.Duration(24)*time.Hour),
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) hasInProgressActivity() bool {
	return u.ongoingChat != nil
}

func (u *User) fulfillsSilenceBarrierAsOf(when time.Time) bool {
	if !u.hasInProgressActivity() {
		return true
	}

	silence := when.Sub(u.ongoingChat.EndTime())
	return silence.Minutes() > float64(u.brevityLimits.SilenceBarrierMinutes())
}

func (u *User) lastSeen() time.Time {
	if u.ongoingChat != nil {
		return u.ongoingChat.EndTime()
	} else if !u.history.Empty() {
		return u.history.LastActivityTime()
	} else {
		return arbitraryOldDate
	}
}

func (u *User) DurationSinceLastSeen() time.Duration {
	return time.Since(u.lastSeen())
}

func (u *User) closeExistingIfWarranted(latestTimestamp time.Time) {
	if u.fulfillsSilenceBarrierAsOf(latestTimestamp) {
		// close any existing.
		if u.hasInProgressActivity() {
			u.history.AddChat(u.ongoingChat)
			u.ongoingChat = nil
		}
	}
}

func (u *User) Track(latestTimestamp time.Time) {
	u.closeExistingIfWarranted(latestTimestamp)

	if !u.hasInProgressActivity() {
		// start new
		u.ongoingChat = chat.NewChat(latestTimestamp, u.id)
	}

	u.ongoingChat.MutateLastTime(latestTimestamp)
}

func (u *User) GetExcusalAsOf(latestTimestamp time.Time) string {
	u.closeExistingIfWarranted(latestTimestamp)

	if !u.hasInProgressActivity() {
		return ""
	}

	if u.ongoingChat.Duration().Minutes() > float64(u.brevityLimits.ChatCapMinutes()) {
		return fmt.Sprintf("%s must be getting back to work now :)", u.id)
	}

	return ""
}
