//go:build test
// +build test

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devatherock/simple-slack/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestRunApp(test *testing.T) {
	helper.SetEnvironmentVariable(test, "TEXT", "Build failed!")

	// Test HTTP server
	var capturedRequest []byte
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		capturedRequest, _ = ioutil.ReadAll(request.Body)
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(writer, `{"success":true}`)
	}))
	defer testServer.Close()
	helper.SetEnvironmentVariable(test, "WEBHOOK", testServer.URL)

	runApp([]string{"-x", "dummy"})

	// Verify request
	jsonRequest := make(map[string]interface{})
	json.Unmarshal(capturedRequest, &jsonRequest)
	assert.Equal(test, 1, len(jsonRequest))

	var attachments []interface{}
	attachments = jsonRequest["attachments"].([]interface{})
	var attachment map[string]interface{}
	attachment = attachments[0].(map[string]interface{})

	assert.Equal(test, 1, len(attachments))
	assert.Equal(test, 2, len(attachment))
	assert.Equal(test, "Build failed!", attachment["text"])
	assert.Equal(test, "#cfd3d7", attachment["color"])
}

func TestRun(test *testing.T) {
	// Test HTTP server
	var capturedRequest []byte
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		capturedRequest, _ = ioutil.ReadAll(request.Body)
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(writer, `{"success":true}`)
	}))
	defer testServer.Close()

	set := flag.NewFlagSet("test", 0)
	set.String("text", "Build failed!", "")
	set.String("color", "red", "")
	set.String("title", "Build notification", "")
	set.String("channel", "general", "")
	set.String("webhook", testServer.URL, "")

	context := cli.NewContext(nil, set, nil)
	actual := run(context)

	// Verify no error
	assert.Nil(test, actual)

	// Verify request
	jsonRequest := make(map[string]interface{})
	json.Unmarshal(capturedRequest, &jsonRequest)
	assert.Equal(test, 2, len(jsonRequest))

	var attachments []interface{}
	attachments = jsonRequest["attachments"].([]interface{})
	var attachment map[string]interface{}
	attachment = attachments[0].(map[string]interface{})

	assert.Equal(test, 1, len(attachments))
	assert.Equal(test, 3, len(attachment))
	assert.Equal(test, "Build failed!", attachment["text"])
	assert.Equal(test, "red", attachment["color"])
	assert.Equal(test, "Build notification", attachment["title"])
	assert.Equal(test, "general", jsonRequest["channel"])
}

func TestRunSlackCallFailed(test *testing.T) {
	cases := []int{
		400,
		500,
	}

	for _, statusCode := range cases {
		// Test HTTP server
		var capturedRequest []byte
		testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			capturedRequest, _ = ioutil.ReadAll(request.Body)
			writer.WriteHeader(statusCode)
		}))
		defer testServer.Close()

		set := flag.NewFlagSet("test", 0)
		set.String("text", "Build failed!", "")
		set.String("color", "red", "")
		set.String("title", "Build notification", "")
		set.String("channel", "general", "")
		set.String("webhook", testServer.URL, "")

		context := cli.NewContext(nil, set, nil)
		actual := run(context)

		// Verify error
		assert.NotNil(test, actual)
		assert.Equal(test, "HTTP request to Slack failed", actual.Error())

		// Verify request
		jsonRequest := make(map[string]interface{})
		json.Unmarshal(capturedRequest, &jsonRequest)
		assert.Equal(test, 2, len(jsonRequest))

		var attachments []interface{}
		attachments = jsonRequest["attachments"].([]interface{})
		var attachment map[string]interface{}
		attachment = attachments[0].(map[string]interface{})

		assert.Equal(test, 1, len(attachments))
		assert.Equal(test, 3, len(attachment))
		assert.Equal(test, "Build failed!", attachment["text"])
		assert.Equal(test, "red", attachment["color"])
		assert.Equal(test, "Build notification", attachment["title"])
		assert.Equal(test, "general", jsonRequest["channel"])
	}
}
