package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "simple slack plugin"
	app.Before = validate
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "color, c",
			Usage:  "Color with which the message block will be highlighted",
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

func validate(context *cli.Context) error {
	if context.String("text") == "" || context.String("webhook") == "" {
		cli.ShowAppHelp(context)
		os.Exit(1)
	}

	return nil
}

func run(context *cli.Context) {
	attachments := [1]map[string]string{
		{
			"color": context.String("color"),
			"text":  context.String("text"),
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
