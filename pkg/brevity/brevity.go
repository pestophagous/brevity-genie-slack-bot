package brevity

import (
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/trace"
)

type BrevityBot struct{}

func NewBrevityBot() *BrevityBot {
	return &BrevityBot{}
}

// Note: UserActivity is an interface defined in user_activity.go
func (bb *BrevityBot) Track(activity UserActivity) {
	trace.Trace("package-brevity", activity.Timestamp())
	trace.Trace("package-brevity", activity.UserID())
}
