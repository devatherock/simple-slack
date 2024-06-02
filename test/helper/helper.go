package helper

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sets an environment variable that will be cleaned up when the test ends
func SetEnvironmentVariable(test *testing.T, variable string, value string) {
	os.Setenv(variable, value)

	test.Cleanup(func() {
		os.Unsetenv(variable)
	})
}

func VerifySlackRequest(test *testing.T, request []byte, expected map[string]string) {
	jsonRequest := make(map[string]interface{})
	json.Unmarshal(request, &jsonRequest)
	assert.Equal(test, 2, len(jsonRequest))

	attachments := jsonRequest["attachments"].([]interface{})
	attachment := attachments[0].(map[string]interface{})

	assert.Equal(test, 1, len(attachments))
	assert.Equal(test, 3, len(attachment))
	assert.Equal(test, expected["text"], attachment["text"])
	assert.Equal(test, expected["color"], attachment["color"])
	assert.Equal(test, expected["title"], attachment["title"])
	assert.Equal(test, "general", jsonRequest["channel"])

}
