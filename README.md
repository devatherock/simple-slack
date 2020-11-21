[![CircleCI](https://circleci.com/gh/devatherock/simple-slack.svg?style=svg)](https://circleci.com/gh/devatherock/simple-slack)
[![Version](https://img.shields.io/docker/v/devatherock/simple-slack?sort=semver)](https://hub.docker.com/r/devatherock/simple-slack/)
[![Coverage Status](https://coveralls.io/repos/github/devatherock/simple-slack/badge.svg?branch=master)](https://coveralls.io/github/devatherock/simple-slack?branch=master)
[![Docker Pulls](https://img.shields.io/docker/pulls/devatherock/simple-slack.svg)](https://hub.docker.com/r/devatherock/simple-slack/)
[![Docker Image Size](https://img.shields.io/docker/image-size/devatherock/simple-slack.svg?sort=date)](https://hub.docker.com/r/devatherock/simple-slack/)
[![Docker Image Layers](https://img.shields.io/microbadger/layers/devatherock/simple-slack.svg)](https://microbadger.com/images/devatherock/simple-slack)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# simple-slack
CI plugin to post messages to [Slack](https://slack.com/) or other chat clients with Slack compatible incoming webhooks like [Rocket.Chat](https://rocket.chat/). For a listing of available options and  usage
samples, please take a look at the [docs](DOCS.md).

## Usage

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
  devatherock/simple-slack:0.4.0
```
