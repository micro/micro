# Micro [![License](https://img.shields.io/badge/license-polyform:shield-blue)](https://polyformproject.org/licenses/shield/1.0.0/) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)

<kbd><img src="https://raw.githubusercontent.com/micro/micro/master/docs/images/banner.png" /></kbd>

## Overview

Micro addresses the key requirements for building services in the cloud. It leverages the microservices
architecture pattern and provides a set of services which act as the building blocks of a platform. Micro deals
with the complexity of distributed systems and provides simpler programmable abstractions to build on. 

## Contents

- [Introduction](https://micro.mu/introduction) - A high level introduction to Micro
- [Getting Started](https://micro.mu/getting-started) - The helloworld quickstart guide
- [Upgrade Guide](https://micro.mu/upgrade-guide) - Update your go-micro project to use micro v3.
- [Architecture](https://micro.mu/architecture) - Describes the architecture, design and tradeoffs
- [Reference](https://micro.mu/reference) - In-depth reference for Micro CLI and services
- [Resources](https://micro.mu/resources) - External resources and contributions
- [Roadmap](https://micro.mu/roadmap) - Stuff on our agenda over the long haul
- [Users](https://micro.mu/users) - Developers and companies using Micro in production
- [FAQ](https://micro.mu/faq) - Frequently asked questions
- [Blog](https://micro.mu/blog) - For the latest from us

## Features

Below are the core components that make up Micro.

**Server**

Micro is built as a microservices architecture and abstracts away the complexity of the underlying infrastructure. We compose 
this as a single logical server to the user but decompose that into the various building block primitives that can be plugged 
into any underlying system. 

The server is composed of the following services.

- **API** - HTTP Gateway which dynamically maps http/json requests to RPC using path based resolution
- **Auth** - Authentication and authorization out of the box using jwt tokens and rule based access control.
- **Broker** - Ephemeral pubsub messaging for asynchronous communication and distributing notifications
- **Config** - Dynamic configuration and secrets management for service level config without the need to restart
- **Events** - Event streaming with ordered messaging, replay from offsets and persistent storage
- **Network** - Inter-service networking, isolation and routing plane for all internal request traffic
- **Proxy** - An identity aware proxy used for remote access and any external grpc request traffic
- **Runtime** - Service lifecyle and process management with support for source to running auto build
- **Registry** - Centralised service discovery and API endpoint explorer with feature rich metadata
- **Store** - Key-Value storage with TTL expiry and persistent crud to keep microservices stateless

**Framework**

Micro additionaly now contains the incredibly popular Go Micro framework built in for service development. 
The Go framework makes it drop dead simple to write your services without having to piece together lines and lines of boilerplate. Auto 
configured and initialised by default, just import and get started quickly.

**Command Line**

Micro brings not only a rich architectural model but a command line experience tailored for that need. The command line interface includes 
dynamic command mapping for all services running on the platform. Turns any service instantly into a CLI command along with flag parsing 
for inputs. Includes support for multiple environments and namespaces, automatic refreshing of auth credentials, creating and running 
services, status info and log streaming, plus much, much more.

**Environments**

Finally Micro bakes in the concept of `Environments` and multi-tenancy through `Namespaces`. Run your server locally for 
development and in the cloud for staging and production, seamlessly switch between them using the CLI commands `micro env set [environment]` 
and `micro user set [namespace]`.

## Install

**From Source**

```sh
go get github.com/micro/micro/v3
```

**Using Docker**

```sh
# install
docker pull micro/micro

# run it
docker run -p 8080-8081:8080-8081/tcp micro/micro server
```

**Helm Chart**

```
helm repo add micro https://micro.github.io/helm
helm install micro micro/micro
```

**Release binaries**

```sh
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

Set the environment to local (127.0.0.1:8081)

```
micro env set local
```

Login to the server

```
# user: admin pass: micro
micro login
```

Create a service

```sh
# generate a service (follow instructions in output)
micro new helloworld

# run the service
micro run helloworld

# check the status
micro status

# list running services
micro services

# call the service
micro helloworld --name=Alice

# curl via the api
curl -d '{"name": "Alice"}' http://localhost:8080/helloworld
```

## Example Service

Micro includes a Go framework for writing services wrapping gRPC for the core IDL and transport. 

Define services in proto:

```proto
syntax = "proto3";

package helloworld;

service Helloworld {
	rpc Call(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}
```

Write them using Go:

```go
package main

import (
	"context"
  
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/helloworld/proto"
)

type Helloworld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (h *Helloworld) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logger.Info("Received Helloworld.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	// Create service
	srv := service.New(
		service.Name("helloworld"),
	)

	// Register Handler
	srv.Handle(new(Helloworld))

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
```

Call with the client:

```go
import (
	"context"
  
	"github.com/micro/micro/v3/service/client"
	pb "github.com/micro/services/helloworld/proto"
)

// create a new helloworld service client
helloworld := pb.NewHelloworldService("helloworld", client.DefaultClient) 

// call the endpoint Helloworld.Call
rsp, err := helloworld.Call(context.Background(), &pb.Request{Name: "Alice"})
```

## Usage

See the [docs](https://micro.mu) for detailed information on the architecture, installation and use of the platform.

## License

See [LICENSE](LICENSE) which makes use of [Polyform Shield](https://polyformproject.org/licenses/shield/1.0.0/).

## Hosting

For the hosted Micro Platform aka M3O see [m3o.com](https://m3o.com).

## Community

Join us in GitHub [Discussions](https://github.com/micro/micro/discussions), [Slack](https://slack.micro.mu) or follow on [Twitter](https://twitter.com/microhq) for updates.
