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
	}{
		{
			"TEXT",
			"WEBHOOK",
		},
		{
			"PLUGIN_TEXT",
			"PLUGIN_WEBHOOK",
		},
		{
			"PARAMETER_TEXT",
			"SLACK_WEBHOOK",
		},
	}

	for _, data := range cases {
		parameters := map[string]string{
			data.textVariable:    "Build failed!",
			data.webhookVariable: testServer.URL,
		}
		arguments := []string{"run", "--rm", "--network", "host"}

		for key, value := range parameters {
			arguments = append(arguments, "-e")
			arguments = append(arguments, fmt.Sprintf("%s=%s", key, value))
		}
		arguments = append(arguments, "devatherock/simple-slack:latest")

		exitCode, output := helper.ExecuteCommand(exec.Command("docker", arguments...))
		assert.Equal(test, 0, exitCode)
		fmt.Print(output)

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
}
