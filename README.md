# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)

Micro is a distributed systems runtime for the Cloud and beyond.

## Overview

Micro addresses the key requirements for building distributed systems in the Cloud and beyond. It leverages the microservices
architecture pattern and provides a set of services which act as the building blocks of a platform. Micro deals
with the complexity of distributed systems and provides simpler programmable abstractions to build on.

<img src="https://micro.mu/runtime01.svg?load" />

Technology is constantly evolving. The infrastructure stack is always changing. Micro is a platform which
addresses these issues with a pluggable foundation and strongly defined apis to build on. Plug into any stack or cloud.

Follow us on [Twitter](https://twitter.com/microhq) or join the [Slack](https://micro.mu/slack) community.

## Features

The runtime is composed of the following features:

- **api:** An api gateway. A single entry point with dynamic request routing using service discovery. The API gateway allows you to build a scalable
microservice architecture on the backend and consolidate serving a public api on the frontend. The micro api provides powerful routing
via discovery and pluggable handlers to serve http, grpc, websockets, publish events and more.

- **broker:** A message broker allowing for async messaging. Microservices are event driven architectures and should provide messaging as a first
class citizen. Notify other services of events without needing to worry about a response.

- **network:** Build multi-cloud networks with the micro network service. Simply drop-in and connect the network services across any environment
and create a single flat network to route globally. The micro network dynamically builds routes based on your local registry in each datacenter
ensuring queries are routed based on locality.

- **new:** A service template generator. Create new service templates to get started quickly. Micro provides predefined templates for writing micro services.
Always start in the same way, build identical services to be more productive.

- **proxy:** A transparent service proxy built on [Go Micro](https://github.com/micro/go-micro). Offload service discovery, load balancing,
fault tolerance, message encoding, middleware, monitoring and more to a single a location. Run it standalone or alongside your service.

- **registry:** The registry provides service discovery to locate other services, store feature rich metadata and endpoint information. It's a
service explorer which lets you centrally and dynamically store this info at runtime.

- **store:** State is a fundamental requirement of any system. We provide a key-value store to provide simple storage of state which can be shared
between services or offload long term to keep microservices stateless and horizontally scalable.

- **web:** The web dashboard allows you to explore your services, describe their endpoints, the request and response formats and even
query them directly. The dashboard also includes a built in CLI like experience for developers who want to drop into the terminal on the fly.

Additionally micro provides a Go development framework:

- **go-micro:** Leverage the powerful [Go Micro](https://github.com/micro/go-micro) framework to develop microservices easily and quickly.
Go Micro abstracts away the complexity of distributed systems and provides simpler abstractions to build highly scalable microservices.

## Install

From source

```
go get github.com/micro/micro/v2
```

Docker image

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

Boot the entire runtime environment locally

```
micro server
```

### Create a service

```
# generate a service (follow instructions in output)
micro new example

# run the service
micro run --server example

# list services
micro list services

# call a service
micro call go.micro.service.example Example.Call '{"name": "John"}'
```

## Usage

See all the options

```
micro --help
```

See the [docs](https://micro.mu/docs/runtime.html) for detailed information on the architecture, installation and use of the platform.
