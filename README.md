# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)

Micro is a runtime for microservice development.

## Overview

Micro addresses the key requirements for building scalable systems. It takes the microservice architecture pattern and transforms it into 
a set of tools which act as the building blocks of a platform. Micro deals with the complexity of distributed systems and provides 
simple abstractions already understood by developers.

<img src="https://micro.mu/runtime3.svg" />

Technology is constantly evolving. The infrastructure stack is always changing. Micro is a pluggable platform which addresses these issues. 
Plug in any stack or underlying technology. Build future-proof systems using micro.

## Features

The runtime is composed of the following features:

- **api**: An api gateway. A single entry point with dynamic request routing using service discovery. The API gateway allows you to build a scalable 
microservice architecture on the backend and consolidate serving a public api on the frontend. The micro api provides powerful routing 
via discovery and pluggable handlers to serve http, grpc, websockets, publish events and more.

- **bot:** A slackbot which runs on your platform and lets you manage your applications from Slack itself. The micro bot enables ChatOps 
and gives you the ability to do everything with your team via messaging. It also includes ability to create slack commands as services which 
are discovered dynamically.

- **cli:** An interactive CLI to describe, query and interact directly with your platform and services from the terminal. The CLI 
gives you all the commands you expect to understand what's happening with your micro services. It also includes an interactive mode.

- **network:** Build multi-cloud networks with the micro network service. Simply drop-in and connect the network services across any environment 
and create a single flat network to route globally. The micro network dynamically builds routes based on your local registry in each datacenter 
ensuring queries are routed based on locality.

- **new:** A service template generator. Create new service templates to get started quickly. Micro provides predefined templates for writing micro services. 
Always start in the same way, build identical services to be more productive.

- **proxy:** A transparent service proxy built on [Go Micro](https://github.com/micro/go-micro). Offload service discovery, load balancing, 
fault tolerance, message encoding, middleware, monitoring and more to a single a location. Run it standalone or alongside your service.

- **tunnel:** A network tunnel to get access to services in any environment without the need for a vpn. The micro tunnel provides point to 
point tunnelling with a built-in proxy to query services in remote environments. Query production systems from your local laptop. 

- **web:** The web dashboard allows you to explore your services, describe their endpoints, the request and response formats and even 
query them directly. The dashboard also includes a built in CLI like experience for developers who want to drop into the terminal on the fly.

Additionally micro provides a Go development framework:

- **go-micro:** Leverage the powerful [Go Micro](https://github.com/micro/go-micro) framework to develop microservices easily and quickly. 
Go Micro abstracts away the complexity of distributed systems and provides simpler abstractions to build highly scalable microservices.

## Install

From source

```
go get github.com/micro/micro
```

Docker image

```
docker pull micro/micro
```


## Getting Started

Boot the entire system and connect to the network

```
micro
```

Run without connecting to the network

```
micro --network=local
```

### Create a service

```
# enable go modules
export GO111MODULE=on

# generate a service (follow instructions in output)
micro new example

# run the service
go run example/main.go

# list services
micro list services

# call a service
micro call go.micro.srv.example Example.Call '{"name": "John"}'
```

## Use the network

The micro network is a shared global services network actively in development.

Proxy service calls through the network

```
export MICRO_PROXY=go.micro.network
```

View network services, routes, nodes

```
# List nodes
micro network nodes

# List services
micro network services

# List routes
micro network routes

# Peer graph
micro network graph
```

## Usage

See all the options

```
micro --help
```

See the [docs](https://micro.mu/docs/) for detailed information on the architecture, installation and use of the platform.
