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

- **API Gateway:** A single entry point with dynamic request routing using service discovery. The API gateway allows you to build a scalable 
microservice architecture on the backend and consolidate serving a public api on the frontend. The micro api provides powerful routing 
via discovery and pluggable handlers to serve http, grpc, websockets, publish events and more.

- **Interactive CLI:** A CLI to describe, query and interact directly with your platform and services from the terminal. The CLI 
gives you all the commands you expect to understand what's happening with your micro services. It also includes an interactive mode.

- **Service Proxy:** A transparent proxy built on [Go Micro](https://github.com/micro/go-micro). Offload service discovery, load balancing, 
fault tolerance, message encoding, middleware, monitoring and more to a single a location. Run it standalone or alongside your service.

- **Template Generation:** Create new service templates to get started quickly. Micro provides predefined templates for writing micro services. 
Always start in the same way, build identical services to be more productive.

- **Slack Bot:** A bot which runs on your platform and lets you manage your applications from Slack itself. The micro bot enables ChatOps 
and gives you the ability to do everything with your team via messaging. It also includes ability to create slack commmands as services which 
are discovered dynamically.

- **Web Dashboard:** The web dashboard allows you to explore your services, describe their endpoints, the request and response formats and even 
query them directly. The dashboard also includes a built in CLI like experience for developers who want to drop into the terminal on the fly.

- **Go Framework:** Leverage the powerful [Go Micro](https://github.com/micro/go-micro) framework to develop microservices easily and quickly. 
Go Micro abstracts away the complexity of distributed systems and provides simpler abstractions to build highly scalable microservices.

## Getting Started

See the [docs](https://micro.mu/docs/) for detailed information on the architecture, installation and use of the platform.

### Install

```
go get -u github.com/micro/micro
```

### Create a service

```
micro new github.com/my/service
```
