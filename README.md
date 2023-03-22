# Micro

<p>
    <a href="https://goreportcard.com/report/github.com/micro/micro">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/micro/micro">
    </a>
	<a href="https://pkg.go.dev/github.com/micro/micro/v3?tab=doc"><img
    alt="Go.Dev reference"
    src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white"></a>
    <a href="https://opensource.org/licenses/Apache-2.0"><img
    alt="Apache License"
    src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"></a>
    <a href="https://marketplace.digitalocean.com/apps/micro"><img
    alt="DigitalOcean Droplet"
    src="https://img.shields.io/badge/digitalocean-droplet-blue.svg"></a>    
    <a href="https://twitter.com/MicroDotDev"><img
    alt="Twitter @MicroDotDev"
    src="https://img.shields.io/badge/twitter-follow-blue.svg"></a>    
</p>

API first development platform

## Overview 

Micro addresses the key requirements for building services in the cloud. It leverages the microservices
architecture pattern and provides a set of services which act as the building blocks of a platform. Micro deals 
with the complexity of distributed systems and provides simpler programmable abstractions to build on. 

<img src="docs/images/micro.png?v=1" />

## Features

Below are the core components that make up Micro.

**Server**

Micro is built as a microservices architecture and abstracts away the complexity of the underlying infrastructure. We compose 
this as a single logical server to the user but decompose that into the various building block primitives that can be plugged 
into any underlying system.

The server is composed of the following services.

- **API** - A Gateway which dynamically maps HTTP requests to RPC using path based resolution
- **Auth** - Authentication and authorization out of the box using JWT tokens and rule based access control.
- **Broker** - Ephemeral pubsub messaging for asynchronous communication and distributing notifications
- **Config** - Dynamic configuration and secrets management for service level config without reload
- **Events** - Event streaming with ordered messaging, replay from offsets and persistent storage
- **Network** - Inter-service networking, isolation and routing plane for all internal request traffic
- **Proxy** - An identity aware proxy used for remote access and any external grpc request traffic
- **Runtime** - Service lifecycle and process management with support for source to running auto build
- **Registry** - Centralised service discovery and API endpoint explorer with feature rich metadata
- **Store** - Key-Value storage with TTL expiry and persistent crud to keep microservices stateless

**Framework**

Micro comes with a built in Go microservices framework for service based development. 
The Go framework makes it drop dead simple to write your services without having to piece together endless lines of boilerplate code. Auto 
configured and initialised by default, just import and get started quickly.

**Command Line**

Micro brings not only a rich architectural model but a command line experience tailored for that need. The command line interface includes 
dynamic command mapping for all services running on the platform. Turns any service instantly into a CLI command along with flag parsing 
for inputs. Includes support for multiple environments and namespaces, automatic refreshing of auth credentials, creating and running 
services, status info and log streaming, plus much, much more.

**Dashboard**

Explore, discover and consume services via a browser using Micro Web. The dashboard makes use of your env configuration to locate the server 
and provides dynamic form fill for services.

**Environments**

Micro bakes in the concept of `Environments` and multi-tenancy through `Namespaces`. Run your server locally for 
development and in the cloud for staging and production, seamlessly switch between them using the CLI commands `micro env set [environment]` 
and `micro user set [namespace]`.

## Docs

- [Introduction](https://micro.dev/introduction) - A high level introduction to Micro
- [Getting Started](https://micro.dev/getting-started) - The helloworld quickstart guide
- [Upgrade Guide](https://micro.dev/upgrade-guide) - Update your go-micro project to use micro v3.
- [Architecture](https://micro.dev/architecture) - Describes the architecture, design and tradeoffs
- [Reference](https://micro.dev/reference) - In-depth reference for Micro CLI and services
- [Resources](https://micro.dev/resources) - External resources and contributions
- [Roadmap](https://micro.dev/roadmap) - Stuff on our agenda over the long haul
- [FAQ](https://micro.dev/faq) - Frequently asked questions

## Installation

### From Source

```
go install github.com/micro/micro/v3@latest
```

### Docker Image

```
docker pull ghcr.io/micro/micro:latest
```

### Install Binaries

#### Windows

```sh
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"

```
#### Linux

```sh
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash
```

#### MacOS

```sh
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash
```

### Run the server 

The server starts with a single command ready to use

### Local
```sh
micro server
```

### Docker

```
sudo docker run -p 8080:8080 -p 8081:8081 ghcr.io/micro/micro:latest server
```

Now go to [localhost:8080](http://localhost:8080) and make sure the output is something like `{"version": "v3.10.1"}` 
which is the latest version of micro installed.

## Usage

Set the environment e.g local

```
micro env set local
```

### Login to Micro

Default username/password: `admin/micro`

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

### Create a Service

Generate a service using the template

```
micro new helloworld
```

Output

```
Creating service helloworld

.
├── main.go
├── handler
│   └── helloworld.go
├── proto
│   └── helloworld.proto
├── Makefile
├── README.md
├── .gitignore
└── go.mod


download protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:

visit https://github.com/protocolbuffers/protobuf/releases

compile the proto file helloworld.proto:

cd helloworld
make init
go mod vendor
make proto
```

### Making changes

Edit the protobuf definition in `proto/helloworld.proto` and run `make proto` to recompile

Go to `handler/helloworld.go` to make changes to the response handler

```go
type Helloworld struct{}

func New() *Helloworld {
        return &Helloworld{}
}

func (h *Helloworld) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
        rsp.Msg = "Hello " + req.Name
        return nil
}
```

### Run the service

Run from local dir

```
micro run .
```

Or from a git url

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

### Write a service client

A service client is used within another service and must be run by micro

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service"
	pb "github.com/micro/services/helloworld/proto"
)

func callService(hw pb.HelloworldService) {
	for {
		// call an endpoint on the service
		rsp, err := hw.Call(context.Background(), &pb.CallRequest{
			Name: "John",
		})
		if err != nil {
			fmt.Println("Error calling helloworld: ", err)
			return
		}

		// print the response
		fmt.Println("Response: ", rsp.Message)

		time.Sleep(time.Second)
	}
}

func main() {
	// create and initialise a new service
	srv := service.New(
		service.Name("caller"),
	)

	// new helloworld client
	hw := pb.NewHelloworldService("helloworld", srv.Client())
	
	// run the client caller
	go callService(hw)
	
	// run the service
	service.Run()
}
```

Run it

```
micro run .
```

### Write an api client

An api client is an external app or client which makes requests through the micro api

Get a token
```
export TOKEN=`micro user token`
```

Call helloworld
```go
package main

import (
    "fmt"
    "os"

    "github.com/micro/micro-go"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Msg string `json:"msg"`
}

func main() {
	token := os.Getenv("TOKEN")
	c := micro.NewClient(nil)

	// set your api token
	c.SetToken(token)

   	req := &Request{
		Name: "John",
	}
	
	var rsp Response

	if err := c.Call("helloworld", "Call", req, &rsp); err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println(rsp)
}
```

Run it

```
go run main.go
```

For more see the [getting started](https://micro.dev/getting-started) guide.

## Web Dashboard

Use services via the web with the Micro Web dashboard

```
micro web
```

Browse to `localhost:8082`

## License

See [LICENSE](LICENSE) which makes use of [Apache 2.0](https://opensource.org/licenses/Apache-2.0).

## Development

[1 click deploy](https://marketplace.digitalocean.com/apps/micro) a Micro Dev environment on a DigitalOcean Droplet

Use our [refcode](https://marketplace.digitalocean.com/apps/micro?refcode=1eb1b2aca272&action=deploy) so we get $25 credit too!
