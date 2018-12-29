# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)

Micro is a toolkit for microservice development.

# Overview

Micro addresses the key requirements for building scalable systems. It takes the microservice architecture pattern and transforms it into 
a set of tools which act as the building blocks of a platform. Micro deals with the complexity of distributed systems and provides 
simple abstractions already understood by developers.

<img src="https://micro.mu/micro-diag.png" />

Technology is constantly evolving. The infrastructure stack is always changing. Micro is a pluggable toolkit which addresses these issues. 
Plug in any stack or underlying technology. Build future-proof systems using micro.

## Features

The toolkit is composed of the following features:

- **API Gateway:** A single entry point with dynamic request routing using service discovery. The API gateway allows you to build a scalable 
microservice architecture on the backend and consolidate serving a public api on the frontend. The micro api provides powerful routing 
via discovery and pluggable handlers to serve http, grpc, websockets, publish events and more.

- **Interactive CLI:** A CLI to describe, query and interact directly with your platform and services from the terminal. The CLI 
gives you all the commands you expect to understand what's happening with your micro services. It also includes an interactive mode.

- **Web Dashboard:** The web dashboard allows you to explore your services, describe their endpoints, the request and response formats and even 
query them directly. The dashboard is also includes a built in CLI like experience for developers who want to drop into the terminal on the fly.

- **SlackOps Bot:** A bot which runs on your platform and lets you manage your applications from Slack itself. The micro bot enables ChatOps 
and gives you the ability to do everything with your team via messaging. It also includes ability to create slack commmands as services which 
are discovered dynamically.

- **Service Templates:** Generate new service templates to get started quickly. Micro provides predefined templates for writing micro services. 
Always start in the same way, build identical services to be more productive.


## Getting Started

See the [docs](https://micro.mu/docs) for detailed information on the architecture, installation and use of the toolkit.

## Sponsors

Sixt is an Enterprise Sponsor of Micro

<a href="https://micro.mu/blog/2016/04/25/announcing-sixt-sponsorship.html"><img src="https://micro.mu/sixt_logo.png" width=150px height="auto" /></a>

Become a sponsor by backing micro on [Patreon](https://www.patreon.com/microhq). Finding existing backers in [BACKERS.md](BACKERS.md).
