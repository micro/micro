---
layout:	post
author: Janos Dobronszki
title:	Micro Server - Getting started with microservices
date:	2020-05-05 10:00:00
---
<br />
In this post we will have a look at how to run and manage microservices locally with `micro server` and the Micro CLI in general.
The Micro CLI consists of both the server command and other client commands that enable us to interact with the server.
`micro server` can run microservices in different environments - binaries locally for speed and simplicity, or containers in a more production ready environment.

## Installation

Using Go:

```sh
go install github.com/micro/micro/v2
```

Or by downloading the binary

```sh
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

## Running a service

Before diving into writing a service, let's run an existing one, because it's just a few commands away!


First, we have to start the `micro server`. The command to do that is:

```sh
micro server
```

To talk to this server, we just have to tell Micro CLI to address our server instead of using the default implementations - micro can work without a server too, but [more about that later](#-environments).

The following command tells the CLI to talk to our server:

```
micro env set server
```

Great! We are ready to roll. Just to verify that everything is in order, let's see what services are running:

```
$ micro list services
go.micro.api
go.micro.auth
go.micro.bot
go.micro.broker
go.micro.config
go.micro.debug
go.micro.network
go.micro.proxy
go.micro.registry
go.micro.router
go.micro.runtime
go.micro.server
go.micro.web
```

All those services are ones started by our `micro server`. This is pretty cool, but still it's not something we launched! Let's start a service for which existence we can actually take credit for. If we go to [github.com/micro](https://github.com/micro), we see a bunch of services written by micro authors. One of them is the `helloworld`. Try our luck, shall we?

The command to run services is `micro run`. This command may take a while as it checks out
the repository from GitHub. (@todo this actually fails currently, fix)

```
micro run github.com/micro/helloworld
```


If we take a look at the running `micro server`, we should see something like

```
Creating service helloworld version latest source /tmp/github.com-micro-services/helloworld
Processing create event helloworld:latest
```

We can also have a look at logs of the service to verify it's running.

```sh
$ micro logs helloworld
Starting [service] go.micro.service.helloworld
Server [grpc] Listening on [::]:36577
Registry [service] Registering node: go.micro.service.helloworld-213b807a-15c2-496f-93ac-7949ad38aadf
```

So since our service is running happily, let's try to call it! That's what services are for.

## Calling a service

We have a couple of options to call a service running on our `micro server`.

### Calling a service from CLI

The easiest is perhaps with the CLI:

```sh
$ micro call go.micro.service.helloworld Helloworld.Call '{"name":"Jane"}'
{
	"msg": "Hello Jane"
}

```

That worked! If we wonder what endpoints a service has we can run the following command:

```sh
micro get service go.micro.service.helloworld
```

Otherwise the best place to look is at the [proto definition](https://github.com/micro/blob/master/helloworld/proto/helloworld/helloworld.proto). You can also browse to the UI at [http://localhost:8082](http://localhost:8082/service/go.micro.service.helloworld) to see live info.

### Calling a service with Go Micro

Let's write a small client we can use to call the helloworld service.
Normally you'll make a service call inside another service so this is just a sample of a function you may write. We'll [learn how to write a full fledged service soon](#-writing-a-service).

Let's take the following file:

```go
package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2"
	proto "github.com/micro/helloworld/proto"
)

