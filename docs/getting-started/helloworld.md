---
title: Hello World
keywords: micro, framework, helloworld
tags: [micro, framework, helloworld]
sidebar: home_sidebar
permalink: /helloworld
summary: A hello world example using Micro
---

# Overview

This is a helloworld example of using Micro. We'll run through writing the whole thing.

## Installation

Using Go:

```sh
go get github.com/micro/micro/v3
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

Before interacting with the `micro server`, we need to log in with the id 'admin' and password 'micro':

```sh
$ micro login
Enter email address: admin
Enter Password:
Successfully logged in.
```

If all goes well you'll see log output from the server listing the services as it starts them. Just to verify that everything is in order, let's see what services are running:

```
$ micro services
api
auth
broker
config
network
proxy
registry
runtime
server
store
```

All those services are ones started by our `micro server`. This is pretty cool, but still it's not something we launched! Let's start a service for which existence we can actually take credit for. If we go to [github.com/micro/services](https://github.com/micro/services), we see a bunch of services written by micro authors. One of them is the `helloworld`. Try our luck, shall we?

The command to run services is `micro run`.

```
micro run github.com/micro/services/helloworld
```


If we take a look at the running services using `micro status`, we should see the service listed:
```
NAME		VERSION	SOURCE		STATUS	BUILD	UPDATED	METADATA
helloworld	latest	helloworld	running	n/a	unknown	owner=n/a,group=n/a
```

We can also have a look at logs of the service to verify it's running.

```sh
$ micro logs helloworld
micro@Bens-MBP-3 micro % micro logs helloworld
2020-08-11 15:18:33  file=service/service.go:192 level=info Starting [service] helloworld
2020-08-11 15:18:33  file=grpc/grpc.go:902 level=info Server [grpc] Listening on [::]:49602
2020-08-11 15:18:33  file=grpc/grpc.go:728 level=info Registry [service] Registering node: helloworld-c49ee2a3-e9d0-4411-9b9b-5fe6aea6b49d
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

That worked! If we wonder what nodes and endpoints a service has we can run the following command:

```sh
micro get service helloworld
```

### With the framework

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

Save the example locally. For ease of following this guide, name the folder `example-service`.
After doing a `cd example-service && go mod init example`, we are ready to run this service with `micro run`:

```
micro run .
```

`micro run`s, when successful, do not print any output. A useful command to see what is running, is `micro status`. At this point we should have two services running:

```
$ micro status
NAME							VERSION		SOURCE																	STATUS		BUILD	UPDATED		METADATA
example-service		latest		example-service													starting	n/a		4s ago		owner=n/a,group=n/a
helloworld				latest		github.com/micro/services/helloworld		running		n/a		unknown		owner=n/a,group=n/a
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

## Creating a service

To create a new service, use the `micro new` command. It should output something reasonably similar to the following:

```sh
$ micro new helloworld
Creating service helloworld in helloworld

.
├── main.go
├── generate.go
├── plugin.go
├── handler
│   └── helloworld.go
├── subscriber
│   └── helloworld.go
├── proto/helloworld
│   └── helloworld.proto
├── Dockerfile
├── Makefile
├── README.md
├── .gitignore
└── go.mod


download protobuf for micro:

brew install protobuf
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/micro/v3/cmd/protoc-gen-micro@master

compile the proto file helloworld.proto:

cd helloworld
protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/helloworld/helloworld.proto
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
protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/helloworld.proto
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

#### With the CLI

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

After this, we can list both `key1` and `key2` keys as they both share commond prefixes:

```
$ micro store read --prefix --verbose key
KEY    VALUE   EXPIRY
key1   val1    None
key2   val2    None
```

There is more to the store, but this knowledge already enables us to be dangerous!

#### With the framework

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
example-service	latest	example-service.tar.gz			running	n/a	unknown	owner=n/a,group=n/a
```

that it was updated.

If things for some reason go haywire, we can try the time tested "turning it off and on again" solution and do:

```
micro kill example-service
micro run example-service
```

to start with a clean slate.

So once we did update the example service, we should see the following in the logs:

```
$ micro logs example-service
key: mykey, value: Hi there
```

## Config

Configuration and secrets is an essential part of any production system - let's see how the Micro config works.

### CLI

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

### With the framework

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
	fmt.Println("Value of key.subkey: ", config.Get("key", "subkey").String(""))
}
```

Assuming the folder name for this service is still `example-service` (to update the existing service, [see updating a service](#-updating-a-service)):
```
$ micro logs example-service
Value of key.subkey:  val
```

