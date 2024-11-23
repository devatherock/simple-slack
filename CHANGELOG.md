# Changelog

## [Unreleased]
### Changed
- Used image from dockerhub for deployment
- fix(deps): update module github.com/urfave/cli/v2 to v2.27.5
- fix(deps): update module github.com/stretchr/testify to v1.10.0

## [1.3.0] - 2024-09-22
### Added
- [#78](https://github.com/devatherock/simple-slack/issues/78): Included CircleCI workflow name in message

### Changed
- Used machine executor in deployment step
- Upgraded `flyctl` to `0.2.93`

## [1.2.0] - 2024-09-22
### Added
- Health check to docker compose, to prevent transient build failures

### Changed
- Fixed bug in deployment step
- Upgraded go to `1.22`
- Upgraded `circleci-templates` to `0.7.0`
- Upgraded `golang.org/x/crypto` to `0.21.0`
- chore(deps): update alpine docker tag to v3.20.3
- [#73](https://github.com/devatherock/simple-slack/issues/73): Handled `failing` CircleCI status

## [1.1.0] - 2024-06-07
### Added
- `circleci-templates` orb for common tasks

### Changed
- Made only HIGH bolt vulnerabilities create issues
- Upgraded go to `1.20`
- chore(deps): update alpine docker tag to v3.19.1
- fix(deps): update module github.com/stretchr/testify to v1.9.0
- Moved functional logic away from the main file
- feat: Added endpoint to send notification for CircleCI builds
- chore: Added fly.io deployment configuration
- fix(deps): update module github.com/urfave/cli/v2 to v2.27.2

## [1.0.0] - 2023-06-16
### Added
- [#30](https://github.com/devatherock/simple-slack/issues/30): functional tests

### Changed
- Updated dockerhub readme in CI pipeline
- Upgraded `go` to `1.18`
- Set alpine version to `3.17.4`
- [#39](https://github.com/devatherock/simple-slack/issues/39): Built a multi-arch docker image

## [0.7.0] - 2021-04-06
### Changed
- feat: Quit with non-zero exit code when API call to Slack fails([#35](https://github.com/devatherock/simple-slack/issues/35))

## [0.6.0] - 2021-02-14
### Added
- feat: Used `VELA_BUILD_STATUS` environment variable to choose message highlight color in vela
- feat: Added support for sprig functions within the text template([#32](https://github.com/devatherock/simple-slack/issues/32))

## [0.5.0] - 2020-11-22
### Added
- make file
- First unit test
- Code coverage using coveralls
- [#10](https://github.com/devatherock/simple-slack/issues/10): Unit tests

### Changed
-   Refactored code for easier unit testing

## [0.4.0] - 2020-06-13
### Added
- Some additional fields to the outgoing webhook payload. This is to support Zulip messaging.

## [0.3.0] - 2020-06-12
### Added
- A log statement on successful completion

### Changed
- Upgraded to `urfave/cli/v2`

## [0.2.0] - 2020-04-24
### Changed
- [Issue 4](https://github.com/devatherock/simple-slack/issues/4): Highlight color based on build status

## [0.1.0] - 2020-04-24
### Added
- Support for templating with environment variables injected in camel case

## [0.0.1] - 2020-04-24
### Added
- Initial version. Posts provided text to slack, in specified color