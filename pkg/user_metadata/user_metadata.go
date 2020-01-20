package user_metadata

import (
	"time"

	"github.com/pestophagous/brevity-genie-slack-bot/pkg/backend_contract"
)

type UserMetadata interface {
	Id() string
	HandleForAtMessages() string
	HumanReadableName() string

	MetadataAsOfDate() time.Time
}

type userMetadata struct {
	id        string
	atHandle  string
	humanName string

	lastFetched time.Time
}

func MetadataForUser(userId string, backend backend_contract.Api) UserMetadata {
	return dataForUser(userId, backend)
}

func (u *userMetadata) Id() string {
	// TODO: staleness check and async re-query
	return u.id
}

func (u *userMetadata) HandleForAtMessages() string {
	// TODO: staleness check and async re-query
	return u.atHandle
}

func (u *userMetadata) HumanReadableName() string {
	// TODO: staleness check and async re-query
	return u.humanName
}

func (u *userMetadata) MetadataAsOfDate() time.Time {
	// TODO: staleness check and async re-query
	return u.lastFetched
}
