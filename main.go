package main

import (
	"os"

	"github.com/morhayn/gitlab2jira/internal/jira"
	"github.com/morhayn/gitlab2jira/internal/webhook"
)

func main() {
	// fmt.Println(webhook.Tocken)
	// fmt.Println(jira.Tocken)
	if webhook.Tocken == "" {
		if tocken := os.Getenv("GITLAB_TOCKEN"); tocken != "" {
			webhook.Tocken = tocken
		}
	}
	if jira.Tocken == "" {
		if tocken := os.Getenv("JIRA_TOCKEN"); tocken != "" {
			jira.Tocken = tocken
		}
	}
	if jira.UrlJira == "" {
		if jiraUrl := os.Getenv("JIRA_URL"); jiraUrl != "" {
			jira.UrlJira = jiraUrl
		}
	}
	webhook.Webhook()
}
