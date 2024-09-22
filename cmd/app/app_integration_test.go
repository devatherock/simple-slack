//go:build api
// +build api

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/devatherock/simple-slack/test/helper"
	"github.com/stretchr/testify/assert"
)

const baseUrl string = "http://localhost:8082"

var client = &http.Client{}
var inProgressIndexes = []int{0, 3, 6}

func TestSendNotificationWithBuildId(test *testing.T) {
	var requestIndex = 0
	var capturedRequest []byte
	var requestUrl = ""
	var headers map[string][]string

	// Test HTTP server that listens on a specific port
	listener, _ := net.Listen("tcp", "127.0.0.1:8085")
	testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		capturedRequest, _ = ioutil.ReadAll(request.Body)
		writer.Header().Set("Content-Type", "application/json")

		if slices.Contains(inProgressIndexes, requestIndex) {
			requestUrl = request.Host + request.URL.Path
			headers = request.Header

			fmt.Fprintln(writer, `{"project_slug":"gh/devatherock/email-sender","pipeline_number":292,"status":"running"}`)
		} else if requestIndex == 1 {
			fmt.Fprintln(writer, `{"project_slug":"gh/devatherock/email-sender","pipeline_number":292,"status":"success"}`)
		} else if requestIndex == 4 {
			fmt.Fprintln(writer, `{"project_slug":"gh/devatherock/email-sender","pipeline_number":292,"status":"failed"}`)
		} else if requestIndex == 7 {
			fmt.Fprintln(writer, `{"project_slug":"gh/devatherock/email-sender","pipeline_number":292,"status":"failing"}`)
		} else { // 2, 5, 8 For the slack send call
			fmt.Fprintln(writer, `{"success":true}`)
		}

		requestIndex++
	}))

	// NewUnstartedServer creates a listener. Close that listener and replace with the one we created.
	testServer.Listener.Close()
	testServer.Listener = listener

	// Start the server
	testServer.Start()

	// Stop the server on return from the function
	defer testServer.Close()

	cases := []struct {
		expectedStatus string
		expectedColor  string
	}{
		{
			"Success",
			"#33ad7f",
		},
		{
			"Failed",
			"#a1040c",
		},
		{
			"Failing",
			"#a1040c",
		},
	}

	progressIndex := 0
	for _, data := range cases {
		notificationRequest := map[string]interface{}{
			"text":     "",
			"channel":  "general",
			"color":    "blue",
			"title":    "Build completed",
			"webhook":  "http://localhost:8085",
			"build_id": "1234",
		}

		jsonStr, _ := json.Marshal(&notificationRequest)
		request, _ := http.NewRequest("POST", baseUrl+"/api/notification", bytes.NewBuffer(jsonStr))

		response, err := client.Do(request)
		io.Copy(ioutil.Discard, response.Body)
		defer response.Body.Close()

		assert.Nil(test, err)
		assert.Equal(test, 204, response.StatusCode)

		// Wait for request to be processed completely
		for requestIndex-progressIndex < 2 {
			time.Sleep(1 * time.Second)
		}
		progressIndex = requestIndex

		// Verify first CircleCI request
		assert.Equal(test, "localhost:8085/api/v2/workflow/1234", requestUrl)
		assert.Equal(test, "dummy", headers["Circle-Token"][0])

		// Verify slack request
		helper.VerifySlackRequest(test, capturedRequest, map[string]string{
			"text":  fmt.Sprintf("%s: <https://app.circleci.com/pipelines/gh/devatherock/email-sender/292|gh/devatherock/email-sender-292>", data.expectedStatus),
			"color": data.expectedColor,
			"title": "Build completed",
		})
	}
}

func TestSendNotificationNoBuildId(test *testing.T) {
	var capturedRequest []byte

	// Test HTTP server
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		capturedRequest, _ = ioutil.ReadAll(request.Body)
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(writer, `{"success":true}`)
	}))

	// Stop the server on return from the function
	defer testServer.Close()

	cases := []struct {
		text         string
		color        string
		expectedText string
	}{
		{
			"Failed",
			"red",
			"Failed",
		},
		{
			"",
			"green",
			"Build completed",
		},
	}

	for _, data := range cases {
		notificationRequest := map[string]interface{}{
			"text":    data.text,
			"channel": "general",
			"color":   data.color,
			"title":   "some title",
			"webhook": testServer.URL,
		}

		jsonStr, _ := json.Marshal(&notificationRequest)
		request, _ := http.NewRequest("POST", baseUrl+"/api/notification", bytes.NewBuffer(jsonStr))

		response, err := client.Do(request)
		io.Copy(ioutil.Discard, response.Body)
		defer response.Body.Close()

		assert.Nil(test, err)
		assert.Equal(test, 200, response.StatusCode)

		// Verify slack request
		helper.VerifySlackRequest(test, capturedRequest, map[string]string{
			"text":  data.expectedText,
			"color": data.color,
			"title": "some title",
		})
	}
}

func TestSendNotificationNoWebhook(test *testing.T) {
	notificationRequest := map[string]interface{}{
		"text":    "",
		"channel": "general",
		"color":   "red",
		"title":   "some title",
	}

	jsonStr, _ := json.Marshal(&notificationRequest)
	request, _ := http.NewRequest("POST", baseUrl+"/api/notification", bytes.NewBuffer(jsonStr))

	response, err := client.Do(request)
	io.Copy(ioutil.Discard, response.Body)
	defer response.Body.Close()

	assert.Nil(test, err)
	assert.Equal(test, 400, response.StatusCode)
}

func TestSendNotificationErrorFromSlack(test *testing.T) {
	notificationRequest := map[string]interface{}{
		"text":    "",
		"channel": "general",
		"color":   "red",
		"title":   "some title",
		"webhook": "http://localhost:8085",
	}

	jsonStr, _ := json.Marshal(&notificationRequest)
	request, _ := http.NewRequest("POST", baseUrl+"/api/notification", bytes.NewBuffer(jsonStr))

	response, err := client.Do(request)
	io.Copy(ioutil.Discard, response.Body)
	defer response.Body.Close()

	assert.Nil(test, err)
	assert.Equal(test, 400, response.StatusCode)
}

func TestSendNotificationInvalidJson(test *testing.T) {
	request, _ := http.NewRequest("POST", baseUrl+"/api/notification", bytes.NewBuffer([]byte("some text")))

	response, err := client.Do(request)
	io.Copy(ioutil.Discard, response.Body)
	defer response.Body.Close()

	assert.Nil(test, err)
	assert.Equal(test, 400, response.StatusCode)
}

func TestCheckHealth(test *testing.T) {
	request, _ := http.NewRequest("GET", baseUrl+"/api/health", nil)

	response, err := client.Do(request)
	defer response.Body.Close()

	assert.Nil(test, err)
	assert.Equal(test, 200, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(test, "UP", string(responseBody))
}
