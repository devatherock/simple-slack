package main

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

	"github.com/urfave/cli/v2"
)

// Presorted for contains check to work
var secretEnvVariables = []string{"PLUGIN_WEBHOOK", "SLACK_WEBHOOK", "WEBHOOK"}

const defaultColor string = "#cfd3d7" // grey

func main() {
	runApp(os.Args)
}

// Initializes and runs the app
func runApp(args []string) {
	app := cli.NewApp()
	app.Name = "simple slack plugin"
	app.Before = validate
	app.Action = run
	app.Flags = []cli.Flag{
		createStringCliFlag(
			"color",
			[]string{"c"},
			"Color in which the message block will be highlighted",
			[]string{"COLOR", "PLUGIN_COLOR", "PARAMETER_COLOR"},
		),
		createStringCliFlag(
			"text",
			[]string{"t"},
			"The message content",
			[]string{"TEXT", "PLUGIN_TEXT", "PARAMETER_TEXT"},
		),
		createStringCliFlag(
			"title",
			[]string{"ti"},
			"The message title",
			[]string{"TITLE", "PLUGIN_TITLE", "PARAMETER_TITLE"},
		),
		createStringCliFlag(
			"channel",
			[]string{"ch"},
			"The slack channel name",
			[]string{"CHANNEL", "PLUGIN_CHANNEL", "PARAMETER_CHANNEL"},
		),
		createStringCliFlag(
			"webhook",
			[]string{"u"},
			"The slack webhook URL",
			[]string{"WEBHOOK", "PLUGIN_WEBHOOK", "SLACK_WEBHOOK"},
		),
	}

	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

// Creates a String CLI parameter
func createStringCliFlag(name string, aliases []string, usage string, envVars []string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    name,
		Aliases: aliases,
		Usage:   usage,
		EnvVars: envVars,
	}
}

// Validates the input parameters
func validate(context *cli.Context) error {
	if context.String("text") == "" || context.String("webhook") == "" {
		return errors.New("Required parameters not specified")
	}

	return nil
}

// Sends the input text to slack
func run(context *cli.Context) error {
	payload := buildPayload(context)
	data, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", context.String("webhook"), bytes.NewBuffer(data))
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

	return nil
}

// Builds the Slack HTTP request payload
func buildPayload(context *cli.Context) (payload map[string]interface{}) {
	// Build attachments section
	text := parseTemplate(context.String("text"))
	attachments := [1]map[string]string{
		{
			"color": getHighlightColor(context.String("color")),
			"text":  text,
		},
	}

	title := context.String("title")
	if title != "" {
		attachments[0]["title"] = title
	}

	// Build complete payload
	payload = map[string]interface{}{
		"attachments": attachments,
	}

	channel := context.String("channel")
	if channel != "" {
		payload["channel"] = channel
	}

	return
}

// Decides the highlight color based on build status
func getHighlightColor(inputColor string) (outputColor string) {
	if inputColor != "" {
		outputColor = inputColor
	} else if os.Getenv("DRONE") == "true" {
		buildStatus := os.Getenv("DRONE_BUILD_STATUS")

		switch buildStatus {
		case "success":
			outputColor = "#33ad7f" // green
		case "failure", "error", "killed":
			outputColor = "#a1040c" // red
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
func parseTemplate(templateText string) string {
	var templateContext = make(map[string]string)
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")

		// Inject all environment variables other than secrets
		if !contains(secretEnvVariables, variable[0]) {
			templateContext[envVariableToCamelCase(variable[0])] = variable[1]
		}
	}

	buffer := new(bytes.Buffer)
	parsedTemplate, err := template.New("test").Parse(templateText)
	err = parsedTemplate.Execute(buffer, templateContext)
	if err != nil {
		log.Fatal(err)
	}

	return buffer.String()
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
