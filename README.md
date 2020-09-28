# Micro [![License](https://img.shields.io/badge/license-polyform:shield-blue)](https://polyformproject.org/licenses/shield/1.0.0/) [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/micro/micro/v3?tab=overview) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro) [<img src="https://img.shields.io/badge/slack-micro-yellow.svg?logo=slack" />](https://slack.micro.mu)

Micro is a Go cloud services development framework.

## Overview

Micro addresses the key requirements for building cloud native services. It leverages the microservices
architecture pattern and provides a set of services which act as the building blocks of a platform. Micro deals
with the complexity of distributed systems and provides simpler programmable abstractions to build on. 

## Features

The framework is composed of the following features:

- **Server:** A distributed systems runtime composed of building block services which abstract away the underlying infrastructure 
and provide a programmable abstraction layer. Authentication, configuration, messaging, storage and more built in.

- **Clients:** Multiple entrypoints through which you can access your services. Write services once and access them through every means 
you've already come to know. A HTTP api, gRPC proxy and commmand line interface.

- **Library:** A Go library which makes it drop dead simple to write your services without having to piece together lines and lines of 
boilerplate. Auto configured and initialised by default, just import and get started quickly.

- **Plugins:** Micro is runtime and infrastructure agnostic. Each underlying building block service uses the Go Micro standard library 
to provide a pluggable foundation. We make it simple to use by pre-initialising for local use and the cloud.

## Install

Install from source

```
go get github.com/micro/micro/v3
```

Using a docker image

```
docker pull micro/micro
```

Latest release binaries

```
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

## Getting Started

Run the server locally

```
micro server
```

Create a service

```
# generate a service (follow instructions in output)
micro new helloworld

# run the service
micro run helloworld

# list services
micro services

# call a service
micro helloworld --name=Alice

# curl via the api
curl -d '{"name": "Alice"}' http://localhost:8080/helloworld
```

## Usage

See all the options

```
micro --help
```

See the [docs](https://github.com/micro/docs) for detailed information on the architecture, installation and use of the platform.

## License

See [LICENSE](LICENSE) which makes use of [Polyform Shield](https://polyformproject.org/licenses/shield/1.0.0/).

## Hosting

If you're interested in a hosted version of Micro see [m3o.com](https://m3o.com). Docs at [m3o.dev](https://m3o.dev).

## Commercial Use

If you want to sell or offer Micro as a Service please email [contact@m3o.com](mailto:contact@m3o.com)

## Community

Join us on [slack](https://slack.micro.mu) or follow us on [Twitter](https://twitter.com/microhq)
