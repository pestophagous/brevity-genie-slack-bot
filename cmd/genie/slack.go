package main

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"

	"github.com/pestophagous/brevity-genie-slack-bot/pkg/backend_contract"
	"github.com/pestophagous/brevity-genie-slack-bot/pkg/trace"
)

type slackAdapter struct {
	client *slack.Client
}

type slackUserAdapter struct {
	slackUser slack.User
}

func (s *slackAdapter) GetUsers() ([]backend_contract.UserData, error) {
	users, err := s.client.GetUsers()
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return []backend_contract.UserData{}, err
	}

	results := []backend_contract.UserData{}

	for _, u := range users {
		trace.Trace("user-metadata", "----- (querying users)")
		trace.Trace("user-metadata", "ID:", u.ID)
		trace.Trace("user-metadata", "Name:", u.Name)
		trace.Trace("user-metadata", "RealName:", u.RealName)
		trace.Trace("user-metadata", "IsBot:", u.IsBot)
		trace.Trace("user-metadata", "IsAdmin:", u.IsAdmin)
		trace.Trace("user-metadata", "IsOwner:", u.IsOwner)
		trace.Trace("user-metadata", "IsPrimaryOwner:", u.IsPrimaryOwner)

		results = append(results, &slackUserAdapter{slackUser: u})
	}

	return results, nil
}

func (u *slackUserAdapter) ID() string {
	return u.slackUser.ID
}

func (u *slackUserAdapter) HumanReadableName() string {
	return u.slackUser.RealName
}

func (u *slackUserAdapter) AtHandle() string {
	return fmt.Sprintf("<@%s>", u.slackUser.ID)
}
