# Framework

Go Micro is a framework for microservices development

## Overview

Go Micro provides the fundamental building blocks for creating distributed systems. It serves primarily for writing microservices but can be used to construct any kind of low level distributed system such as a database or network. Go Micro 
is a pluggable framework using Go interfaces.

## Design

Below are the packages and interfaces which we're developing or planning to add in the future.

Package	|	Function	|	Description
-------	|	--------	|	-----------
[Service](https://godoc.org/github.com/micro/go-micro#Service)	|	Communication	| Request/Response, Streaming, PubSub
[Auth](https://godoc.org/github.com/micro/go-micro/auth) | Authentication | Authentication and authorization
[Config](https://godoc.org/github.com/micro/go-micro/config)	|	Configuration	|	Dynamic config, safe fallbacks, etc
[Debug](https://godoc.org/github.com/micro/go-micro/debug)	|	Debugging	|	Logging, tracing, metrics, healthchecks
[Network](https://godoc.org/github.com/micro/go-micro/network) |	Networking	|	Multi-DC networking
[Runtime](https://godoc.org/github.com/micro/go-micro/runtime)	|	Runtime	|	Service runtime
[Sync](https://godoc.org/github.com/micro/go-micro/sync)	|	Synchronisation	|	Locking, leadership election
Events  | Event Streaming | Event streaming and timeseries database
Flow |	Orchestration	|	State machine for managing complex workflows of business logic
Model | Data Model  | A data modeling and crud interface


## Micro

Micro is use for building microservices or more practically distributed systems. It's core concern is communication. 
It's pluggable and runtime agnostic with very simple abstractions for development. 

On the road to v2 defaults should look to support gRPC and kubernetes more natively along with graphql and nats.

Supported registries:
- [x] Multicast DNS (Local)
- [x] Serf gossip (P2P / Mesh)
- [x] Consul (Distributed)

Supported brokers:
- [x] Memory (Local / In Process)
- [x] HTTP (P2P / Registry) => move to grpc?
- [x] NATS (Distributed)

Supported transports:
- [x] HTTP (Local / http/1.1)
- [x] gRPC (P2P)
- [x] Service mesh

## Config

Config is for dynamic configuration. Its used where higher level app configuration is required in process without having to reload 
or restart. Its useful for things related to business logic. 

## Sync

Sync is for distributed synchronisation in all forms; data, time, leadership. Distributed systems and microservices require 
coordination. From a development perspective this can often be difficult. Providing clear abstractions for these that 
are also possible to leverage without external dependencies is valuable. 

Sync ideally becomes the basis for all data storage and any form of synchronisation.

## Runtime

Runtime is the basis for `micro run service`. The library implements the ability to pull the source, build and run the service. 
Ideally a service should be able to declare its own dependencies and they should bootstrap as a DAG. 

Supported sources:
- [x] Git URL

Supported packagers:
- [x] Go binary
- [x] Docker container
- [ ] WASM

Supported runtimes:
- [x] Linux process
- [ ] Docker
- [x] Kubernetes API


