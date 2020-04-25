[![CircleCI](https://circleci.com/gh/devatherock/simple-slack.svg?style=svg)](https://circleci.com/gh/devatherock/simple-slack)
[![Docker Pulls](https://img.shields.io/docker/pulls/devatherock/simple-slack.svg)](https://hub.docker.com/r/devatherock/simple-slack/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Microbadger](https://images.microbadger.com/badges/image/devatherock/simple-slack.svg)](https://microbadger.com/images/devatherock/simple-slack)
# simple-slack
CI plugin to post messages to slack. For a listing of available options and  usage
samples, please take a look at the [docs](DOCS.md).

## Usage

Execute from the working directory:

```
docker run --rm \
  -e SLACK_WEBHOOK=https://hooks.slack.com/services/... \
  -e PARAMETER_COLOR=#33ad7f \
  -e PARAMETER_TEXT="Success: {{.BuildLink}} ({{.BuildRef}}) by {{.BuildAuthor}}" \
  -e BUILD_REF="refs/heads/master" \
  -e BUILD_AUTHOR=octocat \
  -e BUILD_LINK=http://github.com/octocat/hello-world \
  devatherock/simple-slack:latest
```