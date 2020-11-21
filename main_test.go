package main

import (
	"flag"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func setEnvironmentVariable(test *testing.T, variable string, value string) {
	os.Setenv(variable, value)

	test.Cleanup(func() {
		os.Unsetenv(variable)
	})
}

func TestValidateError(test *testing.T) {
	cases := []struct {
		parameters map[string]string
	}{
		{map[string]string{
			"text": "hello",
		}},
		{map[string]string{
			"webhook": "https://secreturl",
		}},
	}
	expected := "Required parameters not specified"

	for _, data := range cases {
		set := flag.NewFlagSet("test", 0)
		for key, value := range data.parameters {
			set.String(key, value, "dummy")
		}

		context := cli.NewContext(nil, set, nil)
		actual := validate(context)
		if actual.Error() != expected {
			test.Logf("Expected: %s, Actual: %s", expected, actual)
			test.Fail()
		}
	}
}

func TestValidateSuccess(test *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("text", "hello", "dummy")
	set.String("webhook", "https://secreturl", "dummy")

	context := cli.NewContext(cli.NewApp(), set, nil)
	actual := validate(context)
	if actual != nil {
		test.Logf("Expected: nil, Actual: %s", actual)
		test.Fail()
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
		if actual != data.expected {
			test.Logf("Expected: %s, Actual: %s", data.expected, actual)
			test.Fail()
		}
	}
}

func TestGetHighlightColorForCIOtherThanDrone(test *testing.T) {
	cases := []struct{ inputColor, expected string }{
		{"yellow", "yellow"},
		{"", "#cfd3d7"},
	}

	for _, data := range cases {
		actual := getHighlightColor(data.inputColor)
		if actual != data.expected {
			test.Logf("Expected: %s, Actual: %s", data.expected, actual)
			test.Fail()
		}
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
		if actual != data.expected {
			test.Logf("Expected: %s, Actual: %s", data.expected, actual)
			test.Fail()
		}
	}
}

func TestEnvVariableToCamelCase(test *testing.T) {
	cases := []struct{ inputVariable, expected string }{
		{"BUILD_STATUS", "BuildStatus"},
		{"DRONE", "Drone"},
	}

	for _, data := range cases {
		actual := envVariableToCamelCase(data.inputVariable)
		if actual != data.expected {
			test.Logf("Expected: %s, Actual: %s", data.expected, actual)
			test.Fail()
		}
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
		if actual != data.expected {
			test.Logf("Expected: %t, Actual: %t", data.expected, actual)
			test.Fail()
		}
	}
}
