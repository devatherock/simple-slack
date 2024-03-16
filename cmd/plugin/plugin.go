package main

import (
	"log"
	"os"

	"github.com/devatherock/simple-slack/pkg/slack"
	"github.com/urfave/cli/v2"
)

func main() {
	runApp(os.Args)
}

// Initializes and runs the app
func runApp(args []string) {
	app := cli.NewApp()
	app.Name = "simple slack plugin"
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
	handleError(err)
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

// Sends the input text to slack
func run(context *cli.Context) error {
	return slack.Notify(buildRequest(context))
}

// Forms a Slack request from the supplied parameters
func buildRequest(context *cli.Context) slack.SlackRequest {
	slackRequest := slack.SlackRequest{}
	slackRequest.Text = context.String("text")
	slackRequest.Color = context.String("color")
	slackRequest.Title = context.String("title")
	slackRequest.Channel = context.String("channel")
	slackRequest.Webhook = context.String("webhook")

	return slackRequest
}

// Logs the error and exits the application
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
