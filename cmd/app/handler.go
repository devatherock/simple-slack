//go:build !plugin && !integration
// +build !plugin,!integration

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type NotificationRequest struct {
	Text    string `json:",omitempty"`
	Channel string `json:",omitempty"`
	Color   string `json:",omitempty"`
	Title   string `json:",omitempty"`
	Webhook string `json:",omitempty"`
	Token   string `json:",omitempty"`
	BuildId string `json:"build_id,omitempty"`
}

// Handles /api/notification endpoint. Waits for the supplied build
// to complete and then triggers a notification
func sendNotification(writer http.ResponseWriter, request *http.Request) {
	// Read request
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error("Error reading request: ", err)
		writer.WriteHeader(400)
		return
	}

	// Parse request
	notificationRequest := NotificationRequest{}
	err = json.Unmarshal(requestBody, &notificationRequest)
	if err != nil {
		log.Error("Error parsing request: ", err)
		writer.WriteHeader(400)
		return
	}

	// Use webhook from environment variable if available
	if notificationRequest.Webhook == "" {
		notificationRequest.Webhook = os.Getenv("SLACK_WEBHOOK")
	}

	successStatus, err := notify(notificationRequest)
	if err != nil {
		log.Error("Error sending notification: ", err)
		writer.WriteHeader(400)
	} else {
		writer.WriteHeader(successStatus)
	}
}

// Handles /api/health endpoint. Indicates the health of the application
func checkHealth(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("UP"))
}
