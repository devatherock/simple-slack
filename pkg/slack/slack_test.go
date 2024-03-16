//go:build test
// +build test

package slack

import (
	"testing"

	"github.com/devatherock/simple-slack/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestValidateError(test *testing.T) {
	cases := []struct {
		request SlackRequest
	}{
		{
			SlackRequest{
				Text: "hello",
			},
		},
		{
			SlackRequest{
				Webhook: "https://secreturl",
			},
		},
	}
	expected := "Required parameters not specified"

	for _, data := range cases {
		actual := validate(data.request)

		assert.Equal(test, expected, actual.Error())
	}
}

func TestValidateSuccess(test *testing.T) {
	request := SlackRequest{
		Text:    "hello",
		Webhook: "https://secreturl",
	}
	actual := validate(request)

	assert.Nil(test, actual)
}

func TestBuildPayload(test *testing.T) {
	cases := []struct {
		request  SlackRequest
		expected map[string]interface{}
	}{
		{
			SlackRequest{
				Text:    "Build failed!",
				Color:   "red",
				Title:   "Build notification",
				Channel: "general",
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
			SlackRequest{
				Text: "Build failed!",
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
		actual, err := buildPayload(data.request)

		assert.Nil(test, err)
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
		helper.SetEnvironmentVariable(test, "DRONE", "true")
		helper.SetEnvironmentVariable(test, "DRONE_BUILD_STATUS", data.buildStatus)
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
		helper.SetEnvironmentVariable(test, "VELA", "true")
		helper.SetEnvironmentVariable(test, "VELA_BUILD_STATUS", data.buildStatus)
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
		{
			"Success: {{.CircleBuildUrl}}, Slack URL: {{.Webhook}}",
			"Success: https://someurl, Slack URL: <no value>",
		},
		{
			"Success: {{.CircleBuildUrl}}",
			"Success: https://someurl",
		},
	}

	for _, data := range cases {
		helper.SetEnvironmentVariable(test, "CIRCLE_BUILD_URL", "https://someurl")
		helper.SetEnvironmentVariable(test, "WEBHOOK", "https://secreturl")
		actual, err := parseTemplate(data.template)

		assert.Nil(test, err)
		assert.Equal(test, data.expected, actual)
	}
}

func TestParseSprigTemplate(test *testing.T) {
	cases := []struct{ template, expected string }{
		{
			`Success: {{list .CircleBuildUrl ", Slack URL: " .Webhook | join "" }}`,
			"Success: https://someurl, Slack URL: ",
		},
	}

	for _, data := range cases {
		helper.SetEnvironmentVariable(test, "CIRCLE_BUILD_URL", "https://someurl")
		helper.SetEnvironmentVariable(test, "WEBHOOK", "https://secreturl")
		actual, err := parseTemplate(data.template)

		assert.Nil(test, err)
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