func main() {
	// create and initialise a new service
	service := micro.NewService()
	service.Init()

	// create the proto client for helloworld
	client := proto.NewHelloworldService("go.micro.service.helloworld", service.Client())

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

Save the example locally. For ease of following this guide, name the folder `example-service`.
After doing a `cd example-service && go mod init example`, we are ready to run this service with `micro run`:

```
micro run .
```

An other useful command to see what is running, is `micro status`. At this point we should have two services running:

```
$ micro status
NAME			VERSION	SOURCE										STATUS		BUILD	UPDATED		METADATA
example-service	latest	/home/username/example-service				starting	n/a		4s ago		owner=n/a,group=n/a
helloworld		latest	/tmp/github.com-micro-services/helloworld	running		n/a		6m5s ago	owner=n/a,group=n/a
```

Now, since our example-service client is also running, we should be able to see it's logs:
```sh
$ micro logs example-service
# some go build output here
Response:  Hello John
```

Great! That response is coming straight from the helloworld service we started earlier!

### From other languages

In the [clients repo](https://github.com/micro/clients) there are Micro clients for various languages and frameworks. They are designed to connect easily to the live Micro environment or your local one, but more about environments later.

## Writing a service

To scaffold a new service, the `micro new` command can be used. It should output something
reasonably similar to the following:

```sh
$ micro new foobar
Creating service go.micro.service.foobar in foobar

.
├── main.go
├── generate.go
├── plugin.go
├── handler
│   └── foobar.go
├── subscriber
│   └── foobar.go
├── proto/foobar
│   └── foobar.proto
├── Dockerfile
├── Makefile
├── README.md
├── .gitignore
└── go.mod


download protobuf for micro:

brew install protobuf
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/micro/v2/cmd/protoc-gen-micro@master

compile the proto file foobar.proto:

cd foobar
protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/foobar/foobar.proto
```

As can be seen from the output above, before building the first service, the following tools must be installed:
* [protoc](http://google.github.io/proto-lens/installing-protoc.html)
* [protobuf/proto](github.com/golang/protobuf/protoc-gen-go)
* [protoc-gen-micro](github.com/golang/protobuf/protoc-gen-go)

They are all needed to translate proto files to actual Go code.
Protos exist to provide a language agnostic way to describe service endpoints, their input and output types, and to have an efficient serialization format at hand.

Currently Micro is  Go focused (apart from the [clients](#-from-other-languages) mentioned before), but this will change soon.

So once all tools are installed, being inside the service root, we can issue the following command to generate the Go code from the protos:

```
protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/foobar/foobar.proto
```

The generated code must be committed to source control, to enable other services to import the proto when making service calls (see previous section [Calling a service](#-calling-a-service).

At this point, we know how to write a service, run it, and call other services too.
We have everything at our fingertips, but there are still some missing pieces to write applications. One of such pieces is the store interface, which helps with persistent data storage even without a database.

## Storage

Amongst many other useful built-in services Micro includes a persistent storage service for storing data.

### Interfaces as building blocks

A quick side note. Micro (the server/CLI) and Go Micro (the framework) are centered around strongly defined interfaces which are pluggable and provide an abstraction for underlying distributed systems concepts. What does this mean?

Let's take our current case of the [store interface](https://github.com/micro/go-micro/blob/master/store/store.go). It's aimed to enable service writers data storage with a couple of different implementations:

* in memory
* file storage (default when running `micro server`)
* cockroachdb

Similarly, the [runtime](https://github.com/micro/go-micro/blob/master/runtime/runtime.go) interface, that allows you to run services in a completely runtime agnostic way has a few implementations:

* local, which just runs actual processes - aimed at local development
* kubernetes - for running containers in a highly available and distributed way

This is a recurring theme across Micro interfaces. Let's take a look at the default store when running `micro server`.

### Using the Store

#### Using the store with CLI

First, let's go over the more basic store CLI commands.

To save a value, we use the write command:

```sh
micro store write key1 value1
```

The UNIX style no output meant it was happily saved. What about reading it?

```
$ micro store read key1
val1
```

Or to display it in a fancier way, we can use the `--verbose` or `-v` flags:

```
KEY    VALUE   EXPIRY
key1   val1    None
```

This view is especially useful when we use the `--prefix` or `-p` flag, which lets us search for entries which key have certain prefixes.

To demonstrate that first let's save an other value:

```
micro store write key2 val2
```

After this, we can list both `key1` and `key2` keys as they both share commond prefixes:

```
$ micro store read --prefix --verbose key
KEY    VALUE   EXPIRY
key1   val1    None
key2   val2    None
```

There are more to the store, but this knowledge already enables us to be dangerous!

#### Using the Store with Go-Micro

Accessing the same data we have just manipulated from our Go Micro services could not be easier.
First let's create an entry that our service can read. This time we will specify the table for the `micro store write` command too, as each service has its own table in the store:


```
micro store write --table go.micro.service.example mykey "Hi there"
```

Let's modify [the example service we wrote previously](#-calling-a-service-with-go-micro) so instead of calling a service, it reads the above value from a store.

```go
package main

import (
	"fmt"
	"time"

	"github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService()

	service.Init(micro.Name("go.micro.service.example"))

	records, err := service.Options().Store.Read("mykey")
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

We are almost done! But first we have to learn how to update a service.

## Updating and killing a service

Now since the example service is running (can be easily verified by `micro status`), we should not use `micro run`, but rather `micro update` to deploy it.

We can simply issue

```
micro update .
```

And verify both with the micro server output:

```
Updating service example-service version latest source /home/username/example-service
Processing update event example-service:latest in namespace default
```

and micro status:

```
$ micro status example-service
NAME			VERSION	SOURCE							STATUS		BUILD	UPDATED		METADATA
example-service	latest	/home/username/example-service	starting	n/a		10s ago		owner=n/a,group=n/a
```

that it was updated.

If things for some reason go haywire, we can try the time tested "turning it off and on again" solution and do:

```
micro kill .
micro run .
```

to start with a clean slate.

So once we did update the example service, we should see the following in the logs:

```
$ micro logs example-service
key: mykey, value: Hi there
```

Nice! The example service read the value from the store successfully.

## Clients

Beyond this we're working on multi-language clients which you can find and contribute to 
on github at [github.com/micro/clients](https://github.com/micro/clients). We'd love to 
discuss this further but it's not quite ready.

## Further reading

This is just a brief getting started guide for quickly getting up and running with Micro. 
Come back from time to time to learn more as this guide gets continually upgraded. If you're 
interested in learning more Micro magic, have a look at the following sources:

- Read the [docs](https://m3o.com.docs)
- Learn by [examples](https://github.com/micro/examples)
- Come join us on [Slack](https://slack.m3o.com) and ask quesions

Cheers

From the team at Micro
