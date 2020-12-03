---
title: Reference
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference
summary: Reference - a comprehensive guide to Micro
---

## Reference
{: .no_toc }

This reference doc is an in depth guide for the technical details and usage of Micro

## Contents
{: .no_toc }

* TOC
{:toc}

## Overview

Micro is a platform for cloud native development. It consists of a server, command line interface and 
service framework which enables you to build, run, manage and consume Micro services. This reference 
walks through the majority of Micro in depth and attempts to help guide you through any usage. It 
should be thought of much like a language spec and will evolve over time.

## Installation

Below are the instructions for installing micro locally, in docker or on kubernetes

### Local 

Micro can be installed locally in the following way. We assume for the most part a Linux env with Go and Git installed.

#### Go Get

```
go get github.com/micro/micro/v3
```

#### Docker

```sh
docker pull micro/micro
```

#### Release Binaries

```sh
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

### Kubernetes

Micro can be installed onto a Kubernetes cluster using helm. Micro will be deployed in full and leverage zero-dep implementations designed for Kubernetes. For example, micro store will internally leverage a file store on a persistent volume, meaning there are no infrastructure dependencies required.

#### Dependencies

You will need to be connected to a Kubernetes cluster

#### Install

Install micro with the following commands:

```shell
helm repo add micro https://micro.github.io/helm
helm install micro micro/micro
```

#### Uninstall

Uninstall micro with the following commands:

```shell
helm uninstall micro
helm repo remove micro
```

## Server

The micro server is a distributed systems runtime for the Cloud and beyond. It provides the building 
blocks for distributed systems development as a set of services, command line and service framework. 
The server is much like a distributed operating system in the sense that each component runs 
independent of each other but work together as one system. This composition allows us to use a 
microservices architecture pattern even for the platform.

### Features

The server provides the below functionality as built in primitives for services development.

- Authentication
- Configuration
- PubSub Messaging
- Event Streaming
- Service Discovery
- Service Networking
- Key-Value Storage
- HTTP API Gateway
- gRPC Identity Proxy

### Usage

To start the server simply run

```sh
micro server
```

This will boot the entire system and services including a http api on :8080 and grpc proxy on :8081

### Help

Run the following command to check help output
```
micro --help
```

### Commands

Run helloworld and check its status

```
micro env	# should point to local
micro run github.com/micro/services/helloworld # run helloworld
micro status 	# wait for status running
micro services	# should display helloworld
```

Call the service and verify output

```sh
$ micro helloworld --name=John
{
        "msg": "Hello John"
}
```

Remove the service

```
micro kill helloworld
```

## Command Line

The command line interface is the primary way to interact with a micro server. It's a simple binary that 
can either be interacted with using simple commands or an interactive prompt. The CLI proxies all commands 
as RPC calls to the Micro server. In many of the builtin commands it will perform formatting and additional 
syntactic work.

### Builtin Commands

Built in commands are system or configuration level commands for interacting with the server or 
changing user config. For the most part this is syntactic sugar for user convenience. Here's a 
subset of well known commands.

```
signup
login
run
update
kill
services
logs
status
env
user
```

The micro binary and each subcommand has a --help flag to provide a usage guide. The majority should be 
obvious to the user. We will go through a few in more detail.

#### Signup

Signup is a command which attempts to query a "signup" to register a new account, this is env specific and requires a signup service to be 
running. By default locally this will not exist and we expect the user to use the admin/micro credentials to administer the system. 
You can then choose to run your own signup service conforming to the proto in micro/proto or use `micro auth create account`. 

Signup is seen as a command for those who want to run their own micro server for others and potentially license the software to take payment.

#### Login

Login authenticates the user and stores credentials locally in a .micro/tokens file. This calls the micro auth service to authenticate the 
user against existing accounts stored in the system. Login asks for a username and password at the prompt.

### Dynamic Commands

When issuing a command to the Micro CLI (ie. `micro command`), if the command is not a builtin, Micro will try to dynamically resolve this command and call
a service running. Let's take the `micro registry` command, because although the registry is a core service that's running by default on a local Micro setup,
the `registry` command is not a builtin one.

With the `--help` flag, we can get information about available subcommands and flags

```sh
$ micro registry --help
NAME:
	micro registry

VERSION:
	latest

USAGE:
	micro registry [command]

COMMANDS:
	deregister
	getService
	listServices
	register
	watch
