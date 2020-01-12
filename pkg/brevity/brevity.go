package brevity

import (
	"log"
)

type BrevityBot struct{}

func NewBrevityBot() *BrevityBot {
	return &BrevityBot{}
}

// Note: UserActivity is an interface defined in user_activity.go
func (bb *BrevityBot) Track(activity UserActivity) {
	log.Print(activity.Timestamp())
	log.Print(activity.UserID())
}
