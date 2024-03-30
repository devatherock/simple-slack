//go:build !plugin && !integration
// +build !plugin,!integration

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/devatherock/simple-slack/pkg/slack"
	log "github.com/sirupsen/logrus"
)

var httpClient = &http.Client{}

type CircleCiWorkFlow struct {
	Project        string `json:"project_slug,omitempty"`
	Status         string `json:",omitempty"`
	PipelineNumber int    `json:"pipeline_number,omitempty"`
}

func notify(notificationRequest NotificationRequest) (statusCode int, err error) {
	statusCode = 200

	slackRequest := slack.SlackRequest{}
	slackRequest.Text = notificationRequest.Text
	slackRequest.Color = notificationRequest.Color
	slackRequest.Title = notificationRequest.Title
	slackRequest.Channel = notificationRequest.Channel
	slackRequest.Webhook = notificationRequest.Webhook

	if slackRequest.Webhook == "" {
		statusCode = 400
		err = errors.New("webhook not specified")
		return
	}

	if notificationRequest.BuildId == "" {
		defaultTextIfMissing(&slackRequest)
		err = slack.Notify(slackRequest)
	} else {
		statusCode, err = notifyOnBuildCompletion(notificationRequest.BuildId, notificationRequest.Token, slackRequest)
	}

	return
}

func notifyOnBuildCompletion(buildId string, token string, slackRequest slack.SlackRequest) (int, error) {
	if token == "" {
		token = os.Getenv("CIRCLECI_TOKEN")
	}

	if token == "" {
		log.Warn("No token found, but build id specified. Build id: ", buildId)
		defaultTextIfMissing(&slackRequest)
		err := slack.Notify(slackRequest)

		if err != nil {
			return 400, err
		} else {
			return 200, nil
		}
	} else {
		go monitor(buildId, token, slackRequest)
	}

	return 204, nil
}

func monitor(buildId string, token string, slackRequest slack.SlackRequest) {
	log.Info("Monitoring build ", buildId)
	buildStatus := "running"

	for buildStatus == "running" {
		circleCiRequest, _ := http.NewRequest("GET", getCircleCiUrl()+"/api/v2/workflow/"+buildId, nil)
		circleCiRequest.Header.Add("Circle-Token", token)

		// Make request
		circleCiResponse, err := httpClient.Do(circleCiRequest)
		if err != nil {
			log.Error("error: ", err)
			return
		}

		// Read response
		circleCiResponseBody, _ := io.ReadAll(circleCiResponse.Body)
		circleCiWorkFlow := CircleCiWorkFlow{}
		json.Unmarshal(circleCiResponseBody, &circleCiWorkFlow)
		circleCiResponse.Body.Close()

		buildStatus = circleCiWorkFlow.Status
		if buildStatus == "success" || buildStatus == "failed" {
			if slackRequest.Text == "" {
				slackRequest.Text = fmt.Sprintf(
					"%s: <https://app.circleci.com/pipelines/%s/%d|%s-%d>",
					strings.ToUpper(buildStatus[:1])+buildStatus[1:],
					circleCiWorkFlow.Project,
					circleCiWorkFlow.PipelineNumber,
					circleCiWorkFlow.Project,
					circleCiWorkFlow.PipelineNumber,
				)
			}

			if buildStatus == "success" {
				slackRequest.Color = "#33ad7f"
			} else {
				slackRequest.Color = "#a1040c"
			}

			slack.Notify(slackRequest)
			break
		} else {
			// Wait if the build hasn't completed yet
			log.Debug("Waiting for build ", buildId)
			time.Sleep(5 * time.Second)
		}
	}

	log.Info("Sent notification for build ", buildId)
}

func getCircleCiUrl() (circleCiUrl string) {
	circleCiUrl = os.Getenv("CIRCLECI_API_HOST")

	if circleCiUrl == "" {
		circleCiUrl = "https://circleci.com"
	}

	return
}

func defaultTextIfMissing(slackRequest *slack.SlackRequest) {
	if slackRequest.Text == "" {
		slackRequest.Text = "Build completed"
	}
}
