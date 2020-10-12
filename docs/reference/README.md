---
title: Reference
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference
summary: Reference - a comprehensive guide to Micro
---

## Reference

Reference entries are in depth look at the technical details and usage of Micro

## Contents

* TOC
{:toc}

## CLI Overview

Micro is driven entirely through a CLI experience. This reference highlights the CLI design.

The CLI speaks to the `micro server` through the gRPC proxy running locally by default on :8081. All requests are proxied based on your environment 
configuration. The CLI provides the sole interaction for controlling services and environments.

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
running. By default locally this will not exist and we expect the user to use the admin/micro credentials to administrate the system. 
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
$ micro env add myown stunningproject.com
$ micro env
* local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
  foobar    example.com
```

### Set Environment

The `*` marks wich environment is selected. Let's select the newly added:

```sh
$ micro env set myown
$ micro env
$ micro env
  local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
* foobar     example.com
```

### Login to an Environment

Each environment is effectively an isolated deployment with its own authentication, storage, etc. So each env requires signup and login. At this point we have to log in to the `example` env with `micro login`. If you don't have the credentials to the environment, you have to ask the admin.

## Installation

### Helm

Micro can be installed onto a Kubernetes cluster using helm. Micro will be deployed in full and leverage zero-dep implementations designed for Kubernetes. For example, micro store will internally leverage a file store on a persistant volume, meaning there are no infrastructure dependancies required.

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

### Local 

Micro can be installed locally in the following way. We assume for the most part a Linux env with Go and Git installed.

#### Go Get

```
go get github.com/micro/micro/v3
```

#### Docker

```
docker pull micro/micro
```

#### Release Binaries

```
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

## Server

The micro service is a distributed systems runtime for the Cloud and beyond. It provides the building 
blocks for distributed systems development as a set of microservices and framework.

### Usage

To start the server simply run

```
micro server
```

This will boot the entire system and services including a http api on :8080 and grpc proxy on :8081

### Verify Status

Check help text is output with no errors
```
micro --help
```

Run helloworld

```
micro env	# should point to local
micro status	# returns empty response
micro services	# returns empty response
micro run github.com/micro/services/helloworld # run helloworld
micro status 	# wait for status running
```

Call the service and verify output

```shell
$ micro helloworld --name=John
{
        "msg": "Hello John"
}
```

Remove the service

```
micro kill helloworld
```

## Services

The Micro Server is not a monolithic process. Instead it is composed of many separate services.

Below we describe the list of services provided by the Micro Server. Each service is considered a 
building block primitive for a platform and distributed systems development. The proto 
interfaces for each can be found in [micro/proto/auth](https://github.com/micro/micro/blob/master/proto/auth/auth.proto) 
and the Go library, client and server implementations in [micro/service/auth](https://github.com/micro/micro/tree/master/service/auth).

### Auth

The auth service provides both authentication and authorization.
The auth service stores accounts and access rules. It provides the single source of truth for all authentication 
and authorization within the Micro runtime. Every service and user requires an account to operate. When a service 
is started by the runtime an account is generated for it. Core services and services run by Micro load rules 
periodically and manage the access to their resources on a per request basis.

### Config

The config service provides dynamic configuration for services. Config can be stored and loaded separately to 
the application itself for configuring business logic, api keys, etc. We read and write these as key-value 
pairs which also support nesting of JSON values. The config interface also supports storing secrets by 
defining the secret key as an option at the time of writing the value.

### Broker

TODO

### Events

TODO

### Network

TODO

### Registry

TODO

### Runtime

TODO

### Store

Micro's store interface is for persistent key-value storage.

For a good beginner level doc on the Store, please see the [helloworld tutorial](/helloworld).

#### Overview

Key-value stores that support ordering of keys can be used to build complex applications.
Due to their very limited feature set, key-value stores generally scale easily and reliably, often linearly with the number of nodes added.

This scalability comes at the expense of inconvenience and mental overhead when writing business logic. For usecases where linear scalability is important, this tradeoff is preferred.

#### Query by ID

Reading by ID is the archetypical job for key value stores. Storing data to enable this ID works just like in any other database:

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
The problem with the above data structure is that it's not very useful to ask "find me keys coming in the order after "id2". To enable other kind of queries, the data must be saved with different keys.

In the case of the schoold students, let's say we wan't to list by class. To do this, having the query in mind, we can copy the data over to an other table named after the query we want to do:

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

Micro services only have access to one Store table. This means all keys take live in the same namespace and can collide. A very useful pattern is to separate the entries by their intended query pattern, ie taking the "users by id" and users by class records above:

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
