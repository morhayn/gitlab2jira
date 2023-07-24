package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/morhayn/gitlab2jira/internal/jira"
	"github.com/morhayn/gitlab2jira/internal/telegram"

	"github.com/xanzy/go-gitlab"
)

var Tocken string

// webhook is a HTTP Handler for Gitlab Webhook events.
var regex = `.*(SO\w+VO-\d*).*`

type webhook struct {
	Secret         string
	EventsToAccept []gitlab.EventType
}

// webhookExample shows how to create a Webhook server to parse Gitlab events.
func Webhook() {
	wh := webhook{
		Secret: Tocken,
		EventsToAccept: []gitlab.EventType{
			// gitlab.EventTypePush,
			// gitlab.EventTypeIssue,
			gitlab.EventTypeMergeRequest,
			// gitlab.EventTypePipeline,
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/webhook/", wh)
	if err := http.ListenAndServe("0.0.0.0:3000", mux); err != nil {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// ServeHTTP tries to parse Gitlab events sent and calls handle function
// with the successfully parsed events.
func (hook webhook) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("webhook OK")
	event, err := hook.parse(request)
	if err != nil {
		writer.WriteHeader(500)
		fmt.Println(writer, "could parse the webhook event: %v", err)
		return
	}
	if event == gitlab.EventTypeMergeRequest {
		if err := MergeWebhook(request); err != nil {
			writer.WriteHeader(500)
			fmt.Println("Error Parse Push Event, ", err)
			return
		}
	}
	// Write a response when were done.
	writer.WriteHeader(204)
}

// parse verifies and parses the events specified in the request and
// returns the parsed event or an error.
func (hook webhook) parse(r *http.Request) (gitlab.EventType, error) {
	var event gitlab.EventType
	// if r.Method != http.MethodPost {
	// fmt.Println(r.Method)
	// return event, errors.New("invalid HTTP Method")
	// }
	// If we have a secret set, we should check if the request matches it.
	if len(hook.Secret) > 0 {
		signature := r.Header.Get("X-Gitlab-Token")
		if signature != hook.Secret {
			return event, errors.New("token validation failed")
		}
	}
	e := r.Header.Get("X-Gitlab-Event")
	if strings.TrimSpace(e) == "" {
		return event, errors.New("missing X-Gitlab-Event Header")
	}
	eventType := gitlab.EventType(e)
	if !isEventSubscribed(eventType, hook.EventsToAccept) {
		return event, errors.New("event not defined to be parsed")
	}
	return eventType, nil
}
func MergeWebhook(r *http.Request) error {
	defer func() {
		if _, err := io.Copy(ioutil.Discard, r.Body); err != nil {
			log.Printf("could discard request body: %v", err)
		}
		if err := r.Body.Close(); err != nil {
			log.Printf("could not close request body: %v", err)
		}
	}()
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return errors.New("error reading request body")
	}
	push := gitlab.MergeEvent{}
	if err := json.Unmarshal(payload, &push); err != nil {
		return err
	}
	attr := push.ObjectAttributes
	// fmt.Println("STATE=", attr.State)
	fmt.Println(attr.State)
	if attr.State == "opened" {
		if telegram.Tocken != "" {
			c := telegram.New()
			chartId := -1001514590541
			// fmt.Println(int64(chartId))
			t1, err := time.Parse("2006-01-02 15:04:05 -0700", attr.CreatedAt)
			if err != nil {
				fmt.Println(err)
			}
			last := attr.LastCommit.Message
			messg := fmt.Sprintf("[Link to request](%s)\n*created* %s\n*creator* %s\n*last commit* ```%s```\n*assignees* %s\n*rewiewers* %s\n",
				attr.URL, t1.Format("2006-01-02 15:04:05"), push.User.Name, last, push.Assignees[0].Name,
				push.Reviewers[0].Name)
			// messg = strings.Replace(messg, "_", "\\_", -1)
			// messg = strings.Replace(messg, "*", "\\*", -1)
			// messg = strings.Replace(messg, "`", "\\`", -1)
			// messg = strings.Replace(messg, "[", "\\[", -1)
			// messg = strings.Replace(messg, "+", "\\+", -1)
			// messg = strings.Replace(messg, "-", "\\-", -1)
			// messg = strings.Replace(messg, "=", "\\=", -1)
			// messg = strings.Replace(messg, "|", "\\|", -1)
			// messg = strings.Replace(messg, ".", "\\.", -1)
			// messg = strings.Replace(messg, "!", "\\!", -1)
			// messg = strings.Replace(messg, "(", "\\(", -1)
			// messg = strings.Replace(messg, "~", "\\~", -1)
			// messg = strings.Replace(messg, "#", "\\#", -1)
			// messg = strings.Replace(messg, ">", "\\>", -1)
			fmt.Println(messg)
			err = c.SendMessage(messg, int64(chartId))
			if err != nil {
				fmt.Println(err)
			}
		}
		if jira.Tocken != "" {
			fmt.Printf("message: %s, url: %s \n", attr.Description, attr.URL)
			re := regexp.MustCompile(regex)
			matchDesc := re.FindStringSubmatch(attr.Description)
			matchTitle := re.FindStringSubmatch(attr.Title)

			if len(matchDesc) > 1 {
				ticket := matchDesc[1]
				jira.SendComment(ticket, attr.URL, attr.Description, attr.State, push.User.Name)
			} else if len(matchTitle) > 1 {
				ticket := matchTitle[1]
				jira.SendComment(ticket, attr.URL, attr.Title, attr.State, push.User.Name)
			}
		}
	}
	return nil
}

func isEventSubscribed(event gitlab.EventType, events []gitlab.EventType) bool {
	fmt.Println(event, events)
	for _, e := range events {
		if event == e {
			return true
		}
	}
	return false
}
