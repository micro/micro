# Micro

<p align="center">
    <a href="https://discord.gg/TBR9bRjd6Z"><img src="https://img.shields.io/badge/chat-discord-brightgreen.svg?logo=discord&%20style=flat"></a>
    <a href="https://goreportcard.com/report/github.com/micro/micro">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/micro/micro">
    </a>
	<a href="https://pkg.go.dev/github.com/micro/micro/v3?tab=doc"><img
    alt="Go.Dev reference"
    src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white"></a>
    <a href="https://opensource.org/licenses/Apache-2.0"><img
    alt="Apache License"
    src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"></a>
</p>

An API first development platform

## Overview 

Micro addresses the key requirements for building services in the cloud. It leverages the microservices
architecture pattern and provides a set of services which act as the building blocks of a platform. 

Micro deals with the complexity of distributed systems and provides simpler programmable abstractions to build on. 

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
- **Runtime** - Service lifecycle and process management with support for source to running auto build
- **Registry** - Centralised service discovery and API endpoint explorer with feature rich metadata
- **Store** - Key-Value storage with TTL expiry and persistent crud to keep microservices stateless
- **Web** - Simple web dashboard with dynamic forms to describe and query services in the browser

**Framework**

Micro additionally contains a built in Go framework for service development. 
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




## Installation
### Install Binaries
Windows
```sh
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```
Linux
```sh
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash
```
MacOS
```sh
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash
```

### Run the server 

```sh
micro server
```
Now go to [localhost:8080](http://localhost:8080) and make sure the output is something like `{"version": "v3.10.1"}` which is the laest version of micro installed.

## Usage

### Login to Micro
default username: `admin`

default password: `micro`

```sh
$ micro login
Enter username: admin
Enter password:
Successfully logged in.
```

See what's running:

```sh
$ micro services
api
auth
broker
config
events
network
proxy
registry
runtime
server
store
```

View in browser at localhost:8082

### Run a service

```sh
micro run github.com/micro/services/helloworld
```

### Check status of running service

```sh
$ micro status
NAME		VERSION	SOURCE					STATUS	BUILD	UPDATED	METADATA
helloworld	latest	github.com/micro/services/helloworld	running	n/a	4s ago	owner=admin, group=micro
```

### View logs of the service to verify it's running.

```sh
$ micro logs helloworld
2020-10-06 17:52:21  file=service/service.go:195 level=info Starting [service] helloworld
2020-10-06 17:52:21  file=grpc/grpc.go:902 level=info Server [grpc] Listening on [::]:33975
2020-10-06 17:52:21  file=grpc/grpc.go:732 level=info Registry [service] Registering node: helloworld-67627b23-3336-4b92-a032-09d8d13ecf95
```

### Call the service

```sh
$ micro helloworld call --name=Jane
{
	"msg": "Hello Jane"
}
```

Curl it

```
curl "http://localhost:8080/helloworld?name=John"
```

### Write a client

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service"
	proto "github.com/micro/services/helloworld/proto"
)

func main() {
	// create and initialise a new service
	srv := service.New()

	// create the proto client for helloworld
	client := proto.NewHelloworldService("helloworld", srv.Client())

	// call an endpoint on the service
	rsp, err := client.Call(context.Background(), &proto.CallRequest{
		Name: "John",
	})
	if err != nil {
		fmt.Println("Error calling helloworld: ", err)
		return
	}

	// print the response
	fmt.Println("Response: ", rsp.Message)
	
	// let's delay the process for exiting for reasons you'll see below
	time.Sleep(time.Second * 5)
}
```

Run it

```
micro run .
```

For more see the [getting started](https://micro.dev/getting-started) guide.

## 📃 Docs

- [Introduction](https://micro.dev/introduction) - A high level introduction to Micro
- [Getting Started](https://micro.dev/getting-started) - The helloworld quickstart guide
- [Upgrade Guide](https://micro.dev/upgrade-guide) - Update your go-micro project to use micro v3.
- [Architecture](https://micro.dev/architecture) - Describes the architecture, design and tradeoffs
- [Reference](https://micro.dev/reference) - In-depth reference for Micro CLI and services
- [Resources](https://micro.dev/resources) - External resources and contributions
- [Roadmap](https://micro.dev/roadmap) - Stuff on our agenda over the long haul
- [FAQ](https://micro.dev/faq) - Frequently asked questions

## 📄 License

See [LICENSE](LICENSE) which makes use of [Apache 2.0](https://opensource.org/licenses/Apache-2.0).
