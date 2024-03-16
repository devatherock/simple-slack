package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

// Presorted for contains check to work
var secretEnvVariables = []string{"PLUGIN_WEBHOOK", "SLACK_WEBHOOK", "WEBHOOK"}

const defaultColor string = "#cfd3d7" // grey
const successColor string = "#33ad7f" // green
const failureColor string = "#a1040c" // red

type SlackRequest struct {
	Text    string `json:",omitempty"`
	Channel string `json:",omitempty"`
	Color   string `json:",omitempty"`
	Title   string `json:",omitempty"`
	Webhook string `json:",omitempty"`
}

func Notify(request SlackRequest) error {
	err := validate(request)
	if err != nil {
		return err
	}

	payload, err := buildPayload(request)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", request.Webhook, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println("Message posted to webhook with http status", res.StatusCode)

	if res.StatusCode > 399 {
		return errors.New("HTTP request to Slack failed")
	}

	return nil
}

// Builds the Slack HTTP request payload
func buildPayload(request SlackRequest) (payload map[string]interface{}, err error) {
	text, err := parseTemplate(request.Text)
	if err != nil {
		return
	}

	// Build attachments section
	attachments := [1]map[string]string{
		{
			"color": getHighlightColor(request.Color),
			"text":  text,
		},
	}

	if request.Title != "" {
		attachments[0]["title"] = request.Title
	}

	// Build complete payload
	payload = map[string]interface{}{
		"attachments": attachments,
	}

	if request.Channel != "" {
		payload["channel"] = request.Channel
	}

	return
}

// Validates the input parameters
func validate(request SlackRequest) error {
	if request.Text == "" || request.Webhook == "" {
		return errors.New("Required parameters not specified")
	}

	return nil
}

// Decides the highlight color based on build status
func getHighlightColor(inputColor string) (outputColor string) {
	if inputColor != "" {
		outputColor = inputColor
	} else if os.Getenv("DRONE") == "true" {
		buildStatus := os.Getenv("DRONE_BUILD_STATUS")

		switch buildStatus {
		case "success":
			outputColor = successColor
		case "failure", "error", "killed":
			outputColor = failureColor
		default:
			outputColor = defaultColor
		}
	} else if os.Getenv("VELA") == "true" {
		buildStatus := os.Getenv("VELA_BUILD_STATUS")

		switch buildStatus {
		case "success", "running": // When none of the previous steps have failed, VELA_BUILD_STATUS has the value running within a step
			outputColor = successColor
		case "failure", "error":
			outputColor = failureColor
		default:
			outputColor = defaultColor
		}
	} else {
		outputColor = defaultColor
	}

	return
}

// Processes the input text as a template with environment variables as the
// context
func parseTemplate(templateText string) (string, error) {
	var templateContext = make(map[string]string)
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")

		// Inject all environment variables other than secrets
		if !contains(secretEnvVariables, variable[0]) {
			templateContext[envVariableToCamelCase(variable[0])] = variable[1]
		}
	}

	buffer := new(bytes.Buffer)
	parsedTemplate, err := template.New("test").Funcs(sprig.TxtFuncMap()).Parse(templateText)
	if err != nil {
		return "", err
	}

	err = parsedTemplate.Execute(buffer, templateContext)
	return buffer.String(), err
}

// Converts an environment variable name into a camelcase string. For example,
// BUILD_MESSAGE would be converted to BuildMessage
func envVariableToCamelCase(envVar string) (camelCase string) {
	isToUpper := true
	for _, runeValue := range envVar {
		if runeValue == '_' {
			isToUpper = true
		} else if isToUpper {
			camelCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			camelCase += strings.ToLower(string(runeValue))
		}
	}
	return
}

// Checks if a string is contained in a sorted string slice
func contains(stringArray []string, searchTerm string) bool {
	index := sort.SearchStrings(stringArray, searchTerm)
	return index < len(stringArray) && stringArray[index] == searchTerm
}
