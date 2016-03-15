package main

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
)

type Slack struct {
	Webhook string
	Proxy string
	Channel string
	Owners  string
}

func (s *Slack) Notify(check AppCheck) {
	attachment := slack.Attachment{
		Text:  &check.Message,
		Color: s.resultToColor(check),
	}
	attachment.
		AddField(slack.Field{Title: "App", Value: check.App, Short: true}).
		AddField(slack.Field{Title: "Check", Value: check.CheckName, Short: true}).
		AddField(slack.Field{Title: "Result", Value: s.resultToString(check), Short: true})

	destination := GetString(check.Labels, "alerts.slack.channel", s.Channel)

	appSpecificOwners := GetString(check.Labels, "alerts.slack.owners", s.Owners)
	var owners []string
	if appSpecificOwners != "" {
		owners = strings.Split(appSpecificOwners, ",")
	} else {
		owners = []string{"@here"}
	}

	alertSuffix := "Please check!"
	if check.Result == Pass {
		alertSuffix = "Check Resolved, thanks!"
	}
	mainText := fmt.Sprintf("%s, %s", s.parseOwners(owners), alertSuffix)

	payload := slack.Payload(mainText,
		"marathon-alerts",
		"",
		destination,
		[]slack.Attachment{attachment})

	webhooks := strings.Split(GetString(check.Labels, "alerts.slack.webhook", s.Webhook), ",")
	proxy := GetString(check.Labels, "alerts.slack.proxy", s.Proxy)

	for _, webhook := range webhooks {
		err := slack.Send(webhook, proxy, payload)
		if err != nil {
			fmt.Printf("Unexpected Error - %v", err)
		}
	}
}

func (s *Slack) resultToColor(check AppCheck) *string {
	color := "black"
	switch check.Result {
	case Pass:
		color = "good"
	case Warning:
		color = "warning"
	case Critical:
		color = "danger"
	}

	return &color
}

func (s *Slack) resultToString(check AppCheck) string {
	value := "Unknown"
	switch check.Result {
	case Pass:
		value = "Passed"
	case Warning:
		value = "Warning"
	case Critical:
		value = "Failed"
	}

	return value
}

func (s *Slack) parseOwners(owners []string) string {
	parsedOwners := []string{}
	for _, owner := range owners {
		if owner != "@here" {
			owner = fmt.Sprintf("@%s", owner)
		}
		parsedOwners = append(parsedOwners, owner)
	}

	return strings.Join(parsedOwners, ", ")
}
