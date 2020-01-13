package chat

import (
	"time"
)

// An interface providing 'getters' but no setters (no mutators)
type ReadonlyChat interface {
	StartTime() time.Time
	EndTime() time.Time
	Duration() time.Duration
	User() string
}

// This says: a pointer (nil) cast-to-Chat is assignable to interface type ReadonlyChat.
// In other words: compile-time enforcement/documentation that Chat satisfies ReadonlyChat.
var _ ReadonlyChat = (*Chat)(nil)

type Chat struct {
	start        time.Time
	lastActivity time.Time
	user         string
}

func NewChat(startTime time.Time, userID string) *Chat {
	return &Chat{
		start:        startTime,
		lastActivity: startTime,
		user:         userID,
	}
}

func (c *Chat) MutateLastTime(incomingLastTime time.Time) {
	if incomingLastTime.After(c.lastActivity) {
		c.lastActivity = incomingLastTime
	}
}

func (c *Chat) StartTime() time.Time {
	return c.start
}

func (c *Chat) EndTime() time.Time {
	return c.lastActivity
}

func (c *Chat) Duration() time.Duration {
	return c.lastActivity.Sub(c.start)
}

func (c *Chat) User() string {
	return c.user
}
