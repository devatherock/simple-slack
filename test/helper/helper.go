package helper

import (
	"os"
	"testing"
)

// Sets an environment variable that will be cleaned up when the test ends
func SetEnvironmentVariable(test *testing.T, variable string, value string) {
	os.Setenv(variable, value)

	test.Cleanup(func() {
		os.Unsetenv(variable)
	})
}
