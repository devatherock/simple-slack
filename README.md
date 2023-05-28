[![CircleCI](https://circleci.com/gh/devatherock/simple-slack.svg?style=svg)](https://circleci.com/gh/devatherock/simple-slack)
[![Version](https://img.shields.io/docker/v/devatherock/simple-slack?sort=semver)](https://hub.docker.com/r/devatherock/simple-slack/)
[![Coverage Status](https://coveralls.io/repos/github/devatherock/simple-slack/badge.svg?branch=master)](https://coveralls.io/github/devatherock/simple-slack?branch=master)
[![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=simple-slack&metric=alert_status)](https://sonarcloud.io/component_measures?id=simple-slack&metric=alert_status&view=list)
[![Docker Pulls](https://img.shields.io/docker/pulls/devatherock/simple-slack.svg)](https://hub.docker.com/r/devatherock/simple-slack/)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=simple-slack&metric=ncloc)](https://sonarcloud.io/component_measures?id=simple-slack&metric=ncloc)
[![Docker Image Size](https://img.shields.io/docker/image-size/devatherock/simple-slack.svg?sort=date)](https://hub.docker.com/r/devatherock/simple-slack/)
# simple-slack
CI plugin to post messages to [Slack](https://slack.com/) or other chat clients with Slack compatible incoming webhooks like [Rocket.Chat](https://rocket.chat/)

## Config

The following parameters/secrets can be set to configure the plugin.

### Parameters
* **color** - Color in which the message block will be highlighted.
* **text** - The message content. The text uses go templating. Any environment variable available at runtime can be used within the text, after converting it to camel case. For example, to use the environment variable `DRONE_BUILD_STATUS`, the syntax will be `{{.DroneBuildStatus}}`

### Secrets

The following secret values can be set to configure the plugin.

* **SLACK_WEBHOOK** - The slack webhook to post the message to

## Usage

### Docker:
Execute from the working directory:

```
docker run --rm \
  -e SLACK_WEBHOOK=https://hooks.slack.com/services/... \
  -e PARAMETER_COLOR=#33ad7f \
  -e PARAMETER_TEXT="Success: {{.BuildLink}} ({{.BuildRef}}) by {{.BuildAuthor}}" \
  -e PARAMETER_CHANNEL="xyz" \
  -e PARAMETER_TITLE="Build completed" \
  -e BUILD_REF="refs/heads/master" \
  -e BUILD_AUTHOR=octocat \
  -e BUILD_LINK=http://github.com/octocat/hello-world \
  devatherock/simple-slack:latest
```

### Drone:

```yaml
pipeline:
  notify_slack:
    when:
      event: [ push ]
      status: [ success, failure, error ]
    image: devatherock/simple-slack:latest
    secrets: [ slack_webhook ]
    settings:
      color: "#33ad7f"
      text: |-
        {{.DroneBuildStatus}} {{.DroneBuildLink}} ({{.DroneCommitRef}}) by {{DroneCommitAuthor}}
        {{.DroneCommitMessage}}
```

### Vela:

```yaml
steps:
  - name: notify_slack
    ruleset:
      event: [ push ]
      status: [ success ]
    image: devatherock/simple-slack:latest
    secrets: [ slack_webhook ]
    parameters:
      color: "#33ad7f"
      text: |-
        Success: {{.BuildLink}} ({{.BuildRef}}) by {{.BuildAuthor}}
        {{.BuildMessage}}
```
