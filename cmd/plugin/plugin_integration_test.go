//go:build integration
// +build integration

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	localhelper "github.com/devatherock/simple-slack/test/helper"
	"github.com/devatherock/vela-template-tester/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestMinimalParameters(test *testing.T) {
	// Test HTTP server
	var capturedRequest []byte
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		capturedRequest, _ = ioutil.ReadAll(request.Body)
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(writer, `{"success":true}`)
	}))
	defer testServer.Close()

	cases := []struct {
		textVariable    string
		webhookVariable string
		message         string
	}{
		{
			"TEXT",
			"WEBHOOK",
			"Build failed! Slack URL: {{.Webhook}}",
		},
		{
			"PLUGIN_TEXT",
			"PLUGIN_WEBHOOK",
			"Build failed! Slack URL: {{.PluginWebhook}}",
		},
		{
			"PARAMETER_TEXT",
			"SLACK_WEBHOOK",
			"Build failed! Slack URL: {{.SlackWebhook}}",
		},
	}

	for _, data := range cases {
		parameters := map[string]string{
			data.textVariable:    data.message,
			data.webhookVariable: testServer.URL,
		}
		arguments := []string{"run", "--rm", "--network", "host"}

		for key, value := range parameters {
			arguments = append(arguments, "-e")
			arguments = append(arguments, fmt.Sprintf("%s=%s", key, value))
		}
		arguments = append(arguments, "devatherock/simple-slack:latest")

		exitCode, output := helper.ExecuteCommand(exec.Command("docker", arguments...))
		fmt.Print(output)
		assert.Equal(test, 0, exitCode)

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
		assert.Equal(test, "Build failed! Slack URL: <no value>", attachment["text"])
		assert.Equal(test, "#cfd3d7", attachment["color"])
	}
}

func TestRequiredParametersMissing(test *testing.T) {
	arguments := []string{"run", "--rm", "devatherock/simple-slack:latest"}
	exitCode, output := helper.ExecuteCommand(exec.Command("docker", arguments...))
	fmt.Print(output)

	assert.Equal(test, 1, exitCode)
}

func TestAllParameters(test *testing.T) {
	// Test HTTP server
	var capturedRequest []byte
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		capturedRequest, _ = ioutil.ReadAll(request.Body)
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(writer, `{"success":true}`)
	}))
	defer testServer.Close()

	cases := []struct {
		textVariable    string
		webhookVariable string
		colorVariable   string
		titleVariable   string
		channelVariable string
		message         string
	}{
		{
			"TEXT",
			"WEBHOOK",
			"COLOR",
			"TITLE",
			"CHANNEL",
			`Failure: {{list .CircleBuildUrl ", Slack URL: " .Webhook | join "" }}`,
		},
		{
			"PLUGIN_TEXT",
			"PLUGIN_WEBHOOK",
			"PLUGIN_COLOR",
			"PLUGIN_TITLE",
			"PLUGIN_CHANNEL",
			`Failure: {{list .CircleBuildUrl ", Slack URL: " .PluginWebhook | join "" }}`,
		},
		{
			"PARAMETER_TEXT",
			"SLACK_WEBHOOK",
			"PARAMETER_COLOR",
			"PARAMETER_TITLE",
			"PARAMETER_CHANNEL",
			`Failure: {{list .CircleBuildUrl ", Slack URL: " .SlackWebhook | join "" }}`,
		},
	}

	for _, data := range cases {
		parameters := map[string]string{
			data.textVariable:    data.message,
			data.webhookVariable: testServer.URL,
			data.colorVariable:   "red",
			data.titleVariable:   "Build notification",
			data.channelVariable: "general",
		}
		parameters["CIRCLE_BUILD_URL"] = "https://someurl"

		arguments := []string{"run", "--rm", "--network", "host"}
		for key, value := range parameters {
			arguments = append(arguments, "-e")
			arguments = append(arguments, fmt.Sprintf("%s=%s", key, value))
		}
		arguments = append(arguments, "devatherock/simple-slack:latest")

		exitCode, output := helper.ExecuteCommand(exec.Command("docker", arguments...))
		fmt.Print(output)
		assert.Equal(test, 0, exitCode)

		// Verify request
		localhelper.VerifySlackRequest(test, capturedRequest, map[string]string{
			"text":  "Failure: https://someurl, Slack URL: ",
			"color": "red",
			"title": "Build notification",
		})
	}
}

func TestSlackError(test *testing.T) {
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

		parameters := map[string]string{
			"TEXT":             "Failure: {{.CircleBuildUrl}}",
			"WEBHOOK":          testServer.URL,
			"CIRCLE_BUILD_URL": "https://someurl",
		}
		arguments := []string{"run", "--rm", "--network", "host"}

		for key, value := range parameters {
			arguments = append(arguments, "-e")
			arguments = append(arguments, fmt.Sprintf("%s=%s", key, value))
		}
		arguments = append(arguments, "devatherock/simple-slack:latest")

		exitCode, output := helper.ExecuteCommand(exec.Command("docker", arguments...))
		fmt.Print(output)
		assert.Equal(test, 1, exitCode)

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
		assert.Equal(test, "Failure: https://someurl", attachment["text"])
		assert.Equal(test, "#cfd3d7", attachment["color"])
	}
}

func TestDynamicHighlightColor(test *testing.T) {
	// Test HTTP server
	var capturedRequest []byte
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		capturedRequest, _ = ioutil.ReadAll(request.Body)
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(writer, `{"success":true}`)
	}))
	defer testServer.Close()

	cases := []struct {
		buildStatus         string
		expectedColor       string
		ciVariable          string
		buildStatusVariable string
	}{
		{
			"success",
			"#33ad7f",
			"DRONE",
			"DRONE_BUILD_STATUS",
		},
		{
			"failure",
			"#a1040c",
			"DRONE",
			"DRONE_BUILD_STATUS",
		},
		{
			"error",
			"#a1040c",
			"DRONE",
			"DRONE_BUILD_STATUS",
		},
		{
			"killed",
			"#a1040c",
			"DRONE",
			"DRONE_BUILD_STATUS",
		},
		{
			"pending",
			"#cfd3d7",
			"DRONE",
			"DRONE_BUILD_STATUS",
		},
		{
			"success",
			"#33ad7f",
			"VELA",
			"VELA_BUILD_STATUS",
		},
		{
			"failure",
			"#a1040c",
			"VELA",
			"VELA_BUILD_STATUS",
		},
		{
			"error",
			"#a1040c",
			"VELA",
			"VELA_BUILD_STATUS",
		},
		{
			"pending",
			"#cfd3d7",
			"VELA",
			"VELA_BUILD_STATUS",
		},
	}

	for _, data := range cases {
		parameters := map[string]string{
			"TEXT":                   "Build completed",
			"WEBHOOK":                testServer.URL,
			data.buildStatusVariable: data.buildStatus,
			data.ciVariable:          "true",
		}
		arguments := []string{"run", "--rm", "--network", "host"}

		for key, value := range parameters {
			arguments = append(arguments, "-e")
			arguments = append(arguments, fmt.Sprintf("%s=%s", key, value))
		}
		arguments = append(arguments, "devatherock/simple-slack:latest")

		exitCode, output := helper.ExecuteCommand(exec.Command("docker", arguments...))
		fmt.Print(output)
		assert.Equal(test, 0, exitCode)

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
		assert.Equal(test, "Build completed", attachment["text"])
		assert.Equal(test, data.expectedColor, attachment["color"])
	}
}
