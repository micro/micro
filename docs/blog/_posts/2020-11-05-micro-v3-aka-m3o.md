---
layout:	post
author: Asim Aslam
title:	"M3O aka Micro 3.0 is a platform for cloud native development"
date:	2020-11-05 10:00:00
---

This is the official announcement for the release of Micro 3.0 better known as M3O. Our 3.0 release is a 
major refactor, overhaul and consolidation of the existing tooling into something that looks more like a platform for 
cloud native development. 

Read on to learn more or go straight to the 
[latest release](https://github.com/micro/micro/releases/latest). Head to [m3o.com](https://m3o.com) for 
the hosted offering.

## Overview

Micro focuses on developer productivity for the backend. It's clear that the Cloud has become infinitely more complex 
over the past few years. Micro attempts to create order out of that chaos by distilling it all down to a handful of 
primitives for distributed systems development.

Why should you care? If you're reading this you've no doubt encountered the tedious nature of infrastructure management, 
wrangling a kubernetes cluster on AWS or the thousands of things you need to do to cobble together a platform before 
starting to build a product. We think we've nailed the solution for that just as Android did for Mobile. Keep reading 
if you want to find out more.

## Quick Flashback

Micro started out as a [toolkit for microservices](/blog/2016/03/20/micro.html) development, 
incorporating an api gateway, web dashboard and cli to interact with services built using a Go RPC framework. 
Back then it felt like getting anyone to buy into PaaS again was going to be a losing battle. So we chose 
to write single purpose tools around an RPC framework thinking it might allow people to adopt it piece by piece 
until they saw the need for a platform. It was really straight forward right until it wasn't. 

There was a simple Go framework plus some surrounding 
components to query and interact with them, but like any long lived project, the complexity grew as we 
tried to solve for that platform experience that just couldn't be done with a swiss army knife. The repo 
exploded with a number of independent libraries. To the creator its obvious what these are all for but to 
the user there is nothing but cognitive overload. 

In 2019 we went through a [consolidation](/blog/2019/06/10/the-great-consolidation.html) of all those libraries 
which helped tremendously but there was still always one outstanding question. What's the difference between 
[micro](https://github.com/micro/micro) and [go-micro](https://github.com/micro/go-micro)? It's a good 
question and one we've covered before. We saw go-micro as a framework and micro as a toolkit but these 
words were basically empty and meaningless because multiple projects working in coordination really need a 
crisp story that makes sense and we didn't have one.

In 2020 we're looking to rectify that but let's first let's talk about platforms.

## PaaS in 2020

5 years ago the world exploded with a proliferation of "cloud native" tooling as containers and 
container orchestration took centre stage. More specifically, Docker and Kubernetes redefined the 
technology landscape along with a more conscious move towards building software in the cloud.

Micro took a forward looking view even as far back as 2015. It was clear distributed systems and cloud native 
was going to become the dominant model for backend services development over the coming years but, what wasn't clear 
is just how long we'd spend wrangling all sorts tools like docker, kubernetes, grpc, istio and everything else. 
It felt like we were rebuilding the stack and weren't really ready to talk about development aspects of it all.

In fact at that time, people mostly wanted to kick the tyres on all these tools and piece something together. 
Running kubernetes yourself became all the rage and even using service mesh as the holy grail for solving 
all your distributed systems problems. Many of us have come to realise while all of this tech is fun 
it's not actually solving development problems.  

We've gotten to the point of managed kubernetes and even things like Google Cloud Run or DigitalOcean App 
Platform, but none of these things are helping with a development model for a cloud native era. Our 
frustrations with the existing developer experience have grown and Micro felt like something that 
could solve for all that, but only if we took a drastic step to overhaul it.

We think PaaS 3.0 is not just about running your container or even your source code but something that 
encapsulates the entire developer experience including a model for writing code for the cloud. Based on that 
Micro 3.0 aka M3O is a platform for cloud native development.

## What even is Cloud Native?

What is cloud native? What does it mean to build for the cloud? What is a cloud service?

Cloud native is basically a descriptive term for something that was built to run in the cloud. That's it. It's not 
magic, it might sound like a buzzword, but the reality is it simply means, that piece of software was built 
to run in the cloud. How does that differ from the way we used to build before? Well the idea behind the cloud 
is that its ephemeral, scalable and everything can be accessed via an API.

Our expectation for services running in the cloud is that they're mostly stateless, leveraging external services 
for the persistence, that they are identified by name rather than IP address and they themselves provide an 
API that can be consumed by multiple clients such as web, mobile and cli or other services. 

Cloud native applications are horizontally scalable and operate within domain boundaries that divide them as 
separate apps which communicate over the network via their APIs rather than as one monolithic entity. 
We think cloud services require a fundamentally different approach to software creation and why Micro 3.0 
was designed with this in mind.

## M3O aka Micro 3.0

Micro 3.0 reimagines Micro as a platform for cloud native development. What does that mean? Well we think of 
it as PaaS 3.0, a complete solution for source to running and beyond. Micro has moved from just being a Go 
framework to incorporating a standalone server and hosted platform. Our hosted offering is called 
[M3O](https://m3o.com), a hat tip to Micro 3.0 or M[icr]o, whichever way you want to see it. 

### Server

The server is our abstraction for cloud infrastructure and underlying systems you might need for writing 
distributed systems. The server encapsulates all of these concerns as gRPC services which you can 
query via any language. The goal here is to say developers don't really need to be thinking about infrastructure 
but what they do need is design patterns and primitives for building distributed systems. 

The server includes the following:

- **Authentication**: Auth whether its authentication or authorization is part of the system. Create JWT tokens, define access rules, use one system to govern everything in a simple and straight forward manner. Whether it’s for a user or a service.

- **Configuration**: Dynamic config management allows you to store relevant config that needs to be updated without having to restart services. Throw API keys and business logic related configuration into the secure config service and let your services pick up the changes.

- **Key-Value Storage**: We’re focused on best practices for microservices development which means keeping services mostly stateless. To do this we’re providing persistent storage on the platform. Key-Value allows you to rapidly write code and store data in the format you care about.

- **Event Streaming**: Distributed systems are fundamentally in need of an event driven architecture to breakdown the tight dependencies between them. Using event streaming and pubsub allows you to publish and subscribe to relevant events async.

- **Service Registry**: Micro and M3O bake in service discovery so you can browse a directory of services to explore your service APIs and enable you to query services by name. Micro is all about microservices and multi-service development.

- **Service Network**: Because you don't want to have to resolve those service names to addresses and deal with the load balancing aspect, the server bakes in a "service mesh" which will handle your inter-service requests (as gRPC) and route to the 
appropriate instance.

- **Identity Proxy**: We include a separate identity proxy for external requests using gRPC via the CLI and other means. This enables you to query from your local machine or anywhere else using valid auth credentials and have it seamlessly work as if 
you were in the platform itself.

- **API Gateway**: Finally there’s an API gateway that automatically exposes your services to the outside world over HTTP. Internally writing service to service using gRPC makes sense, but at the end of the day we want to build APIs consumed from clients via HTTP.

### Clients

The server provides inter-service communication and two means of external communication with a HTTP API and gRPC proxy but that 
experience is made much better when there's user experience on the client side that works. Right now we've got two ways of doing this.

- **Command Line**: The CLI provides a convenient and simple way to talk to the server via gRPC requests through the proxy. 
The most convenient commands are builtin but every service you write also gets beautiful dynamic generated commands 
for each endpoint. 

- **gRPC SDKs**: Every service in the server is accessible via gRPC. We're code generating clients for the server itself 
so you can access them from any language. What this enables is a wide array of experiences on the client side without 
having to handcraft libraries for each language.

- **Web Interface**: Coming soon is a dynamically generated web interface that creates a simple query mechanism through a 
browser for any of your services. We've got a http api, gRPC proxy and command line interface but feel like the browser 
could use some love too.

### Library

One thing we really understood from our time working on go-micro was that the developer experience really matters. We 
see Go as the dominant language for the cloud and believe most backend services in the cloud will be written in Go. For 
that reason we continue to include a Service Library which acts as a framework for building your services and accessing 
the underlying systems of the server.

The Service Library provides pre-initialised packages for all of the features of the server and creates a convenient 
initialiser for defining your own services starting with `service.New`. A Service has a name, endpoints, contains 
a server of its own and a client to query other services. The library does enough for you but then attempts to 
get out of your way so the rest is up to you.

A main package for a Micro service looks something like this

```
package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/helloworld/handler"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("helloworld"),
	)

	// Register Handler
	srv.Handle(new(handler.Helloworld))

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
```

When you want to make use of something like the Config service just import it like so.

```
import "github.com/micro/micro/v3/service/config"

val, err := config.Get("key")
```

You can find many more examples in [github.com/micro/services](https://github.com/micro/services).

### Environments

From our experience writing software isn't constraing to a single environment. Most of the time we're doing 
some form of local development followed by a push to staging and then production. We don't really see tools 
capturing that workflow effectively. Thinking about how to do this now we've built in environments as 
a first class system.

The platform offers 3 builtin environments for now; local, dev and platform.

- **Local** - is Micro running on your local machine
- **Dev** - is a free development environment in the cloud
- **Platform** - is a paid secure, scalable and supported production environment

Our goal here is to really direct the flow from local > dev > platform as the lifecycle for any backend service 
development. Start by running the server locally, writing your code and getting it to work. Ship it to the 
dev environment for further testing but also to collaborate with others and serve it publicly. Then if you're 
interested in a scalable and supported production environment, pay for the platform environment. That's it.

Interact with the environments like so.

```
# view the environments
micro env

# set the environment
micro env set platform

# add a new environment
micro env add foobar proxy.foo.com:443
```

### Multi-Tenancy

With the advent of a system like kubernetes and a push towards the cloud we can see that there's really a need 
to move towards shared resource usage. The cloud isn't cheap and we don't all need to be running separate 
kubernetes clusters. In fact wouldn't it be great if we could share that? Well Micro is doing it. We build in 
multi-tenancy using the same logic kubernetes does called Namespaces. 

We've mapped this same experience locally so you get a rudimentary form of namespacing for local dev but 
mostly we're making use of kubernetes namespaces in production along with a whole host of custom written 
isolation mechanisms for authentication, storage, configuration, event streaming, etc so Micro 3.0 
can be used to host more than one tenant.

Whether you decide to self host and share your cluster for dev, staging and production we felt like 
multi-tenancy needs to become a defacto standard in 2020. How it works in practice. Each tenant 
get's a namespace. That namespace has its own isolated set of users and resources in each subsystem. 
When you make any request as a user or service, a JWT token is passed with that so the 
underlying systems can route to the appropriate resources.

### Source to Running 

Micro was built out of a frustration with the existing tools out there. One of the things I've really 
been saying for a long time is that I wanted "source to running" in just one command. With 
Heroku we sort of got that but it really took too much away from us. Back in 2010 Heroku was focused 
on monolithic Rails development. Since then I've really said Heroku took too much away and AWS gave 
too much back. We needed something in between.

Micro can take your source code, from a local directory or a repo thats hosted on github, gitlab or bitbucket. 
In one command it will upload or pull from the relevant place, package it as a container and 
run it. That's it. Source to running in just one command. No more need to deal with the pipeline, 
no more hacking away at containers and the container registries. Write some code and run it.

### Development Model

Source to running is cool. It's what a PaaS is really for but one thing that's really been lacking even 
with the new PaaS boom is a development model. As I eluded to, Heroku takes too much away and AWS 
gives too much back. We're looking for a happy medium. One that doesn't require us to rely on VMs or 
containers but on the other side doesn't limit us to monolithic development.

Micro has always focused on the practice of distributed systems development or microservices. The idea 
of breaking down large monolithic apps into smaller separate services that do one thing well. To do 
this we think you really have to bake the development model into the platform.

What we include is the concept of a Service which contains a Client and Server for both handling 
requests and making queries to other services. We focus on standardisation around protobuf for API 
definitions and using gRPC for the networking layering. 

Not only that we're including pubsub for an event streaming architecture and other pieces like 
nosql key-value storage and dynamic config management. We believe there are specific primitives 
required to start building microservices and distributed systems and that's what Micro looks 
to provide.

### Ten Commands

Alright so we talk a good game, but how easy is it? Well lets show you.

```
# Install the micro binary
curl -fsSL https://install.m3o.com/micro | /bin/bash

# Set the env to platform (paid) or dev (free)
micro env set platform

# Signup before getting started
micro signup

# Create a new service (follow the instructions and push to Github)
micro new helloworld

# Deploy the service from github
micro run github.com/micro/services/helloworld

# Check the service status
micro status

# Query the logs
micro logs helloworld

# Call the service
micro helloworld --name=Alice

# Get your namespace
NAMESPACE=$(micro user namespace)

# Call service via the public http API
curl "https://$NAMESPACE.m3o.app/helloworld?name=Alice"
```

Easy right? We see this as the common flow for most service development. Its a fast iterative loop
from generating a new template to shipping it and querying to make sure it works. There's 
additional stuff in the developer experience like actually writing the service but we think that's 
a separate post.

### Local Environment

If you're using the server locally a few subsitutions.

Set your environment to the local server

```
micro env set local
```

Curl `localhost:8080` rather than m3o.app with your namespace

```
curl -H "Micro-Namespace: $NAMESPACE" "http://localhost:8080/helloworld?name=Alice"
```

### Dev Environment

If you're using the dev environment URLs are `*.m3o.dev`. Find more details at [m3o.dev](https://m3o.dev)

## Documentation

Another thing we really learned from the past is nothing like this works without great documentation 
and tutorials. So we've written a whole suite of docs for Micro available at [micro.mu](https://micro.mu) 
and provide help for using the M3O on [m3o.dev](https://m3o.dev/).

You can find other interesting resources at [Awesome Micro](https://github.com/micro/awesome-micro).

## Licensing

Micro continues to remain open source but licensed using [Polyform Shield](https://polyformproject.org/licenses/shield/1.0.0/) 
which prevents the software for being picked up and run as a service. This is to contend with AWS and others 
running open source for profit without contributing back. It's a longer conversation for another day.

## Motivations

We really believe that writing software for the cloud is too hard. That there's far too much choice and 
time wasted focusing on how to piece everything together. There are tradeoffs to adopting a PaaS but 
ultimately our focus is **developer productivity**. By choosing one tool and one way we stop thinking 
about the how and just get down to what we're trying to build.

M3O and Micro 3.0 look at the state of distributed systems development in the cloud native era and 
try to drastically simplify that experience with a platform that bakes in the development model 
so you can just get back to writing code.

## Go Micro

We will now be ending support for [go-micro](https://github.com/micro/go-micro). Having personally 
spent 6 years since inception on go-micro I feel as though its time to finally let it go. What 
started as a tiny library to help write Kubernetes-as-a-Service back in 2014 turned into a widely 
used open source framework for Go microservices development. Having now amassed more than 14k stars 
you might wonder why we leave it behind. The truth is, while it solved a problem for many it never 
became what I had expected it to be.

Go Micro was built on the premise that developers need a simpler way to build distributed systems. 
With strongly defined abstractions and a pluggable architecture it did that well but that became 
really unweidly to manage. With an MxN matrix of complexity, Go Micro became the thing it was 
trying to fight against. As we attempted to hone on this platform effort, it just became very 
clear that to do that we'd need to start fresh.

Go Micro will live on as an independent library under my own personal account on GitHub but 
it will no longer be supported as an official Micro project. Hopefully it finds second life in 
some other ways but for now we say goodbye.

## Next Steps

You can use the Micro 3.O as a self-hosted open source solution locally, on a VPS or managed kubernetes, 
whatever works for you. Our goal is to facilitate a vastly superior developer experience for building 
services in the Cloud. Come join [Discord](https://discord.gg/hbmJEct) or [Slack](https://slack.m3o.com) 
to chat more about it. And lastly head to to [m3o.com](https://m3o.com) if you're tired of the way you're 
building software for today and want to learn of a better way that's going to make you 10x more productive.

<br>
*Written by Asim Aslam*
