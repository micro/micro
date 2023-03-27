---
layout: post
title:  "Platform APIs"
author: Asim Aslam
date:   2023-03-20 10:00:00
---

Platform APIs are a set of building blocks for platform teams to deliver self-service infrastructure to engineers within their orgs. They are a way of delivering the required dependencies for applications without the need for direct interaction with cloud infrastructure, using APIs, SDKs, CLI and Web UI.

As platform engineers begin to deliver cloud infrastructure to dev teams in a self-serve manner they’re attempting to understand the conceptual pieces to provide and where the abstraction layer lies. Providing a programmable set of platform APIs could become the de facto standard method.

Prior attempts have been made to provide such platforms with [Netflix OSS](https://netflix.github.io/) and Spring Cloud, Azure Service Fabric and [Go Micro](https://go-micro.dev). 
Today though, [Dapr](https://dapr.io) and [Micro](https://micro.dev) have emerged as the open source successors, creating a potential new category and market for API platforms.

## Defining the term “Platform APIs”
Platform APIs are in reference to a common set of building blocks used by most applications. This involves some form of data storage, pubsub messaging, config and secrets management, service to service invocation and authentication. The list is non exhaustive but these are the most common.

Dapr is a “distributed application runtime” which groups these components together into a single application binary consumed as an API via HTTP, gRPC, command line or multi-language SDKs.

Netflix pioneered the idea of an internal platform with Netflix OSS and a stack of components for Java that required client side integration into the applications themselves. This was repackaged as Spring Cloud and became a cloud hosted provider managed by Pivotal later acquired by VMWare.

More client side implementations have appeared such as Go Micro, a Go microservices framework, which attempts to provide the same abstractions for modern day programmers in the programming language Go. Despite its popularity, Go Micro suffered the same problems as past implementations. 

Client side frameworks require developers not only to rewrite applications to use the framework but also require them to reason about distributed systems principles. It was not the abstraction people needed or wanted, especially as app developers. Service mesh tried to take this on as a language agnostic proxy but failed to formulate a strategy beyond service discovery and load balancing.

Based on learnings from Azure Service Fabric, Dapr has managed to codify the platform needs as a single set of APIs which can be consumed largely without a client framework or direct knowledge of cloud infrastructure. It has in principle defined what the abstraction layer for the Cloud should be. 

Micro predates Dapr and has long since been focused on the same primitives as it learned from prior art in Go Micro and the industry at large. Ultimately it doesn't 
matter which came first, but moreso that we're converging on a similar set of building blocks for the next layer on top of container orchestration platforms, service mesh and CNCF technologies.

## APIs Primitives

Both Micro and Dapr define a core set of primitives required for all application development. These are provided by most cloud providers or through open source infrastructure that is self hosted. In the majority of cases the underlying dependencies exist across cloud providers.

Here’s a look at the separation between user application code and the platform runtimes. You can see that Dapr and Micro create an abstraction layer on top of cloud providers and define the core APIs required by all applications. This effectively makes them portable across clouds.

<center>
  <b>Dapr APIs</b>
</center>

<img src="https://github.com/dapr/dapr/raw/master/img/overview.png" />

<center>
  <b>Micro APIs</b>
</center>

<img src="https://github.com/micro/micro/raw/master/docs/images/micro.png" />

Both include the following set of building blocks built as gRPC APIs accessible via API, CLI, SDKs and Web UI.

### Service to Service Invocation
Dapr and Micro both provide service to service invocation. This codifies the service oriented architecture or microservices patterns that are prominent in the cloud today. Everything is now being decoupled on the backend into separate services run on kubernetes or elsewhere, so these services need a way to call each other.

### State Management
Included is a basic key-value storage engine which can be backed by etcd, redis, postgres or other solutions. Given that service based architectures are supposed to have stateless apps, this enables developers to offload that state to some sort of database or key-value storage.

### PubSub Messaging
Event driven architectures are at the heart of most platforms now. So pubsub is built in as a first class citizen, again backable by redis, kafka or any of the cloud based services like SQS, etc. This includes things like fire-and-forget and event streaming for persistence and queuing.

### Config and Secrets
Config and secrets end up being driven by environment variables everywhere, this often isn’t great for things that need to change dynamically or sensitive data like API keys and passwords. Both Dapr and Micro bake in a Config API for this. It can be managed via the CLI.

### Authentication & Authorization
Micro includes service to service authentication via JWT keys, role based access control and scoping requests by endpoint through an API gateway. This is something that appears to be lacking from Dapr but is pretty core to all platform concerns. As teams scale, each team has to gatekeep what’s publicly available vs internal. The same for external requests by end users.

### Service Runtime 
Another missing piece from Dapr is a service runtime, it relies on Kubernetes to handle this management whereas Micro bakes this in at the heart of its build, run, manage and consume philosophy. Micro abstracts away the runtime, whether it’s local or remote using environments, and unifies different systems under a single API. Local as Go builds and remote as kubernetes.

### Pluggable Architecture
The other key component to platform APIs is the pluggable architecture. They are essentially an abstraction layer for cloud infrastructure, in different companies and teams, architectural choices will be different, so plugins and components are made available for this reason.

## Components, Modules and Plugins

How do platform APIs handle plugins, modules or components? Platform APIs are strongly defined API interfaces. In the case of Go Micro they were Go interfaces, for Dapr and Micro, they are defined as protobuf APIs. Given the need for cloud primitives to become building blocks, it requires a model for distribution, adoption and extension. The idea of plugins, modules or components is not new, but the way in which they’re defined has evolved over the years.

Wasm components are slowly emerging as a new way to define building blocks and API primitives.

Here's a great explainer from [Fermyon](https://www.fermyon.com/blog/webassembly-component-model) on the subject. 

> The WebAssembly Component Model is a proposal to build upon the core WebAssembly standard by defining how modules may be composed within an application or library. The repository where it is being developed includes an informal specification and related documents, but that level of detail can be a bit overwhelming at first glance.

Given Wasm’s goal of being a language agnostic runtime that enables cross language reuse of code, this feels like a perfect fit for platform APIs.
If we take something as simple as a key-value storage API, with Wasm components, this can easily be defined as a single reusable interface, implemented in any language, and imported to be used within the code directly in an application or in the platform itself. In the case of Dapr or Micro this could be an internal abstraction layer for how to plugin various runtime components.

At the moment plugins within Dapr and Micro are Go libraries, defined as interfaces, which must then be imported into the server and require a rebuild of the binary. This is very cumbersome. While the end user interfaces with a gRPC or HTTP API, the internal system must create its own plugin/abstraction layer for cloud infrastructure to easily swap out the components.

There is prior work in using RPC/gRPC but this can introduce an extra hop which adds latency. If Wasm components are used, it creates a new portable, reusable and cross architecture compatible model to build native plugins as a set of platform APIs that everyone can use as a standard for application and platform development.

## Ecosystem Potential

The ecosystem potential of platform APIs is pretty phenomenal. In some ways it's a new paradigm and interface to build on which we have not seen for a decade or more. With so many tools now 
in the cloud infrastructure stack, it's key that we're looking for a way to tame this complexity. Platform APIs offer the ability to level up and give back the power to developers while also 
enabling operators to manage the infrastructure in a way that they choose. This potentially brings us one step close to true application and platform portability.

Let's look at the opportunities from both perspectives.

### Developer Opportunity
[Docker](https://docker.com) as a developer centric platform has mostly remained a container runtime and workflow which enables delivering apps or services in a common format, but the internals of the container is largely a black box. It enables opening ports, some networking and storage, but otherwise doesn't care much about what is in the container. 

[Compose](https://github.com/docker/compose) took Docker a step further by helping the developer reason about the application architecture through a multi-service deployment specification e.g compose yaml, and while it helps to define some interdependencies between the applications, it avoids dealing with the development aspect or with the specifics of the infrastructure dependencies themselves.

In most cases, developers are using Redis, Postgres, etc as dependencies for their apps. They’re also defining keys, passwords, and related information as environment variables. Significant amounts of these concerns are related to platform APIs and the deployment infrastructure itself. Docker could potentially define core primitives as part of the compose experience.

How would that work in practice? Let’s assume the user is defining some configuration via environment variables, this can be elevated to a first class citizen type such as Config or Secret in the case of security sensitive values like API credentials. The delivery mechanism can then be abstracted away by a platform API or Wasm component that defines how it is loaded.

Kubernetes has defined high level types through APIs like this, but ultimately the system is very complex and is not portable across different architectures e.g you must be a user of Kubernetes to gain this benefit. Whereas the user of Dapr or Micro can deploy to any architecture and take their config, secrets, etc with them. 

Another example is key-value storage or caching. Redis is predominantly used as a Cache. Docker could define another type known as a Cache with the API Get/Set/Delete implemented as a platform API component and loaded at runtime for the user. The cache implementation could be swapped out without affecting the user’s application in the case of mocking or local platforms.

The same principle applies to all platform APIs. Grouped together this can become a powerful building block for Docker apps known as "Docker Blocks"... Or whatever name makes most sense.
This is Docker's potential opportunity from a platform API perspective.

### Operator Opportunities
The likes of Hashicorp are very much focused on the DevOps ecosystem and they’ve built the building blocks for platforms, but have yet to create a cohesive strategy across their product portfolio.
Platform APIs could provide the glue layer for their suite of services along with many other companies, creating a new layer of abstraction that targets the developer as the consumer, allowing 
the operator to deliver self-serve infrastructure via a centralised platform engineering team.

The dividing line in the developer workflow is the platform API. Hashicorp, Docker and others in the ecosystem have the opportunity to create a new way to bridge the gap between 
both developers and infrastructure using platform APIs as the glue layer. It's not clear whether either will take this approach but Dapr and Micro are already advancing in the field.

Ultimately the driver of success and the goal of a business is to deliver applications faster and more efficiently to end customers. The more productive developers are, the more this translates to direct business value. Part of this success is offering developers self-serve platform APIs that enable them to move rapidly, in leveraging the relevant technologies across varying platforms.

## Closing Thoughts

This is just a quick summary of platform APIs and the relevance to the ecosystem today. We're seeing the emergence of platform engineering as a trend and the need for centralised platforms 
within companies that deliver self serve infrastructure via platform APIs in a programmatic manner, which can be consumed via API, CLI, SDK or Web UI. 

It will be exciting to see how this trend plays out in the coming year or two.
