---
layout:	post
title:	Micro - The great consolidation of 2019
date:	2019-06-10 09:00:00
---
<br>
Micro started it's journey as [**go-micro**](https://github.com/micro/go-micro) - a microservices framework - focused 
on providing the core requirements for microservice development. It creates a simpler experience for building microservices 
by abstracting away the complexity of distributed systems.

Over time we've expanded beyond go-micro into additional tools, libraries and plugins. This has led to fragmentation 
in the way we solve problems and how developers are expected to use micro tools. We're now undergoing a consolidation 
of all these tools to simplify the developer experience.

Micro essentially moves to being a standalone development framework and runtime for microservices development.

Before we talk about the consolidation, let's revisit the journey thus far.

## The Primary Focus

Go-micro started out as primarily being focused on the communication aspects of microservices and we've tried to remain true to that. 
This opinionated approach and focus is what's really driven the success of the framework thus far. Over the years 
we've had numerous requests to solve the next day problems for building production ready software in go-micro itself. 
Much of this has been related to scalability, security, synchronisation and configuration.

<center>
<img src="https://m3o.com/docs/images/go-micro.svg" style="width: 80%; height: auto;"/>
</center>

While there would have been merit to add the additional features requested we really wanted to stay very focused 
on solving one problem really well at first. So we took a different approach to promoting the community to do this.

## Ecosystem & Plugins

Going to production involves much more than vanilla service discovery, message encoding and request-response. 
We really understood this but wanted to enable users to choose the wider platform requirements via pluggable and 
extensible interfaces. Promoting an ecosystem through an [explorer](https://m3o.com/explore/) which aggregates 
micro based open source projects on GitHub along with extensible plugins via the [go-plugins](https://github.com/micro/go-plugins) repo.

<center>
<img src="https://m3o.com/explorer.png" style="width: 80%; height: auto;" />
</center>

Go Plugins has generally been a great success by allowing developers to offload significant complexity to systems built 
for those requirements e.g prometheus for metrics, zipkin for distributed tracing and kafka for durable messaging.

## Points of Interaction

Go Micro was really at the heart of microservice development but as services were written the next questions moved on to; 
how do I query these, how do I interact with them, how do I serve them by traditional means...

Given go-micro used an rpc/protobuf based protocol that was both pluggable and runtime agnostic, we needed some way to address 
this in a way that was true to go-micro itself. This led to the creation of [**micro**](https://github.com/micro/micro), 
the microservice toolkit. Micro provides an api gateway, web dashboard, cli, slack bot, service proxy and much more.

<center>
<img src="https://m3o.com/assets/images/runtime3.svg" style="width: 80%; height: auto;"/>
</center>

The micro toolkit acted as interaction points via http api, browser, slack commands and command line interface. These 
are common ways in which we query and build on applications and it was important for us to provide a runtime that 
really enabled this. Yet still it really focused on communication above all else.
 
## Additional Tools

While the plugins and the toolkit helped users of micro significantly, it still lacked in key areas. It was clear that our community 
wanted us to solve further problems around platform tooling for product development rather than having to do it individually at 
their various companies. We needed the same types of abstractions go-micro provided for things like dynamic configuration, 
distributed synchronization and broader solutions for systems like Kubernetes.

For these we created:

- [micro/go-config](https://github.com/micro/go-config) - a dynamic configuration library
- [micro/go-sync](https://github.com/asim/go-sync) - a distributed synchronisation library
- [micro/kubernetes](https://github.com/micro/kubernetes) - micro on kubernetes initialisation
- [examples](https://github.com/micro/examples) - example usage code for micro
- [microhq](https://github.com/microhq) - a place for prebuilt microservices

These were a few of the repos, libraries and tools used to attempt to solve the wider requirements of our community. Over the 
4 years the number of repos grew and the getting started experience for new users became much more difficult. The barrier to 
entry increased dramatically and with that we knew something needed to change.

## Consolidation

In the past few weeks we've realised [**go-micro**](https://github.com/micro/go-micro) was really the focal point for most of our users 
developing microservices. It's become clear they want additional functionality as part of this library and as a self described 
framework we really need to embrace this by solving those next day concerns without asking a developer to go anywhere else. 

Essentially **go-micro** will be the all encompassing and standalone framework for microservice development.

We've started the consolidation process by moving all our libraries into go-micro and we'll continue to refactor over the 
coming weeks to provide a simpler default getting started experience while also adding further features for logging, tracing, metrics, 
authentication, etc.

<center>
<img src="{{ site.baseurl }}/blog/images/go-micro-repo.png" style="width: 80%; height: auto;" />
</center>

We're not forgetting about [**micro**](https://github.com/micro/micro) either though. In our minds after you've built your microservices 
there's still a need for a way to query, run and manage them. **Micro** is by all accounts going to be the runtime for micro 
service development. We're working on providing a simpler way to manage the end to end flow of microservice development and 
should have more to announce soon.

## Summary

Micro is the simplest way to build microservices and slowly becoming the defacto standard for Go based microservice development in 
the cloud. We're making that process even simpler by consolidating our efforts into a single development framework and runtime. 

<center>...</center>
<br>
To learn more check out the [website](https://m3o.com), follow us on [twitter](https://twitter.com/m3ocloud) or 
join the [slack](https://slack.m3o.com) community.

<h6><a href="https://github.com/micro/micro"><i class="fab fa-github fa-2x"></i> Micro</a></h6>
