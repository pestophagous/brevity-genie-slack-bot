package brevity

import (
	"time"
)

type UserActivity interface {
	Timestamp() time.Time
	UserID() string
}

// ChatActivity is a concrete subtype of the UserActivity interface
type ChatActivity struct {
	timestamp time.Time
	userID    string
}

func NewChatActivity(ts time.Time, uID string) *ChatActivity {
	return &ChatActivity{
		timestamp: ts,
		userID:    uID,
	}
}

func (ca *ChatActivity) Timestamp() time.Time {
	return ca.timestamp
}

func (ca *ChatActivity) UserID() string {
	return ca.userID
}
