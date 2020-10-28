# Plugins

Micro operates as an abstraction layer for microservices and distributed systems development.

## Overview 

We make use of Go's interface system to create a runtime agnostic abstraction 
implemented for a variety of underlying systems. Our goal is ultimately to simplify the development 
experience and provide strongly defined APIs that are understood by the developer.

## Implementation

Our philosophy for plugins and the ecosystem in general is as follows

- Define a go-micro interface. This becomes the building block e.g Registry
- Define a zero dep default implementation e.g memory or mdns
- Define a highly available external implementation e.g etcd
- Define a "service" implementation using the runtime and implement service e.g micro registry

## Feature Matrix

Here's our list of implementations or preferred list

Interface | Default | Highly-Available | Service
--------- | ------- | ---------------- | -------
Auth | None | X (casbin?) | micro auth
Broker | E-NATS | NATS | micro broker
Config | memory | X (github?) | micro config
Registry | mDNS | etcd | micro registry
Store | memory | cockroachdb | micro store

TODO: complete table


## Roadmap

In v2 we plan to streamline how plugins are loaded and used. Rather than specifying a number of core defaults 
all plugins will move to github.com/micro/go-plugins and we'll look to provide a better developer experience 
for loading these at build or runtime.

## Design Ideas

Our future goal will be to generate a `plugins.go` file in top level main package. This will be generated 
via the use of the `--plugins` flag when calling `micro {new, build, run}`. Each command serves a 
purpose of generating a template, building the service or running it. They'll check for 
the plugins file and regen as necessary.

```
micro run service --plugins=broker/rabbitmq ./...
```

Additionally what we'd like to do is generate teh shared objects using GitHub actions against predefined micro/go-micro go.mod deps 
so we can load modules on the fly.
