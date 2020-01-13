package brevity

import (
	"time"

	"github.com/pestophagous/brevity-genie-slack-bot/pkg/trace"
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/user"
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/util"
)

type BrevityBot struct {
	recentUsers map[string]*user.User
}

func NewBrevityBot() *BrevityBot {
	return &BrevityBot{
		recentUsers: make(map[string]*user.User),
	}
}

// Note: UserActivity is an interface defined in user_activity.go
func (bb *BrevityBot) Track(activity UserActivity) []string {
	trace.Trace("package-brevity", activity.Timestamp())
	trace.Trace("package-brevity", activity.UserID())

	user := bb.findUser(activity.UserID())

	user.Track(activity.Timestamp())

	// cull

	return bb.getExcusals([]string{activity.UserID()}, activity.Timestamp())
}

func (bb *BrevityBot) getExcusals(userIds []string, timestamp time.Time) []string {
	var excusals []string
	for _, user := range bb.recentUsers {
		ex := user.GetExcusalAsOf(timestamp)
		if ex != "" && util.MatchStringInSlice(user.Id(), userIds) {
			excusals = append(excusals, ex)
		}
	}
	return excusals
}

func (bb *BrevityBot) findUser(userId string) *user.User {
	if user, ok := bb.recentUsers[userId]; ok {
		return user
	}

	user := user.NewUser(userId)
	bb.recentUsers[userId] = user
	return user
}
