
Micro is an open source platform for cloud native application development.

## Overview

Micro addresses the key requirements for building services in the cloud. It leverages the microservices
architecture pattern and provides a set of services which act as the building blocks of a platform. Micro deals
with the complexity of distributed systems and provides simpler programmable abstractions to build on. 

Micro provides a logical server composed of building block services, a Go framework for development, command 
line interface, API gateway and gRPC Proxy for external and remote access. Each service provides access 
to underlying infrastructure primitives through a standard interface with a development model tying everything 
together.

<center>
  <img src="/images/micro-3.0.png" />
</center>

## Features

Micro is built as a microservices architecture and abstracts away the complexity of the underlying infrastructure. We compose 
this as a single logical server to the user but decompose that into the various building block primitives that can be plugged 
into any underlying system. 

**Server**

The server is composed of the following services.

- **API** - HTTP Gateway which dynamically maps http/json requests to RPC using path based resolution
- **Auth** - Authentication and authorization out of the box using jwt tokens and rule based access control.
- **Broker** - Ephemeral pubsub messaging for asynchronous communication and distributing notifications
- **Config** - Dynamic configuration and secrets management for service level config without the need to restart
- **Events** - Event streaming with ordered messaging, replay from offsets and persistent storage
- **Network** - Inter-service networking, isolation and routing plane for all internal request traffic
- **Proxy** - gRPC identity aware proxy used for remote access and any external grpc request traffic
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

## Source

[GitHub Repo](https://github.com/micro/micro)

## Download

[Latest Release](https://github.com/micro/micro/releases/latest)

## Content

Documentation, guides and quick starts for Micro

- [Introduction](introduction) - A high level introduction to Micro
- [Getting Started](getting-started) - The helloworld quickstart guide
- [Upgrade Guide](upgrade-guide) - Update your go-micro project to use micro v3.
- [Architecture](architecture) - Describes the architecture, design and tradeoffs
- [Reference](reference) - In-depth reference for Micro CLI and services
- [Resources](resources) - External resources and contributions
- [Roadmap](roadmap) - Stuff on our agenda over the long haul
- [Users](users) - Developers and companies using Micro in production
- [FAQ](faq) - Frequently asked questions
- [Blog](blog) - For the latest from us

## Contributing

See the [TODO](/todo) list, open a PR and start hacking away at the docs.

## Community

Join us on [Discord](https://discord.gg/hbmJEct) or [Slack](https://slack.micro.mu). Follow [@microhq](https://twitter.com/microhq) on Twitter for updates.

## Hosting

If you're interested in a hosted version of Micro aka the Micro Platform see [m3o.com](https://m3o.com).

## License

[Polyform Shield](https://polyformproject.org/licenses/shield/1.0.0/)
