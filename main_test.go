package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func setEnvironmentVariable(test *testing.T, variable string, value string) {
	os.Setenv(variable, value)

	test.Cleanup(func() {
		os.Unsetenv(variable)
	})
}

func TestRunApp(test *testing.T) {
	setEnvironmentVariable(test, "TEXT", "Build failed!")

	// Test HTTP server
	var capturedRequest []byte
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		capturedRequest, _ = ioutil.ReadAll(request.Body)
		writer.WriteHeader(400)
	}))
	defer testServer.Close()
	setEnvironmentVariable(test, "WEBHOOK", testServer.URL)

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

func TestValidateError(test *testing.T) {
	cases := []struct {
		parameters map[string]string
	}{
		{
			map[string]string{
				"text": "hello",
			},
		},
		{
			map[string]string{
				"webhook": "https://secreturl",
			},
		},
	}
	expected := "Required parameters not specified"

	for _, data := range cases {
		set := flag.NewFlagSet("test", 0)
		for key, value := range data.parameters {
			set.String(key, value, "dummy")
		}

		context := cli.NewContext(nil, set, nil)
		actual := validate(context)

		assert.Equal(test, expected, actual.Error())
	}
}

func TestValidateSuccess(test *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("text", "hello", "dummy")
	set.String("webhook", "https://secreturl", "dummy")

	context := cli.NewContext(nil, set, nil)
	actual := validate(context)

	assert.Nil(test, actual)
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

func TestBuildPayload(test *testing.T) {
	cases := []struct {
		parameters map[string]string
		expected   map[string]interface{}
	}{
		{
			map[string]string{
				"text":    "Build failed!",
				"color":   "red",
				"title":   "Build notification",
				"channel": "general",
			},
			map[string]interface{}{
				"attachments": [1]map[string]string{
					{
						"color": "red",
						"text":  "Build failed!",
						"title": "Build notification",
					},
				},
				"channel": "general",
			},
		},
		{
			map[string]string{
				"text": "Build failed!",
			},
			map[string]interface{}{
				"attachments": [1]map[string]string{
					{
						"color": "#cfd3d7",
						"text":  "Build failed!",
					},
				},
			},
		},
	}

	for _, data := range cases {
		set := flag.NewFlagSet("test", 0)
		for key, value := range data.parameters {
			set.String(key, value, "")
		}

		context := cli.NewContext(nil, set, nil)
		actual := buildPayload(context)

		assert.Equal(test, data.expected, actual)
	}
}

func TestGetHighlightColorForDrone(test *testing.T) {
	cases := []struct{ buildStatus, inputColor, expected string }{
		{"success", "yellow", "yellow"},
		{"failure", "yellow", "yellow"},
		{"success", "", "#33ad7f"},
		{"failure", "", "#a1040c"},
		{"error", "", "#a1040c"},
		{"killed", "", "#a1040c"},
		{"pending", "", "#cfd3d7"},
	}

	for _, data := range cases {
		setEnvironmentVariable(test, "DRONE", "true")
		setEnvironmentVariable(test, "DRONE_BUILD_STATUS", data.buildStatus)
		actual := getHighlightColor(data.inputColor)

		assert.Equal(test, data.expected, actual)
	}
}

func TestGetHighlightColorForVela(test *testing.T) {
	cases := []struct{ buildStatus, inputColor, expected string }{
		{"success", "yellow", "yellow"},
		{"failure", "yellow", "yellow"},
		{"success", "", "#33ad7f"},
		{"failure", "", "#a1040c"},
		{"error", "", "#a1040c"},
		{"killed", "", "#cfd3d7"},
		{"pending", "", "#cfd3d7"},
	}

	for _, data := range cases {
		setEnvironmentVariable(test, "VELA", "true")
		setEnvironmentVariable(test, "VELA_BUILD_STATUS", data.buildStatus)
		actual := getHighlightColor(data.inputColor)

		assert.Equal(test, data.expected, actual)
	}
}

func TestGetHighlightColorForOtherCI(test *testing.T) {
	cases := []struct{ inputColor, expected string }{
		{"yellow", "yellow"},
		{"", "#cfd3d7"},
	}

	for _, data := range cases {
		actual := getHighlightColor(data.inputColor)
		assert.Equal(test, data.expected, actual)
	}
}

func TestParseTemplate(test *testing.T) {
	cases := []struct{ template, expected string }{
		{"Success: {{.CircleBuildUrl}}, Slack URL: {{.Webhook}}", "Success: https://someurl, Slack URL: <no value>"},
		{"Success: {{.CircleBuildUrl}}", "Success: https://someurl"},
	}

	for _, data := range cases {
		setEnvironmentVariable(test, "CIRCLE_BUILD_URL", "https://someurl")
		setEnvironmentVariable(test, "WEBHOOK", "https://secreturl")
		actual := parseTemplate(data.template)

		assert.Equal(test, data.expected, actual)
	}
}

func TestEnvVariableToCamelCase(test *testing.T) {
	cases := []struct{ inputVariable, expected string }{
		{"BUILD_STATUS", "BuildStatus"},
		{"DRONE", "Drone"},
	}

	for _, data := range cases {
		actual := envVariableToCamelCase(data.inputVariable)
		assert.Equal(test, data.expected, actual)
	}
}

func TestContains(test *testing.T) {
	secretsArray := []string{"PLUGIN_WEBHOOK", "SLACK_WEBHOOK", "WEBHOOK"}
	cases := []struct {
		envVariable string
		expected    bool
	}{
		{"PLUGIN_WEBHOOK", true},
		{"SLACK_WEBHOOK", true},
		{"WEBHOOK", true},
		{"PARAMETER_TEXT", false},
		{"PARAMETER_TITLE", false},
		{"PARAMETER_CHANNEL", false},
	}

	for _, data := range cases {
		actual := contains(secretsArray, data.envVariable)
		assert.Equal(test, data.expected, actual)
	}
}
