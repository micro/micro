---
title: FAQ
keywords: faq
tags: [faq]
sidebar: home_sidebar
permalink: /docs/faq.html
summary:
---

The FAQ should provide quick answers to the most common questions.

## What is Micro?

Micro empowers developers to build, share and collaborate on distributed systems aka microservices.

- Micro is a [Framework](https://github.com/micro/go-micro)
- Micro is a [Runtime](https://github.com/micro/micro)
- Micro is a [Network](https://micro/docs/network.html)
- Micro is a [Community](https://slack.m3o.com/)

### Open Source

Micro is built as open source libraries and tools to help with microservice development.

- **Framework** - A Go framework for writing microservices; service discovery, rpc, pub/sub, etc.
- **Runtime** - A microservice runtime environment; API gateway, cli, slackbot, service proxy, etc.
- **Plugins** - Plugins for the framework and runtime including etcd, kubernetes, nats, grpc, etc.

Find all the tools at [github.com/micro](https://github.com/micro).

### Network

Micro provides an open global always-on cloud environment for sharing and running services. The [network](https://micro.mu/docs/network.html) 
can be joined and used by anyone. It's a great way to unlock public collaboration outside of any one organisation or team.

Explore the [network](https://micro.mu/docs/network.html).

### Community

There's a slack community with thousands of members. 

Invite yourself at [micro.mu/slack/](https://slack.m3o.com/).

## Where do I start?

Start with [go-micro](https://github.com/micro/go-micro). The readme provides a sample microservice.

Learn more by reading the [getting started](https://micro.mu/docs/writing-a-go-service.html) guide or checkout the [examples](https://github.com/micro/examples).

Use the [micro](https://github.com/micro/micro) toolkit to access microservices via the cli, web ui, slack or api gateway.

## Who's using Micro?

See the [users](https://micro.mu/docs/users.html) page with a list of companies using Micro (but note it may be out of date). 

Many more are also using it but not yet publicly listed. Feel free to add your company if you're using Micro.

## How do I use Micro?

It's quite simple.

1. Write services using [go-micro](https://github.com/micro/go-micro).
2. Access them via the [micro](https://github.com/micro/micro) toolkit.
3. Profit.

Checkout the full [greeter](https://github.com/micro/examples/tree/master/greeter) example.

## Zero Dependency Discovery

Multicast DNS is a built in registry for a zero conf service discovery.

You don't need to do anything! It's just built in and enabled by default.

## Can I use something instead of MDNS?

Yes! The registry for service discovery is completely pluggable. Etcd is also included as a default plugin.

```shell
MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=127.0.0.1:2379 myservice
```

## Where can I run Micro?

Micro is runtime agnostic. You can run it anywhere you like. On bare metal, on AWS, Google Cloud. On your favourite container orchestration system like Mesos or Kubernetes.

In fact there's demo config for Micro on Kubernetes. Check out the repo at [github.com/micro/kubernetes](https://github.com/micro/kubernetes)

## What's the difference between API, Web and SRV services?

<img src="images/arch.png" />

As part of the micro toolkit we attempt to define a set of design patterns for a scalable architecture by separating the concerns of the API, Web dashboards and backend services (SRV).

### API Services

API services are a separate tier of services for serving your http/json API. This might be a public facing api or leveraged by your clients. API services 
are served by the `micro api` and conform to an api gateway pattern. They are written as standard Go Micro services but use the namespace 
`go.micro.api` to logically separate themselves.

### Web Services

Web services are served by the micro web with the default namespace go.micro.web. We believe in web apps as first class citizens in the microservice world therefor building web dashboards as microservices. The micro web is a reverse proxy and will forward HTTP requests to the appropriate web apps based on path to service resolution. 

Learn more about it [here](https://github.com/micro/micro/tree/master/web)

### SRV services

SRV services are basically standard RPC services, the usual kind of service you would write. We usually call them RPC or backend services as they should mainly be part of the backend architecture and never be public facing. By default we use the namespace go.micro.srv for these but you should use your domain com.example.srv. 

## How performant is it?

Performance is not a current focus of Micro but always something we strive for. We let the wider Go ecosystem handle 
the heavy lifting and focus on developer productivity first. This means we're using gRPC beneath the covers and NATS 
for messaging.

Where the default plugins are not performant enough you can use [go-plugins](https://github.com/micro/go-plugins) to 
switch out to your tool of choice.

## Does Micro support gRPC?

Yes. In v2 micro makes use of gRPC by default.

## Micro vs Go-Kit

This question comes up a lot. What's the difference between micro and go-kit?

Go-kit describes itself as a standard library for microservices. Like Go, go-kit provides you with individual packages 
which can be used to construct your applications. Go-kit is great where you want complete control over how you define 
your services.

Go Micro is a pluggable RPC framework for microservices. It's an opinionated framework which attempts to simplify 
the communication aspects of distributed systems so you can focus on the business logic itself. Go-micro is great 
where you want to get up and running quickly while having something pluggable to switch out infrastructure without 
code changes.

Micro is a microservice toolkit. It's like a swiss army knife for microservices which builds on go-micro to provide 
traditional entry points like http api gateway, web ui, cli, slack bot, etc. Micro uses tooling to guide the logical 
separation of concerns in your architecture, pushing you to create an API layer of microservices for a public API and 
separately creating a WEB layer of microservices for web UIs.

Use go-kit where you want complete control. Use go-micro where you want an opinionated framework.

## Where Can I Learn More?

- Join the slack community - [micro.mu/slack](https://slack.m3o.com)
- Read the blog - [micro.mu/blog](https://micro.mu/blog)

{% include docs/links.html %}
