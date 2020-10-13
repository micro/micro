---
title: Architecture
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /architecture
summary: Architecture - Design, goals and tradeoffs
---

## Architecture
{: .no_toc }

This document covers the architecture of micro, design decisions and tradeoffs made.

## Contents

* TOC
{:toc}

## Overview

Micro is a framework for cloud native development. It takes the concept of distributed systems and codifies it as a software 
design pattern using microservices and related primitives. The overall goal of Micro is to abstract away cloud infrastructure 
and to define a set of building blocks which can be used to write cloud services aka microservices, APIs or distributed systems.

Micro in v3 has undergone a major overhaul, it encompasses three things:

- Server - A single server which acts as the runtime for a cloud platform
- Clients - Entrypoints via command line, api gateway and gRPC proxy/sdks
- Library - A Go service library specifically designed to write Micro services

For in-depth material see the [Reference](/reference). This doc will otherwise cover things at a high level.

## Server

The server acts as an abstraction for the underlying infrastructure. It provides distributed systems primitives as building 
blocks for writing microservices. We'll outline those below. The rationale behind defining a single server is based on 
the understanding that all systems at scale inevitably need a platform with common building blocks. In the beginning this 
supports monolithic apps for building CRUD and over time evolves to firstly separate the frontend from backend and then 
starts to provide scalable systems for persistence, events, config, auth, etc.

The server encapsulates all these concerns while embracing the distributed systems model and running each core concern 
as a separate service and process. Its the unix philosophy done well for software composition. As a whole though 
micro is a single monolithic codebase after many separate libraries were consolidated and inevitably go-micro got 
folded into micro itself.

### Services

The services provided by the server are below. These are identified as the core concerns for a platform and distributed systems 
development. In time this may evolve to include synchronization (locking & leadership election), sagas pattern, etc but for 
now we want to provide just the core primitives.

- API - public http api gateway served on port :8080
- Auth - authentication and authorization
- Broker - pubsub messaging for async comms
- Config - dynamic configuration and secrets
- Events - event streaming and persistent messaging
- Network - service to service communication
- Proxy - identity aware proxy for external access
- Registry - service discovery and api explorer
- Runtime - service runtime and process manager
- Store - key-value persistent storage

The assumptions made are by building as independent services we can scope the concerns and domain boundaries appropriately but 
cohesively build a whole system that acts as a platform or a cloud operating system. Much like Android for Mobile we think 
Micro could become a definitive server for the Cloud.

## Clients

Clients are effectively entrypoints or forms of access for the Micro server. The server runs both an API gateway and gRPC proxy. 
The API is deemed as a public facing http api gateway that converts http/json with path based routing to rpc requests. The 
world still built around http and cloud services emerging with that pattern, it feels like this is the right approach.

The gRPC proxy is a forward looking entrypoint, where we see it as a way of extending the service-to-service communication 
to the developers laptop or other environments. At the moment it acts as a proxy for command line to remote environments.

The command line interface itself is deemed as a client and the defacto way to interact with Micro. We do not offer a web UI. 
The assumption is web is a context switch and we'd prefer developers to stay in the terminal. The CLI is extensible in the 
sense that every service you run becomes a subcommand. By doing so we take the API path based model and service decomposition 
to the CLI as well.

The final piece that is a continuous work in progress is gRPC generated clients which can be found in micro/client/sdk. 
Eventually they will be published to various package managers but for now are all routed in one directory. These are built from 
the protos defined for each core service within Micro, enabling multi-language access. We do not assume services to be built 
multi-language but consumption of Micro and services may extend outward.

## Library

The Go service library is a core piece which comes from the original go-micro framework started in 2015. This framework offered 
core distributed systems primitives as Go interfaces and made them pluggable. With its complexity and overlap with Micro we 
decided the best thing was to merge the two and create a Service Library within Micro to define the defacto standard for building 
Micro Services in Go. The service library provides pluggable abstractions with pre-initialised defaults. The developer 
will import and use any of these libraries without any initialisation, they in turn speak to the micro server or basically 
the core services via gRPC. 

For the developer, this is their main point of interaction when writing code. We employ a build, run, manage philosophy where 
build actually starts with putting something in the hands of the developer. Import the service library, start writing code, 
import various packages as needed when you have to get config, check auth, store key-value data. Then run your service using 
Micro itself. The server abstracts away the infra, the service is built to run on Micro and everything else is taken care of.

