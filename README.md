# Cohesion

[![License](https://img.shields.io/badge/license-Unlicense-blue.svg?style=flat)](UNLICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/danielkrainas/cohesion)](https://goreportcard.com/report/github.com/danielkrainas/cohesion)

Cohesion is a tool for automatically clustering nodes to one another. It will monitor the cluster agent on the node for when or if the node becomes disconnected or if it was never connected to start, Cohesion will use a preconfigured strategy for locating and joining the cluster or another stray node. The only supported node agent right now is [Hashicorp's Serf](https://github.com/hashicorp/serf).

## Installation

> $ go get github.com/danielkrainas/cohesion

## Usage

> $ cohesion [command] [config_path]

Most commands require a configuration path provided as an argument or in the `COHESION_CONFIG_PATH` environment variable. 

### Agent mode

This is the primary mode for Cohesion. It hosts the HTTP API server and handles monitoring and notifying hooks of container events.

> $ cohesion agent

**Example** - with development config:

> $ cohesion agent ./config.dev.yml

## Configuration

A configuration file is *required* for Cohesion but environment variables can be used to override configuration. A configuration file can be specified as a parameter or with the `COHESION_CONFIG_PATH` environment variable. 

All configuration environment variables are prefixed by `COHESION_` and the paths are separated by an underscore(`_`). Some examples:

- `COHESION_LOG_LEVEL=warn`
- `COHESION_LOG_FORMATTER=json`

A development configuration file is included: `/config.dev.yml` and a `/config.local.yml` has already been added to gitignore to be used for local testing or development.

```yaml
# configuration schema version number, only `0.1`
version: 0.1

# log stuff
log:
  # minimum event level to log: `error`, `warn`, `info`, or `debug`
  level: 'debug'
  # log output format: `text` or `json`
  formatter: 'text'
  # custom fields to be added and displayed in the log
  fields:
    customfield1: 'value'

# node monitoring settings
monitor:
  frequency: 10s

# node connection driver and parameters
node:
  serf:
    addr: 0.0.0.0

# discovery strategy and parameters
discovery: 'echo'

```

`node` and `discovery` only allow specification of *one* driver per configuration. Any additional ones will cause a validation error when the application starts.

## Bugs and Feedback

If you see a bug or have a suggestion, feel free to open an issue [here](https://github.com/danielkrainas/cohesion/issues).

## Contributions

PR's welcome! There are no strict style guidelines, just follow best practices and try to keep with the general look & feel of the code present. All submissions should atleast be `go fmt -s` and have a test to verify *(if applicable)*.

For details on how to extend and develop Cohesion, see the [dev documentation](docs/development/).

## License

[Unlicense](http://unlicense.org/UNLICENSE). This is a Public Domain work. 

[![Public Domain](https://licensebuttons.net/p/mark/1.0/88x31.png)](http://questioncopyright.org/promise)

> ["Make art not law"](http://questioncopyright.org/make_art_not_law_interview) -Nina Paley