# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.1] - 2022-02-16

### Fixed
- Fix constant requeuing of messages (#21).

## [0.3.0] - 2021-02-25

### Added
- Add the `buneary get messages` command.

## [0.2.1] - 2021-02-19

### Changed
- Print success messages on successful resource creations.

## [0.2.0] - 2021-02-10

### Added
- Add the `buneary get exchanges` command.
- Add the `buneary get exchange` command.
- Add the `buneary get queues` command.
- Add the `buneary get queue` command.
- Add the `buneary get bindings` command.
- Add the `buneary get binding` command.
- Add the `--headers` option for specifying message headers.

### Changed
- Use the HTTP API port `15672` instead of the AMQP port `5672`.

## [0.1.1] - 2021-02-07

### Changed
- Enable support for `-v` flag, displaying version information.

### Fixed
- Fix help text for `buneary create binding` command.

## [0.1.0] - 2021-02-05

### First `buneary` release.

## [0.0.0]
