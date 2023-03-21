package jira

import (
	"context"
	"fmt"

	j "github.com/andygrunwald/go-jira/v2/onpremise"
)

var Tocken string
var UrlJira string

func SendComment(ticket, gitlink, message, state, username string) {
	jiraURL := UrlJira

	// See "Using Personal Access Tokens"
	// https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html
	tp := j.BearerAuthTransport{
		Token: Tocken,
	}
	client, err := j.NewClient(jiraURL, tp.Client())
	if err != nil {
		panic(err)
	}

	// Running an empty JQL query to get all tickets
	jql := fmt.Sprintf(`project = "SODIEVO" AND key = "%s"`, ticket)
	_, _, err = client.Issue.Search(context.Background(), jql, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	comment := j.Comment{
		Body: fmt.Sprintf("Owner: %s \n State: %s \n Gitlab: %s \n Description: %s", username, state, gitlink, message),
	}
	_, _, err = client.Issue.AddComment(context.Background(), ticket, &comment)
	if err != nil {
		fmt.Println("ERROR -- ", err)
	}
}