```

The commands listed are endpoints of the `registry` service (see `micro services`).

To see the flags (which are essentially endpoint request parameters) for a subcommand:

```sh
$ micro registry getService --help
NAME:
	micro registry getService

USAGE:
	micro registry getService [flags]

FLAGS:
	--service string
	--options_ttl int64
	--options_domain string

```

At this point it is useful to have a look at the proto of the [registry service here](https://github.com/micro/micro/blob/master/proto/registry/registry.proto).

In particular, let's see the `GetService` endpoint definition to understand how request parameters map to flags:

```proto
message Options {
	int64 ttl = 1;
	string domain = 2;
}

message GetRequest {
	string service = 1;
	Options options = 2;
}
```

As the above definition tells us, the request of `GetService` has the field `service` at the top level, and fields `ttl` and `domain` in an options structure.
The dynamic CLI maps the underscored flagnames (ie. `options_domain`) to request fields, so the following request JSON:

```js
{
    "service": "serviceName",
    "options": {
        "domain": "domainExample"
    }
}
```

is equivalent to the following flags:

```sh
micro registry getService --service=serviceName --options_domain=domainExample
```

### User Config

The command line uses local user config stores in ~/.micro for any form of state such as saved environments, 
tokens, etc. It will always attempt to read from here unless specified otherwise. Currently we store all 
config in a single file `config.json` and any auth tokens in a `tokens` file.

## Environments

Micro is built with a federated and multi-environment model in mind. Our development normally maps through local, staging and production, so Micro takes 
this forward looking view and builds in the notion of environments which are completely isolated micro environments you can interact with through the CLI. 
This reference explains environments.

### View Current

Environments can be displayed using the `micro env` command.

```sh
$ micro env
* local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
```

There are three builtin environments, `local` being the default, and two [`m3o` specific](m3o.com) offerings; dev and platform.
These exist for convenience and speed of development. Additional environments can be created using `micro env add [name] [host:port]`. 
Environment addresses point to the micro proxy which defaults to :8081.

### Add Environment

The command `micro env --help` provides a summary of usage. Here's an example of how to add an environment.

```sh
$ micro env add foobar example.com
$ micro env
* local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
  foobar     example.com
```

### Set Environment

The `*` marks which environment is selected. Let's select the newly added:

```sh
$ micro env set foobar
$ micro env
  local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
