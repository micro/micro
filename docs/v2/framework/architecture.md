---
title: Architecture
keywords: architecture
tags: [architecture]
sidebar: home_sidebar
summary: 
---

Micro provides the fundamental building blocks for microservices. It's goal is to simplify distributed systems development. Because microservices is an architecture pattern, Micro looks to logically separate responsibility through tooling. 

Check out the blog post on architecture [https://micro.mu/blog/2016/04/18/micro-architecture.html](https://micro.mu/blog/2016/04/18/micro-architecture.html) for a detailed 
overview.

This section should explain more about how micro is constructed and how the various libraries/repos relate to each other.

## Runtime

### API

The API acts as a gateway or proxy to enable a single entry point for accessing micro services. It should be run on the edge of your infrastructure. It converts HTTP requests to RPC and forwards to the appropriate service.

<p align="center">
  <img src="images/api.png" />
</p>

### Web

The UI is a web version of go-micro allowing visual interaction into an environmet. In the future it will be a way of aggregating micro web services also. It includes a way to proxy to web apps. /[name] will route to a a service in the registry. The Web UI adds Prefix of "go.micro.web." (which can be configured) to the name, looks 
it up in the registry and then reverse proxies to it.

<p align="center">
  <img src="images/web.png" />
</p>

### Proxy

The proxy is a cli proxy for remote environments.

<p align="center">
  <img src="images/car.png" />
</p>

### Bot

Bot A Hubot style bot that sits inside your microservices platform and can be interacted with via Slack, HipChat, XMPP, etc. It provides the features of the CLI via messaging. Additional commands can be added to automate common ops tasks.

<p align="center">
  <img src="images/bot.png" />
</p>

### CLI

The Micro CLI is a command line version of go-micro which provides a way of observing and interacting with a running environment.

## Plugins

Plugins are a way of adding additional functionality to the runtime. See the [overview](runtime-plugins.html).

