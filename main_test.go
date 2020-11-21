package main

import (
	"testing"
)

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