* foobar     example.com
```

### Login to an Environment

Each environment is effectively an isolated deployment with its own authentication, storage, etc. So each env requires signup and login. At this point we have to log in to the `example` env with `micro login`. If you don't have credentials to the environment, you have to ask the admin.

## Services

Micro is built as a distributed operating system leveraging the microservices architecture pattern.

### Overview

Below we describe the list of services provided by the Micro Server. Each service is considered a 
building block primitive for a platform and distributed systems development. The proto 
interfaces for each can be found in [micro/proto/auth](https://github.com/micro/micro/blob/master/proto/auth/auth.proto) 
and the Go library, client and server implementations in [micro/service/auth](https://github.com/micro/micro/tree/master/service/auth).

### API

The API service is a http API gateway which acts as a public entrypoint and converts http/json to RPC.

#### Overview

The micro API is the public entrypoint for all external access to services to be consumed by frontend, mobile, etc. The api 
accepts http/json requests and uses path based routing to resolve to backend services. It converts the request to gRPC and 
forward appropriately. The idea here is to focus on microservices on the backend and stitch everything together as a single 
API for the frontend. 

#### Usage

In the default `local` [environment](#environments) the API address is `127.0.0.1:8080`.
Each service running is callable through this API.

```sh
$ curl http://127.0.0.1:8080/
{"version": "v3.0.0-beta"}
```

An example call would be listing services in the registry:

```sh
$ curl http://127.0.0.1:8080/registry/listServices
```

The format is 
```
curl http://127.0.0.1:8080/[servicename]/[endpointName]
```

The endpoint name is lower camelcase.

The parameters can be passed on as query params

```sh
$ curl http://127.0.0.1:8080/helloworld/call?name=Joe
{"msg":"Hello Joe"}
```

or JSON body:

```sh
curl -XPOST --header "Content-Type: application/json" -d '{"name":"Joe"}' http://127.0.0.1:8080/helloworld/call
{"msg":"Hello Joe"}
```

To specify a namespace when calling the API, the `Micro-Namespace` header can be used:

```sh
$ curl -H "Micro-Namespace: foobar" http://127.0.0.1:8080/helloworld/call?name=Joe
```

To call a [non-public service/endpoint](#auth), the `Authorization` header can be used:

```sh
MICRO_API_TOKEN=`micro user token`
curl -H "Authorization: Bearer $MICRO_API_TOKEN" http://127.0.0.1:8080/helloworld/call?name=Joe
```

### Auth

The auth service provides both authentication and authorization.

#### Overview

The auth service stores accounts and access rules. It provides the single source of truth for all authentication 
and authorization within the Micro runtime. Every service and user requires an account to operate. When a service 
is started by the runtime an account is generated for it. Core services and services run by Micro load rules 
periodically and manage the access to their resources on a per request basis.

#### Usage

For CLI command help use `micro auth --help` or auth subcommand help such as `micro auth create --help`.

#### Login

To login to a server simply so the following

```
$ micro login
Enter email address: admin
Enter Password: 
Successfully logged in.
```

Assuming you are pointing to the right environment. It defaults to the local `micro server`.

#### Rules

Rules determine what resource a user can access. The default rule is the following:

```sh
$ micro auth list rules
ID          Scope           Access      Resource        Priority
default     <public>        GRANTED     *:*:*           0
```

The `default` rule makes all services callable that appear in the `micro status` output.
Let's see an example of this.

```sh
$ micro run helloworld
# Wait for the service to accept calls
$ curl 127.0.0.1:8080/helloworld/call?name=Alice
{"msg":"Hello Alice"}
```

If we want to prevent unauthorized users from calling our services, we can create the following rule

```sh
# This command creates a rule that enables only logged in users to call the micro server
micro auth create rule  --access=granted --scope='*' --resource="*:*:*" onlyloggedin
```

and delete the default one.
Here, the scope `*` is markedly different from the `<public>` scope we have seen earlier when doing a `micro auth list rules`:

```sh
$ micro auth list rules
ID            Scope         Access       Resource       Priority
onlyloggedin  *             GRANTED      *:*:*          0
default       <public>      GRANTED      *:*:*          0
```

Now, let's remove the default rule.

```sh
# This command deletes the 'default' rule - the rule which enables anyone to call the 'micro server'.
$ micro auth delete rule default
Rule deleted
```

Let's try curling our service again:

```sh
$ curl 127.0.0.1:8080/helloworld/call?name=Alice
{"Id":"helloworld","Code":401,"Detail":"Unauthorized call made to helloworld:Helloworld.Call","Status":"Unauthorized"}
```

Our `onlyloggedin` rule took effect. We can still call the service with a token:

```sh
$ token=$(micro user token)
# Locally:
# curl "Authorization: Bearer $token" 127.0.0.1:8080/helloworld/call?name=Alice
{"msg":"Hello Alice"}
```

(Please note tokens have a limited lifetime so the line `$ token=$(micro user token)` has to be reissued from time to time, or the command must be used inline.)

#### Accounts

Auth service supports the concept of accounts. The default account used to access the `micro server` is the admin account.

```sh
$ micro auth list accounts
ID		Name		Scopes		Metadata
admin		admin		admin		n/a
```

We can create accounts for teammates and coworkers with `micro auth create account`:

```sh
$ micro auth create account --scopes=admin jane
Account created: {"id":"jane","type":"","issuer":"micro","metadata":null,"scopes":["admin"],"secret":"bb7c1a96-c0c6-4ff5-a0e9-13d456f3db0a","name":"jane"}
```

The freshly created account can be used with `micro login` by using the `jane` id and `bb7c1a96-c0c6-4ff5-a0e9-13d456f3db0a` password.

### Broker

The broker is a message broker for asynchronous pubsub messaging.

#### Overview

The broker provides a simple abstraction for pubsub messaging. It focuses on simple semantics for fire-and-forget 
asynchronous communication. The goal here is to provide a pattern for async notifications where some update or 
events occurred but that does not require persistence. The client and server build in the ability to publish 
on one side and subscribe on the other. The broker provides no message ordering guarantees.

While a Service is normally called by name, messaging focuses on Topics that can have multiple publishers and 
subscribers. The broker is abstracting away in the service's client/server which includes message encoding/decoding 
so you don't have to spend all your time marshalling.

##### Client

The client contains the `Publish` method which takes a proto message, encodes it and publishes onto the broker 
on a given topic. It takes the metadata from the client context and includes these as headers in the message 
including the content-type so the subscribe side knows how to deal with it.

##### Server

The server supports a `Subscribe` method which allows you to register a handler as you would for handling requests. 
In this way we can mirror the handler behaviour and deserialize the message when consuming from the broker. In 
this model the server handles connecting to the broker, subscribing, consuming and executing your subscriber
function.

#### Usage

Publisher:
```go
bytes, err := json.Marshal(&Healthcheck{
	Healthy: true,
	Service: "foo",
})
if err != nil {
	return err
}

