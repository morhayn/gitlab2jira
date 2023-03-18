package main

import (
	"github.com/morhayn/gitlab2jira/internal/webhook"
)

func main() {
	// fmt.Println(webhook.Tocken)
	// fmt.Println(jira.Tocken)
	webhook.Webhook()
}
