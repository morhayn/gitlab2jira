### gitlab2jira
### System env
- Gital tocken - GITLAB_TOCKEN
- jira tocken - JIRA_TOCKEN
- jira url - JIRA_URL

Build command
```
export GITLAB="WEBHOOK GITLAB" && \
export JIRA="TOCKEN JIRA" && \
export UrlJira="URL JIRA" && \
CGO_ENABLED=0 go build -o gitlab2jira -ldflags \
"-X github.com/morhayn/gitlab2jira/internal/webhook.Tocken=$GITLAB \
 -X github.com/morhayn/gitlab2jira/internal/jira.UrlJira=$UrlJira \
 -X github.com/morhayn/gitlab2jira/internal/jira.Tocken=$JIRA" main.go
 ```