return broker.Publish("health", &broker.Message{Body: bytes})
```

Subscriber:
```go
handler := func(msg *broker.Message) error {
	var hc Healthcheck
	if err := json.Unmarshal(msg.Body, &hc); err != nil {
		return err
	}
	
	if hc.Healthy {
		logger.Infof("Service %v is healthy", hc.Service)
	} else {
		logger.Infof("Service %v is not healthy", hc.Service)
	}

	return nil
}

sub, err := broker.Subscribe("health", handler)
if err != nil {
	return err
}
```

### Config

The config service provides dynamic configuration for services. 

#### Overview

Config can be stored and loaded separately to 
the application itself for configuring business logic, api keys, etc. We read and write these as key-value 
pairs which also support nesting of JSON values. The config interface also supports storing secrets by 
defining the secret key as an option at the time of writing the value.

#### Usage

Let's assume we have a service called `helloworld` from which we want to read configuration data.
First we have to insert said data with the cli. Config data can be organized under different "paths" with the dot notation.
It's a good convention to save all config data belonging to a service under a top level path segment matching the service name:

```sh
$ micro config set helloworld.somekey hello
$ micro config get helloworld.somekey
hello
```

We can save another key too and read all values in one go with the dot notation:

```sh
$ micro config set helloworld.someotherkey "Hi there!"
$ micro config get helloworld
{"somekey":"hello","someotherkey":"Hi there!"}
```

As it can be seen, the config (by default) stores configuration data as JSONs.
We can save any type:

```sh
$ micro config set helloworld.someboolkey true
$ micro config get helloworld.someboolkey
true
$ micro config get helloworld
{"someboolkey":true,"somekey":"hello","someotherkey":"Hi there!"}
```

So far we have only saved top level keys. Let's explore the advantages of the dot notation.

```sh
$ micro config set helloworld.keywithsubs.subkey1 "So easy!"
{"keywithsubs":{"subkey1":"So easy!"},"someboolkey":true,"somekey":"hello","someotherkey":"Hi there!"}
```

Some of the example keys are getting in our way, let's learn how to delete:

```sh
$ micro config del helloworld.someotherkey
$ micro config get helloworld
{"keywithsubs":{"subkey1":"So easy!"},"someboolkey":true,"somekey":"hello"}
```

We can of course delete not just `leaf` level keys, but top level ones too:

```sh
$ micro config del helloworld.keywithsubs
$ micro config get helloworld
{"someboolkey":true,"somekey":"hello"}
```

##### Secrets

The config also supports secrets - values encrypted at rest. This helps in case of leaks, be it a security one or an accidental copy paste.

They are fairly easy to save:

```sh
$ micro config set --secret helloworld.hushkey "Very secret stuff" 
$ micro config get helloworld.hushkey
[secret]

$ micro config get --secret helloworld.hushkey
Very secret stuff

$ micro config get helloworld
{"hushkey":"[secret]","someboolkey":true,"somekey":"hello"}

$ micro config get --secret helloworld
{"hushkey":"Very secret stuff","someboolkey":true,"somekey":"hello"}
```

Even bool or number values can be saved as secrets, and they will appear as the string constant `[secret]` unless decrypted:

```sh
$ micro config set --secret helloworld.hush_number_key 42
$ micro config get helloworld
{"hush_number_key":"[secret]","hushkey":"[secret]","someboolkey":true,"somekey":"hello"}

# micro config get --secret helloworld
{"hush_number_key":42,"hushkey":"Very secret stuff","someboolkey":true,"somekey":"hello"}
```

#### Service Framework

It is similarly easy to access and set config values from a service.
A good example of reading values is [the config example test service](https://github.com/micro/micro/tree/master/test/service/config-example):

```go
package main

