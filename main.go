package main

import (
	"os"

	"github.com/morhayn/gitlab2jira/internal/jira"
	"github.com/morhayn/gitlab2jira/internal/telegram"
	"github.com/morhayn/gitlab2jira/internal/webhook"
)

func main() {
	// fmt.Println(webhook.Tocken)
	// fmt.Println(jira.Tocken)
	if webhook.Tocken == "" {
		if tocken := os.Getenv("GITLAB"); tocken != "" {
			webhook.Tocken = tocken
		}
	}
	if jira.Tocken == "" {
		if tocken := os.Getenv("JIRA"); tocken != "" {
			jira.Tocken = tocken
		}
	}
	if telegram.Tocken == "" {
		if tocken := os.Getenv("TELE"); tocken != "" {
			telegram.Tocken = tocken
		}
	}
	if jira.UrlJira == "" {
		if jiraUrl := os.Getenv("JIRA_URL"); jiraUrl != "" {
			jira.UrlJira = jiraUrl
		}
	}
	webhook.Webhook()
}
