# Micro [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/micro/micro/v3?tab=doc) [![License](https://img.shields.io/badge/license-apache-blue)](https://opensource.org/licenses/Apache-2.0) 

Micro is an API first development platform.

## Overview

Micro addresses the key requirements for building services in the cloud. It leverages the microservices
architecture pattern and provides a set of services which act as the building blocks of a platform. Micro deals
with the complexity of distributed systems and provides simpler programmable abstractions to build on. 

## Contents

- [Introduction](https://micro.dev/introduction) - A high level introduction to Micro
- [Getting Started](https://micro.dev/getting-started) - The helloworld quickstart guide
- [Upgrade Guide](https://micro.dev/upgrade-guide) - Update your go-micro project to use micro v3.
- [Architecture](https://micro.dev/architecture) - Describes the architecture, design and tradeoffs
- [Reference](https://micro.dev/reference) - In-depth reference for Micro CLI and services
- [Resources](https://micro.dev/resources) - External resources and contributions
- [Roadmap](https://micro.dev/roadmap) - Stuff on our agenda over the long haul
- [FAQ](https://micro.dev/faq) - Frequently asked questions

## Getting Started

Install micro

```sh
go install github.com/micro/micro/v3@latest
```

Run the server 

```sh
micro server
```

Login with the username 'admin' and password 'micro':

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

Run a service

```sh
micro run github.com/micro/services/helloworld
```

Now check the status of the running service

```sh
$ micro status
NAME		VERSION	SOURCE					STATUS	BUILD	UPDATED	METADATA
helloworld	latest	github.com/micro/services/helloworld	running	n/a	4s ago	owner=admin, group=micro
```

We can also have a look at logs of the service to verify it's running.

```sh
$ micro logs helloworld
2020-10-06 17:52:21  file=service/service.go:195 level=info Starting [service] helloworld
2020-10-06 17:52:21  file=grpc/grpc.go:902 level=info Server [grpc] Listening on [::]:33975
2020-10-06 17:52:21  file=grpc/grpc.go:732 level=info Registry [service] Registering node: helloworld-67627b23-3336-4b92-a032-09d8d13ecf95
```

Call the service

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

Write a client

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

## Usage

See the [docs](https://micro.dev/docs) for detailed information on the architecture, installation and use.

## License

See [LICENSE](LICENSE) which makes use of [Apache 2.0](https://opensource.org/licenses/Apache-2.0)
