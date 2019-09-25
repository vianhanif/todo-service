package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const apiURL = "https://slack.com/api"
const defaultHTTPTimeout = 80 * time.Second

// SlackAlertConfig represent the config needed when creating a new slack notifier
type SlackAlertConfig struct {
	Token      string
	Channel    string
	HTTPClient *http.Client
}

// SlackAlert represents the notifier that will notify to slack channel
type SlackAlert struct {
	Token      string
	Channel    string
	HTTPClient *http.Client
}

// Alert alerts message to a slack channel
func (sn *SlackAlert) Alert(message Message) {
	/*
		Examples of calling the slack API:

			curl -X POST -H 'Authorization: Bearer xoxb-1234-56789abcdefghijklmnop' \
			-H 'Content-type: application/json' \
			--data '{
				"channel":"C061EG9SL",
				"text":"I hope the tour went well, Mr. Wonka.",
				"attachments": [{
					"text":"Who wins the lifetime supply of chocolate?",
					"fallback":"You could be telling the computer exactly what it can do with a lifetime supply of chocolate.",
					"color":"#3AA3E3",
					"attachment_type":"default",
					"callback_id":"select_simple_1234",
					"actions":[{
						"name":"winners_list",
						"text":"Who should win?",
						"type":"select",
						"data_source":"users"
					}]
				}]
			}' \
			https://slack.com/api/chat.postMessage
	*/

	payload := map[string]interface{}{
		"channel": sn.Channel,
	}
	if len(message.Trace) > 0 && message.Error != nil {
		var errMessage string
		var traceMessage string
		errMessage = message.Error.Error()
		traceMessage = string(message.Trace)
		payload["attachments"] = []interface{}{
			map[string]interface{}{
				"text": strings.Join([]string{errMessage, traceMessage}, "\n"),
			},
		}
	} else {
		payload["text"] = message.Text
		if len(message.Trace) > 0 {
			payload["attachments"] = []interface{}{
				map[string]interface{}{
					"text": string(message.Trace),
				},
			}
		}
	}
	bs, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error: %v was occured while trying send message to slack.\nMessage was: %v", err, message)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat.postMessage", apiURL), bytes.NewBuffer(bs))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sn.Token))
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		log.Printf("error: %v was occured while trying send message to slack.\nMessage was: %v", err, message)
	}

	res, err := sn.HTTPClient.Do(req)
	if err != nil {
		log.Printf("error: %v was occured while trying send message to slack.\nMessage was: %v", err, message)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("error: %v was occured while trying send message to slack.\nMessage was: %v", err, message)
	}
}

// NewSlackAlert creates a new slack notifier
func NewSlackAlert(config SlackAlertConfig) *SlackAlert {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	return &SlackAlert{
		Token:      config.Token,
		Channel:    config.Channel,
		HTTPClient: config.HTTPClient,
	}
}
