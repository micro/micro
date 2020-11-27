---
title: Introduction
keywords: micro
tags: [micro, overview]
sidebar: home_sidebar
permalink: /introduction
summary: Introduction - A high level introduction to Micro
---

## Introduction
{: .no_toc }

This is a high level introduction to Micro

## Contents
{: .no_toc }

* TOC
{:toc}

## About

Micro is a platform for cloud native development. It addresses the key requirements for building services in the cloud. 
Micro leverages the microservices architecture pattern and provides a set of services which act as the building blocks of a 
platform. Micro deals with the complexity of distributed systems and provides simpler programmable abstractions to build on.

<img src="{{ site.baseurl }}/images/micro-3.0.png" />

Micro is the all encompassing end to end platform experience from source to running and beyond built with a developer first focus.

## Goals

Micro's goal is to abstract away the complexity of building services for the Cloud. The cloud itself has gone through a huge 
boom through managed Compute and infrastructure services from the likes of AWS and others. It's taken what was an operational 
burden and turned it into a suite of fully managed on demand services which can be used via APIs.

This has opened up the category of Cloud as Mobile did in the early 2000s but it has yet to define a development model to 
effectively leverage these services. In fact we think the definition of a "Cloud Service" is yet to be determined and 
the category of Cloud will likely shift to a vertically integrated solution that looks more like an operating system 
bundled with cloud infrastructure rather than self installation and management.

Think of Micro as Android for Cloud.

## Features

Below are the core components that make up Micro.

**Server**

Micro is built as a microservices architecture and abstracts away the complexity of the underlying infrastructure. We compose 
this as a single logical server to the user but decompose that into the various building block primitives that can be plugged 
into any underlying system. 

The server is composed of the following services.

- **API** - HTTP Gateway which dynamically maps http/json requests to RPC using path based resolution
- **Auth** - Authentication and authorization out of the box using jwt tokens and rule based access control.
- **Broker** - Ephemeral pubsub messaging for asynchronous communication and distributing notifications
- **Config** - Dynamic configuration and secrets management for service level config without the need to restart
- **Events** - Event streaming with ordered messaging, replay from offsets and persistent storage
- **Network** - Inter-service networking, isolation and routing plane for all internal request traffic
- **Proxy** - An identity aware proxy used for remote access and any external grpc request traffic
- **Runtime** - Service lifecyle and process management with support for source to running auto build
- **Registry** - Centralised service discovery and API endpoint explorer with feature rich metadata
- **Store** - Key-Value storage with TTL expiry and persistent crud to keep microservices stateless

**Framework**

Micro additionaly now contains the incredibly popular Go Micro framework built in for service development. 
The Go framework makes it drop dead simple to write your services without having to piece together lines and lines of boilerplate. Auto 
configured and initialised by default, just import and get started quickly.

**Command Line**

Micro brings not only a rich architectural model but a command line experience tailored for that need. The command line interface includes 
dynamic command mapping for all services running on the platform. Turns any service instantly into a CLI command along with flag parsing 
for inputs. Includes support for multiple environments and namespaces, automatic refreshing of auth credentials, creating and running 
services, status info and log streaming, plus much, much more.

**Environments**

Finally Micro bakes in the concept of `Environments` and multi-tenancy through `Namespaces`. Run your server locally for 
development and in the cloud for staging and production, seamlessly switch between them using the CLI commands `micro env set [environment]` 
and `micro user set [namespace]`.

## Development

Micro focuses on the microservices development model, which takes from the unix philosophy of writing tools that do one thing well. 
We think the domain boundary you've come to know at the database table level in Rails monolithic web apps moves to service 
boundaries at the network and naming layer.

For example a blog app might consist of posts, comments, tags and so on as tables in a Rails app, but in Micro these 
would be defined as independent services. These as Micro services which do one thing well and communicate over the 
network or via pubsub messaging where necessary.

Micro is built with this Service development model in mind which is why the underlying platform defines the primitives required 
to write those along with accessing them from external means. Micro includes a Go service framework that makes it super 
simple to get started fast.

## Cloud Services

We think the definition of a Cloud Service is one that helps you build the next Twilio or Stripe. Cloud services are looking 
more and more like something built to be consumed entirely as an API. So Micro builds with that model in mind. You write 
microservices on the backend and stitch them together as a single API for the frontend.

Micro provides an API gateway that handles HTTP/JSON requests externally and converts them to gRPC for the backend. This 
massively simplifies the experience of building efficient highly performant services on the backend which are decoupled 
from each other but presenting a single view to the consumers.

## Remote Access

Micro was built with the knowledge that not only do we exist in a multi-environment model but one that's remote first. Because 
of that we build in a gRPC identity proxy for CLI and local services that enables you to remotely connect to any Micro server 
securely and access those services and resources with your credentials stored in the Auth service.

You can assume not only are your services built for a Cloud first era but that your access to them is in that manner also.

## Sample Services

Check out the [micro/services](https://github.com/micro/services) open source repository for example services like the blog.

## Getting Started

Head to the  [getting started](/getting-started) guide to start writing Micro services now.

