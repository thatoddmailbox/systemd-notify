package main

import (
	"fmt"
	"log"

	"github.com/NoteToScreen/teams-go/teams"
	"github.com/PromptBit/slack-go/slack"
)

func notify(unit string, state string, substate string) {
	summary := fmt.Sprintf("Unit %s is now %s (%s)", unit, state, substate)
	log.Println(summary)

	if currentConfig.Notify.Slack.Enabled {
		message := slack.Message{
			Text: fmt.Sprintf("Unit *%s* is now *%s* _(%s)_", unit, state, substate),
		}

		err := slack.PostToWebhook(currentConfig.Notify.Slack.WebhookURL, message)
		if err != nil {
			log.Println("Error while posting to Slack:", err)
		}
	}

	if currentConfig.Notify.Teams.Enabled {
		card := teams.Card{
			Summary:    summary,
			ThemeColor: "0000FF",
			Title:      "Unit status change",
			Sections: []teams.Section{
				teams.Section{
					ActivityTitle: summary,
					Text:          "",
					Facts:         []teams.Fact{},
				},
			},
		}

		err := teams.PostToWebhook(currentConfig.Notify.Teams.WebhookURL, card)
		if err != nil {
			log.Println("Error while posting to Teams:", err)
		}
	}
}