import (
	"fmt"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

type keyConfig struct {
	Subkey  string `json:"subkey"`
	Subkey1 int    `json:"subkey1"`
}

type conf struct {
	Key keyConfig `json:"key"`
}

func main() {
	go func() {
		for {
			time.Sleep(time.Second)
			val, err := config.Get("key.subkey")
			fmt.Println("Value of key.subkey: ", val.String(""), err)

			val, err = config.Get("key", config.Secret(true))
			if err != nil {
				fmt.Println(err)
			}
			c := conf{}
			err = val.Scan(&c.Key)
			fmt.Println("Value of key.subkey1: ", c.Key.Subkey1, err)
		}
	}()

	// run the service
	service.Run()
}
```

The above service will print the value of `key.subkey` and `key.subkey` every second.
By passing in the `config.Secret(true)` option, we tell config to decrypt secret values for us, similarly to the `--secret` CLI flag.

The [config interface](https://github.com/micro/go-micro/blob/master/config/config.go) specifies not just `Get` `Set` and `Delete` to access values,
but a few convenience functions too in the `Value` interface.

It is worth noting that `String` `Int` etc methods will do a best effort try at coercing types, ie. if the value saved is a string, `Int` will try to parse it.
However, the same does not apply to the `Scan` method, which uses `json.Unmarshal` under the hood, which we all know fails when encountering type mismatches.

`Get` should, in all cases, return a non nil `Value`, so even if the `Get` errors, `Value.Int()` and other operations should never panic.

#### Advanced Concepts

##### Merging Config Values

When saving a string with the CLI that is a valid JSON map, it gets expanded to be saved as a proper map structure, instead of a string, ie

```sh
$ micro config set helloworld '{"a": "val1", "b": "val2"}'
$ micro config get helloworld.a
val1
# If the string would be saved as is, `helloworld.a` would be a nonexistent path
```

The advantages of this become particularly visible when `Set`ting a complex type with the library:

```go
type conf struct {
	A string `json:"a"`
	B string `json:"b"`
}

c1 := conf{"val1", "val2"}
config.Set("key", c1)

v, _ := config.Get("key")
c2 := &conf{}
v.Scan(c2)
// c1 and c2 should be equal
```

Or with the following example

```sh
$ micro config del helloworld
$ micro config set helloworld '{"a":1}'
$ micro config get helloworld
{"a":1}
$ micro config set helloworld '{"b":2}'
$ micro config get helloworld
{"a":1,"b":2}
```

#### Secret encryption keys for `micro server`

By default, if not specified, `micro server` generates and saves an encryption key to the location `~/.micro/config_secret_key`. This is intended for local zero dependency use, but not for production.

To specify the secret for the micro server either the envar `MICRO_CONFIG_SECRET_KEY` or the flag `config_secret_key` key must be specified.

### Errors

The errors package provides error types for most common HTTP status codes, e.g. BadRequest, InternalSeverError etc. It's recommended when returning an error to an RPC handler, one of these errors is used. If any other type of error is returned, it's treated as an InternalSeverError.

Micro API detects these error types and will use them to determine the response status code. For example, if your handler returns errors.BadRequest, the API will return a 400 status code. If no error is returned the API will return the default 200 status code.

Error codes are also used when handling retries. If your service returns a 500 (InternalServerError) or 408 (Timeout) then the client will retry the request. Other status codes are treated as client error and won't be retried.

#### Usage

```go
import (
	"github.com/micro/micro/v3/service/errors"
)

func (u *Users) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("users.Read.MissingID", "Missing ID")
	}

	...
}
```

### Events

The events service is a service for event streaming and persistent storage of events.

#### Overview

Event streaming differs from pubsub messaging in that it provides an ordered stream of events that can be consumed 
or replayed from any given point in the past. If you have experience with Kafka then you know it's basically a 
distributed log which allows you to read a file from different offsets and stream it.

The event service and interface provide the event streaming abstraction for writing and reading events along with 
consuming from any given offset. It also supports acking and error handling where appropriate.

Events also different from the broker in that it provides a fixed Event type where you fill in the details and 
handle the decoding of the message body yourself. Events could have large payloads so we don't want to 
unnecessarily decode where you may just want to hand off to a storage system.

#### Functions

The events package has two parts: Stream and Store. Stream is used to Publish and Consume to messages for a given topic. For example, in a chat application one user would Publish a message and another would subscribe. If you later needed to retrieve messages, you could either replay them using the Subscribe function and passing the Offset option, or list them using the Read function.

```go
func Publish(topic string, msg interface{}, opts ...PublishOption) error 
```
The Publish function has two required arguments: topic and message. Topic is the channel you're publishing the event to, in the case of a chat application this would be the chat id. The message is any struct, e.g. the message being sent to the chat. When the subscriber receives the event they'll be able to unmarshal this object. Publish has two supported options, WithMetadata to pass key/value pairs and WithTimestamp to override the default timestamp on the event.

```go
func Consume(topic string, opts ...ConsumeOption) (<-chan Event, error)
```
The Consume function is used to consume events. In the case of a chat application, the client would pass the chat ID as the topic, and any events published to the stream will be sent to the event channel. Event has an Unmarshal function which can be used to access the message payload, as demonstrated below:

```go
for {
	evChan, err := events.Consume(chatID)
	if err != nil {
		logger.Error("Error subscribing to topic %v: %v", chatID, err)
		return err
	}
	for {
		ev, ok := <- evChan
		if !ok {
			break
		}
		var msg Message
		if err :=ev.Unmarshal(&msg); err != nil {
			logger.Errorf("Error unmarshaling event %v: %v", ev.ID, err)
			return err
		}
		logger.Infof("Received message: %v", msg.Subject)
	}
}
```

#### Example

The [Chat Service](https://github.com/micro/services/tree/master/chat) examples usage of the events service, leveraging both the stream and store functions.

### Network

The network is a service to service network for request proxying

#### Overview

The network provides a service to service networking abstraction that includes proxying, authentication, 
tenancy isolation and makes use of the existing service discovery and routing system. The goal here 
is not to provide service mesh but a higher level control plane for routing that can govern access 
based on the existing system. The network requires every service to be pointed to it, making 
an explicit choice for routing.

Beneath the covers cilium, envoy and other service mesh tools can be used to provide a highly 
resilient mesh.

### Registry

The registry is a service directory and endpoint explorer

#### Overview

The service registry provides a single source of truth for all services and their APIs. All services 
on startup will register their name, version, address and endpoints with the registry service. They 
then periodically re-register to "heartbeat", otherwise being expired based on a pre-defined TTL 
of 90 seconds. 

The goal of the registry is to allow the user to explore APIs and services within a running system.

The simplest form of access is the below command to list services.

```
micro services
```

#### Usage

The get service endpoint returns information about a service including response parameters parameters for endpoints:

```sh
$ micro registry getService --service=helloworld
{
	"services": [
		{
			"name": "helloworld",
			"version": "latest",
			"endpoints": [
				{
					"name": "Helloworld.Call",
					"request": {
						"name": "Request",
						"type": "Request",
						"values": [
							{
								"name": "name",
								"type": "string"
							}
						]
					},
					"response": {
						"name": "Response",
						"type": "Response",
						"values": [
							{
								"name": "msg",
								"type": "string"
							}
						]
					}
				}
			],
			"nodes": [
				{
					"id": "helloworld-3a0d02be-f98e-4d9d-a8fa-24e942580848",
					"address": "192.168.43.193:34321",
					"metadata": {
						"broker": "service",
						"protocol": "grpc",
						"registry": "service",
						"server": "grpc",
						"transport": "grpc"
					}
				}
			],
			"options": {}
		}
	]
}
```

This is an especially useful feature for writing custom meta tools like API explorers.

### Runtime

#### Overview

The runtime service is responsible for running, updating and deleting binaries or containers (depending on the platform - eg. binaries locally, pods on k8s etc) and their logs.

#### Running a service

The `micro run` command tells the runtime to run a service. The following are all valid examples:

```sh
micro run github.com/micro/services/helloworld
micro run .  # deploy local folder to your local micro server
micro run ../path/to/folder # deploy local folder to your local micro server
micro run helloworld # deploy latest version, translates to micro run github.com/micro/services/helloworld or your custom base url
micro run helloworld@9342934e6180 # deploy certain version
micro run helloworld@branchname  # deploy certain branch
micro run --name helloworld .
```

#### Specifying Service Name

The service name is derived from the directory name of your application. In case you want to override this specify the `--name` flag.

```
micro run --name helloworld github.com/myorg/helloworld/server
```

#### Running a local folder

If the first parameter is an existing local folder, ie

```sh
micro run ./foobar
```

Then the CLI will upload that folder to the runtime and the runtime runs that.

#### Running a git source

If the first parameter to `micro run` points to a git repository (be it on GitHub, GitLab, Bitbucket or any other provider), then the address gets sent to the runtime and the runtime downloads the code and runs it.

##### Using references

References are the part of the first parameter passed to run after the `@` sign. It can either be a branch name (no reference means version `latest` which equals to master in git terminology) or a commit hash.

When branch names are passed in, the latest commit of the code will run.

#### Listing runtime objects

The `micro status` command lists all things running in the runtime:

```sh
$ micro status
NAME		VERSION	SOURCE					STATUS	BUILD	UPDATED		METADATA
helloworld	latest	github.com/micro/services/helloworld	running	n/a	20h43m45s ago	owner=admin, group=micro
```

The output includes the error if there is one. Commands like `micro kill`, `micro logs`, `micro update` accept the name returned by the `micro status` as the first parameter (and not the service name as that might differ).

#### Updating a service

The `micro update` command makes the runtime pull the latest commit in the branch and restarts the service.

In case of local code it requires not the runtime name (returned by `micro status`) but the local path. For commit hash deploys it just restarts the service.

Examples: `micro update helloworld`, `micro update helloworld@branch`, `micro update helloworld@commit`, `micro update ./helloworld`.

#### Deleting a service

The `micro kill` command removes a runtime object from the runtime. It accepts the name returned by `micro status`.

Examples: `micro kill helloworld`.

#### Logs

The `micro logs` command shows logs for a runtime object. It accepts the name returned by `micro status`.

The `-f` flag makes the command stream logs continuously.

Examples: `micro logs helloworld`, `micro logs -f helloworld`.

### Store

Micro's store interface is for persistent key-value storage.

For a good beginner level doc on the Store, please see the [Getting started tutorial](/getting-started).

#### Overview

Key-value stores that support ordering of keys can be used to build complex applications.
Due to their very limited feature set, key-value stores generally scale easily and reliably, often linearly with the number of nodes added.

This scalability comes at the expense of inconvenience and mental overhead when writing business logic. For use cases where linear scalability is important, this trade-off is preferred.

#### Query by ID

Reading by ID is the archetypal job for key value stores. Storing data to enable this ID works just like in any other database:

```sh
# entries designed for querying "users by id"
KEY         VALUE
id1         {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
id2         {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
id3         {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
id4         {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
```

```go
import "github.com/micro/micro/v3/service/store"

records, err := store.Read("id1")
if err != nil {
	fmt.Println("Error reading from store: ", err)
}
fmt.Println(records[0].Value)
// Will output {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
```

Given this data structure, we can do two queries:

- reading a given key (get "id1", get "id2")
- if the keys are ordered, we can ask for X number of entries after a key (get 3 entries after "id2")

Finding values in an ordered set is possibly the simplest task we can ask a database.
The problem with the above data structure is that it's not very useful to ask "find me keys coming in the order after "id2". To enable other kinds of queries, the data must be saved with different keys.

In the case of the schoold students, let's say we wan't to list by class. To do this, having the query in mind, we can copy the data over to another table named after the query we want to do:

#### Query by Field Value Equality

```sh
# entries designed for querying "users by class"
KEY             VALUE
firstGrade/id1  {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
thirdGrade/id4  {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
```


```go
import "github.com/micro/micro/v3/service/store"

records, err := store.Read("", store.Prefix("secondGrade"))
if err != nil {
	fmt.Println("Error reading from store: ", err)
}
fmt.Println(records[0].Value)
// Will output
// secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
// secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
```

Since the keys are ordered it is very trivial to get back let's say "all second graders".
Key value stores which have their keys ordered support something similar to "key starts with/key has prefix" query. In the case of second graders, listing all records where the "keys start with `secondGrade`" will give us back all the second graders.

This query is basically a `field equals to` as we essentially did a `field class == secondGrade`. But we could also exploit the ordered nature of the keys to do a value comparison query, ie `field avgScores is less than 90` or `field AvgScores is between 90 and 95` etc., if we model our data appropriately:

#### Query for Field Value Ranges

```sh
# entries designed for querying "users by avgScore"
KEY         VALUE
089/id3     {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
092/id2     {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
094/id4     {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
098/id1     {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
```

It's worth remembering that the keys are strings, and that they are ordered lexicographically. For this reason when dealing with numbering values, we must make sure that they are prepended to the same length appropriately.

At the moment Micro's store does not support this kind of query, this example is only here to hint at future possibilities with the store.

#### Tables Usage

Micro services only have access to one Store table. This means all keys live in the same namespace and can collide. A very useful pattern is to separate the entries by their intended query pattern, ie taking the "users by id" and users by class records above:

```sh
KEY         VALUE
# entries designed for querying "users by id"
usersById/id1         		{"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
usersById/id2         		{"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
usersById/id3         		{"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
usersById/id4         		{"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
# entries designed for querying "users by class"
usersByClass/firstGrade/id1  {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
usersByClass/secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
usersByClass/secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
usersByClass/thirdGrade/id4  {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
```

Respective go examples this way become:

```go
import "github.com/micro/micro/v3/service/store"

const idPrefix = "usersById/"

records, err := store.Read(idPrefix + "id1")
if err != nil {
	fmt.Println("Error reading from store: ", err)
}
fmt.Println(records[0].Value)
// Will output {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
```

```go
import "github.com/micro/micro/v3/service/store"

const classPrefix = "usersByClass/"

records, err := store.Read("", store.Prefix(classPrefix + "secondGrade"))
if err != nil {
	fmt.Println("Error reading from store: ", err)
}
fmt.Println(records[0].Value)
// Will output
// secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
// secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
```

### Metadata

Metadata / headers can be passed via the context in RPC calls. The context/metadata package allows services to get and set metadata in a context. The Micro API will add request headers into context, for example if the "Foobar" header is set on an API call to "localhost:8080/users/List", the users service can access this value as follows:

```go
import (
	"context"
	"github.com/micro/micro/v3/service/context/metadata"
)

...

func (u *Users) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	val, ok := metadata.Get(ctx, "Foobar")
	if !ok {
		return fmt.Errorf("Missing Foobar header")
	}

	fmt.Println("Foobar header was set to: %v", val)
	return nil
}
```

Likewise, clients can set metadata in context using the metadata.Set function as follows:

```go
func (u *Users) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	newCtx := metadata.Set(ctx, "Foobar", "mycustomval")
	fRsp, err := u.foosrv.Call(newCtx, &foosrv.Request{})
	...
}
```

## Plugins

Micro is a pluggable architecture built on Go's interface types. Plugins enable swapping out underlying infrastructure.

### Overview

Micro is pluggable, meaning the implementation for each module can be replaced depending on the requirements. Plugins are applied to the micro server and not to services directly, this is done so the underlying infrastructure can change with zero code changes required in your services. 

An example of a pluggable interface is the store. Locally micro will use a filestore to persist data, this is great because it requires zero dependencies and still offers persistence between restarts. When running micro in a test suite, this could be swapped to an in-memory cache which is better suited as it offers consistency between runs. In production, this can be swapped out for standalone infrastructure such as cockroachdb or etcd depending on the requirement.

Let's take an example where our service wants to load data from the store. Our service would call `store.Read(userPrefix + userID)` to load the value, behind the scenes this will execute an RPC to the store service which will in-turn call `store.Read` on the current `DefaultStore` implementation configured for the server. 

### Profiles

Profiles are used to configure multiple plugins at once. Micro comes with a few profiles out the box, such as "local", "kubernetes" and "test". These profiles can be found in `profile/profile.go`. You can configure micro to use one of these profiles using the `MICRO_PROFILE` env var, for example: `MICRO_PROFILE=test micro server`. The default profile used is "local".

### Writing a profile

Profiles should be created as packages within the profile directory. Let's create a "staging" profile by creating `profile/staging/staging.go`. The example below shows how to override the default store to use an in-memory implementation:

```go
// Package staging configures micro for a staging environment
package staging

import (
	"github.com/urfave/cli/v2"

	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"
)

func init() {
	profile.Register("staging", staging)
}

var staging = &profile.Profile{
	Name: "staging",
	Setup: func(ctx *cli.Context) error {
		store.DefaultStore = memory.NewStore()
		return nil
	},
}
```

### Using a custom profile

You can load a custom profile using a couple of commands, the first adds a replace to your go mod, indicating it should look for your custom profile within the profile directory:

```bash
go mod edit -replace github.com/micro/micro/profile/staging/v3=./profile/staging
```

The second command creates a profile.go file which imports your profile. When your profile is imported, the init() function which is defined in staging.go is called, registering your profile.

```
micro init --profile=staging --output=profile.go
```

Now you can start your server using this profile:
```
MICRO_PROFILE=staging go run . server
```
