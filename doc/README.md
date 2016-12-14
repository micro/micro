# Docs

Micro is a microservice toolkit. It's purpose is to simplify distributed systems development. 

Checkout the [roadmap](roadmap.md) to see where it's all going.

View the [summary](SUMMARY.md) to navigate the docs.

# Overview

The goal of Micro is to provide a toolkit for microservice development and management. 
At the core, micro is simple and accessible enough that anyone can easily get started 
writing microservices. As you scale to hundreds of services, micro will provide the 
fundamental tools required to manage a microservice environment.

The toolkit is composed of the following components:

**Go Micro** - A pluggable RPC framework for writing microservices in Go. It provides libraries for 
service discovery, client side load balancing, encoding, synchronous and asynchronous communication.

**API** - An API Gateway that serves HTTP and routes requests to appropriate micro services. 
It acts as a single entry point and can either be used as a reverse proxy or translate HTTP requests to RPC.

**Web** - A web dashboard and reverse proxy for micro web applications. We believe that 
web apps should be built as microservices and therefore treated as a first class citizen in a microservice world. It behaves much the like the API 
reverse proxy but also includes support for web sockets.

**Sidecar** - The Sidecar provides all the features of go-micro as a HTTP service. While we love Go and 
believe it's a great language to build microservices, you may also want to use other languages, so the Sidecar provides a way to integrate 
your other apps into the Micro world.

**CLI** - A straight forward command line interface to interact with your micro services. 
It also allows you to leverage the Sidecar as a proxy where you may not want to directly connect to the service registry.

**Bot** A Hubot style bot that sits inside your microservices platform and can be interacted with via Slack, HipChat, XMPP, etc. 
It provides the features of the CLI via messaging. Additional commands can be added to automate common ops tasks.

<p align="center">
  <img src="images/overview.png" />
</p>
