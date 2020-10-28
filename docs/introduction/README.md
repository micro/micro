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

Micro is something of a distributed operating system made up of many independent services that all act in coordination to 
provide a server to build, run and manage services including your service development and access externally.

Micro includes:

- Authentication
- Configuration
- PubSub Messaging
- Event Streaming
- Service Discovery
- Service Networking
- Key-Value Storage
- HTTP API Gateway
- gRPC Identity Proxy

And many more features. Micro is packed with a CLI interface and a Service Library used to write your applications.

## Development

Micro focuses on the microservices development model, which takes from the unix philosophy of writing tools that do one thing well. 
We think the domain boundary you've come to know at the database table level in Rails monolithic web apps moves to service 
boundaries at the network and naming layer.

For example a blog app might consist of posts, comments, tags and so on as tables in a Rails app, but in Micro these 
would be defined as independent services. These as Micro services which do one thing well and communicate over the 
network or via pubsub messaging where necessary.

Micro is built with this Service development model in mind which is why the underlying platform defines the primitives required 
to write those along with accessing them from external means. Micro includes a Go service library that makes it super 
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

