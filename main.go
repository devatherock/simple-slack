package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

// Presorted for contains check to work
var secretEnvVariables = []string{"PLUGIN_WEBHOOK", "SLACK_WEBHOOK", "WEBHOOK"}

func main() {
	app := cli.NewApp()
	app.Name = "simple slack plugin"
	app.Before = validate
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "color, c",
			Usage:  "Color in which the message block will be highlighted",
			EnvVar: "COLOR,PLUGIN_COLOR,PARAMETER_COLOR",
			Value:  "#cfd3d7",
		},
		cli.StringFlag{
			Name:   "text, t",
			Usage:  "The message content",
			EnvVar: "TEXT,PLUGIN_TEXT,PARAMETER_TEXT",
		},
		cli.StringFlag{
			Name:   "webhook, u",
			Usage:  "The slack webhook URL",
			EnvVar: "WEBHOOK,PLUGIN_WEBHOOK,SLACK_WEBHOOK",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Validates the input parameters
func validate(context *cli.Context) error {
	if context.String("text") == "" || context.String("webhook") == "" {
		cli.ShowAppHelp(context)
		os.Exit(1)
	}

	return nil
}

// Sends the input text to slack
func run(context *cli.Context) {
	attachments := [1]map[string]string{
		{
			"color": context.String("color"),
			"text":  parseTemplate(context.String("text")),
		},
	}
	payload := map[string][1]map[string]string{
		"attachments": attachments,
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", context.String("webhook"), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
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
