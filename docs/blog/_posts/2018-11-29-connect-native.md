---
layout:	post
title:	Consul Connect-Native Go Micro Services
date:	2018-11-29 09:00:00
---
<br>
Today we're announcing support for Connect-Native Go Micro services via a slim initialisation library called [Go Proxy](https://github.com/micro/go-proxy). 
This will provide [Go Micro](https://github.com/micro/go-micro) with the ability to do authorized and secure service-to-service communication.

## What is Consul Connect?

[Consul Connect](https://www.consul.io/docs/connect/index.html) is a feature of [Consul](https://www.consul.io/) which provides 
service-to-service authorization and encryption via mutual TLS. Consul Connect uses [SPIFEE](https://spiffe.io/) compliant 
certificates for identity.

We believe Consul Connect is a powerful mechanism for securing micro services. So how do we integrate?

## Connect-Native

Consul [Connect-Native](https://www.consul.io/docs/connect/native.html) is native integration with the Connect API. This allows 
[Go Micro](https://github.com/micro/go-micro) services to become secure by default. 

Consul Connect provides the ability to use proxies for communication but this can add overhead, Go Micro handles distributed 
systems concerns as a client library, which eliminates this overhead. Native integration with Connect gives us all its benefits 
while maintaining direct point to point connections for performance.

The consul documentation provides an overview of how this works. In Go Micro's case we initialise a consul registry with the 
connect option enabled and setup the broker and transport tls config.

<img src="https://www.consul.io/assets/images/connect-native-overview-cc9dc497.png" />

## Using Connect-Native

We've provided a complete example of how to get started in the [Go Proxy](https://github.com/micro/go-proxy) repository.

But essentially it's a two line change. Import the connect package and create a new service with it. That's it!

<script src="https://gist.github.com/asim/de7a3bcfcd93f6102e6c657ed54b8f2e.js"></script>

## Summary

Connect-Native gives us support for authorization and secure end to end communication for Go Micro apps without the overhead 
of proxies. It's a great addition for micro users and we highly recommend using it.

<center>...</center>

To learn more about micro check out the [website](https://m3o.com), follow us on [twitter](https://twitter.com/m3ocloud) or 
join the [slack](https://slack.m3o.com) community.

<h6><a href="https://github.com/micro/go-proxy"><i class="fab fa-github fa-2x"></i> Go Proxy</a></h6>
