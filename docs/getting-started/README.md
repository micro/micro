---
title: Getting Started with Micro
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /getting-started
summary: A getting started guide for Micro
---

## Getting Started

This is a getting started guide for Micro which teaches you how to go from install to helloworld and beyond.

## Contents
{: .no_toc }


* TOC
{:toc}

## Dependencies

You will need protoc-gen-micro for code generation

- [protobuf](https://github.com/golang/protobuf)
- [protoc-gen-go](https://github.com/golang/protobuf/tree/master/protoc-gen-go)
- [protoc-gen-micro](https://github.com/micro/micro/tree/master/cmd/protoc-gen-micro)

```
# Download latest proto releaes
# https://github.com/protocolbuffers/protobuf/releases
go get github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/micro/v3/cmd/protoc-gen-micro
```

## Install

### Go Get

Using Go:

```sh
go get github.com/micro/micro/v3
```

### Release Binary

Or by downloading the binary

```sh
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

### Docker Image

```
docker pull micro/micro
```

### Helm Chart

```
helm repo add micro https://micro.github.io/helm
helm install micro micro/micro
```

## Running a service

Before diving into writing a service, let's run an existing one, because it's just a few commands away!


First, we have to start the `micro server`. The command to do that is:

```sh
micro server
```

Before interacting with the `micro server`, we need to log in with the id 'admin' and password 'micro':

```sh
$ micro login
Enter email address: admin
Enter Password:
Successfully logged in.
```

If all goes well you'll see log output from the server listing the services as it starts them. Just to verify that everything is in order, let's see what services are running:

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

All those services are ones started by our `micro server`. This is pretty cool, but still it's not something we launched! Let's start a service for which existence we can actually take credit for. If we go to [github.com/micro/services](https://github.com/micro/services), we see a bunch of services written by micro authors. One of them is the `helloworld`. Try our luck, shall we?

The command to run services is `micro run`.

Simply issue the following command

```sh
$ micro run github.com/micro/services/helloworld
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

So since our service is running happily, let's try to call it! That's what services are for.

## Calling a service

We have a couple of options to call a service running on our `micro server`.

### With the CLI

Micro auto-generates CLI commands for your service in the form: `micro [service] [method]`, with the
default method being "Call". Arguments can be passed as flags, hence we can call our service using:

```sh
$ micro helloworld --name=Jane
{
	"msg": "Hello Jane"
}

```

That worked! If we wonder what endpoints a service has we can run the following command:

```sh
$ micro helloworld --help
NAME:
	micro helloworld

VERSION:
	latest

USAGE:
	micro helloworld [command]

COMMANDS:
	call
```

To see the flags for subcommands of `helloworld`:

```sh
$ micro helloworld call --help
NAME:
	micro helloworld call

USAGE:
	micro helloworld call [flags]

FLAGS:
	--name string
```

### With the API

Micro exposes a http API on port 8080 so you can just curl your service like so.

```
curl "http://localhost:8080/helloworld?name=John"
```

### With the Framework

Let's write a small client we can use to call the helloworld service.
Normally you'll make a service call inside another service so this is just a sample of a function you may write. We'll [learn how to write a full fledged service soon](#-writing-a-service).

Let's take the following file:

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
		rsp, err := client.Call(context.Background(), &proto.Request{
			Name: "John",
		})
		if err != nil {
			fmt.Println("Error calling helloworld: ", err)
			return
		}

		// print the response
		fmt.Println("Response: ", rsp.Msg)
		
		// let's delay the process for exiting for reasons you'll see below
		time.Sleep(time.Second * 5)
}
```

Save the example locally. For ease of following this guide, name the folder `example`.
After doing a `cd example && go mod init example`, we are ready to run this service with `micro run`:

```
micro run .
```

`micro run`s, when successful, do not print any output. A useful command to see what is running, is `micro status`. At this point we should have two services running:

```
$ micro status
NAME		VERSION	SOURCE					                STATUS	BUILD	UPDATED		METADATA
example		latest	example.tar.gz				            running	n/a 	2s ago		owner=admin, group=micro
helloworld	latest	github.com/micro/services/helloworld	running	n/a	    5m59s ago	owner=admin, group=micro
```

Now, since our example service client is also running, we should be able to see it's logs:
```sh
$ micro logs example
# some go build output here
Response:  Hello John
```

Great! That response is coming straight from the helloworld service we started earlier!

### Multi-Language Clients

Soon we'll be releasing multi language grpc generated clients to query services and use micro also.

## Creating a service

To create a new service, use the `micro new` command. It should output something reasonably similar to the following:

```sh
$ micro new helloworld
Creating service helloworld

.
├── main.go
├── generate.go
├── handler
│   └── helloworld.go
├── proto
│   └── helloworld.proto
├── Dockerfile
├── Makefile
├── README.md
├── .gitignore
└── go.mod


download protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:

visit https://github.com/protocolbuffers/protobuf/releases

download protobuf for micro:

go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/micro/v3/cmd/protoc-gen-micro

compile the proto file helloworld.proto:

cd helloworld
make proto
```

As can be seen from the output above, before building the first service, the following tools must be installed:
* [protoc](http://google.github.io/proto-lens/installing-protoc.html)
* [protobuf/proto](github.com/golang/protobuf/protoc-gen-go)
* [protoc-gen-micro](github.com/golang/protobuf/protoc-gen-go)

They are all needed to translate proto files to actual Go code.
Protos exist to provide a language agnostic way to describe service endpoints, their input and output types, and to have an efficient serialization format at hand.

So once all tools are installed, being inside the service root, we can issue the following command to generate the Go code from the protos:

```sh
make proto
```

The generated code must be committed to source control, to enable other services to import the proto when making service calls (see previous section [Calling a service](#-calling-a-service).

At this point, we know how to write a service, run it, and call other services too.
We have everything at our fingertips, but there are still some missing pieces to write applications. One of such pieces is the store interface, which helps with persistent data storage even without a database.

## Storage

Amongst many other useful built-in services Micro includes a persistent storage service for storing data.

### With the CLI

First, let's go over the more basic store CLI commands.

To save a value, we use the write command:

```sh
$ micro store write key1 value1
```

The UNIX style no output meant it was happily saved. What about reading it?

```
$ micro store read key1
val1
```

Or to display it in a fancier way, we can use the `--verbose` or `-v` flags.

```
$ micro store read -v key1
KEY    VALUE   EXPIRY
key1   val1    None
```

This view is especially useful when we use the `--prefix` or `-p` flag, which lets us search for entries which key have certain prefixes.

To demonstrate that first let's save an other value:

```
$ micro store write key2 val2
```

After this, we can list both `key1` and `key2` keys as they both share common prefixes:

```
$ micro store read --prefix --verbose key
KEY    VALUE   EXPIRY
key1   val1    None
key2   val2    None
```

There is more to the store, but this knowledge already enables us to be dangerous!

### With the Framework

Accessing the same data we have just manipulated from our Go Micro services could not be easier.
First let's create an entry that our service can read. This time we will specify the table for the `micro store write` command too, as each service has its own table in the store:


```
micro store write --table=example mykey "Hi there"
```

Let's modify [the example service we wrote previously](#-calling-a-service-with-go-micro) so instead of calling a service, it reads the above value from a store.

```go
package main

import (
	"fmt"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/store"
)

func main() {
	srv := service.New(service.Name("example"))
	srv.Init()

	records, err := store.Read("mykey")
	if err != nil {
		fmt.Println("Error reading from store: ", err)
	}

	if len(records) == 0 {
		fmt.Println("No records")
	}
	for _, record := range records {
		fmt.Printf("key: %v, value: %v\n", record.Key, string(record.Value))
	}

	time.Sleep(1 * time.Hour)
}
```

## Updating a service

Now since the example service is running (can be easily verified by `micro status`), we should not use `micro run`, but rather `micro update` to deploy it.

We can simply issue the update command (remember to switch back to the root directory of the example service first):

```sh
micro update .
```

And verify both with micro status:

```sh
$ micro status example
NAME	VERSION	SOURCE	STATUS	BUILD	UPDATED	METADATA
example	latest	n/a		running	n/a		7s ago	owner=admin, group=micro
```

that it was updated.

If things for some reason go haywire, we can try the time tested "turning it off and on again" solution and do:

```sh
micro kill example
micro run .
```

to start with a clean slate.

So once we did update the example service, we should see the following in the logs:

```sh
$ micro logs example
key: mykey, value: Hi there
```

## Config

Configuration and secrets is an essential part of any production system - let's see how the Micro config works.

### With the CLI

The most basic example of config usage is the following:

```sh
$ micro config set key val
$ micro config get key
val
```

While this alone is enough for a great many use cases, for purposes of organisation, Micro also support dot notation of keys. Let's overwrite our keys set previously:

```sh
$ micro config set key.subkey val
$ micro config get key.subkey
val
```

This is fairly straightforward, but what happens when we get `key`?

```sh
$ micro config get key
{"subkey":"val"}
```

As it can be seen, leaf level keys will return only the value, while node level keys return the whole subtree as a JSON document:

```sh
$ micro config set key.othersubkey val2
$ micro config get key
{"othersubkey":"val2","subkey":"val"}
```

### With the Framework

Micro configs work very similarly when being called from [Go code too](https://pkg.go.dev/github.com/micro/go-micro/v3/config?tab=doc):

```go
package main

import (
	"fmt"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

func main() {
	// setup the service
	srv := service.New(service.Name("example"))
	srv.Init()

	// read config value
	val, _ := config.Get("key.subkey")
	fmt.Println("Value of key.subkey: ", val.String(""))
}
```

Assuming the folder name for this service is still `example` (to update the existing service, [see updating a service](#-updating-a-service)):
```
$ micro logs example
Value of key.subkey:  val
```

## Further Resources

This is just a brief getting started guide for quickly getting up and running with Micro. 
Come back from time to time to learn more as this guide gets continually upgraded. If you're 
interested in learning more Micro magic, have a look at the following sources:

- Read the [docs](../)
- Learn by [example](https://github.com/micro/services)
- Join us on [slack](https://slack.micro.mu)
