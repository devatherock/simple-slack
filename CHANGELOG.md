# Changelog

## [Unreleased]
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
