package user_metadata

import (
	"fmt"
	"log"
	"time"

	"github.com/pestophagous/brevity-genie-slack-bot/pkg/backend_contract"
)

var store *backingStore

type backingStore struct {
	backend       backend_contract.Api
	lookup        map[string]*userMetadata
	lastRefreshed time.Time
}

func dataForUser(userId string, backend backend_contract.Api) UserMetadata {
	if store == nil {
		store = &backingStore{
			backend: backend,
			lookup:  make(map[string]*userMetadata),
		}
		refreshMetadataFromBackend()
	}

	if store.backend != backend {
		panic("This code is not designed to switch backends mid-run")
	}

	if data, ok := store.lookup[userId]; ok {
		return data
	}

	// refresh from backend. try one more time
	refreshMetadataFromBackend()

	if data, ok := store.lookup[userId]; ok {
		return data
	}

	return &userMetadata{
		id:          userId,
		atHandle:    fmt.Sprintf("<@%s>", userId),
		humanName:   fmt.Sprintf("NO-NAME,id=%s", userId),
		lastFetched: time.Now(),
	}
}

func refreshMetadataFromBackend() {
	if store == nil {
		panic("backing store MUST be initialized before you ever call into here")
	}

	users, err := store.backend.GetUsers()
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return
	}

	timeOfThisRefresh := time.Now()

	for _, user := range users {
		var metadata *userMetadata
		if _, ok := store.lookup[user.ID()]; !ok {
			// if we didn't find a lookup entry for this user, then make one:
			store.lookup[user.ID()] = &userMetadata{}
		}
		// now use the lookup entry for this user:
		metadata = store.lookup[user.ID()]

		metadata.id = user.ID()
		metadata.atHandle = user.AtHandle()
		metadata.humanName = user.HumanReadableName()
		metadata.lastFetched = timeOfThisRefresh
	}

	store.lastRefreshed = timeOfThisRefresh
}
