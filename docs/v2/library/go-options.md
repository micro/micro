---
title: Options
keywords: go-micro, framework
tags: [go-micro, framework]
sidebar: home_sidebar
permalink: /go-options
summary: Setting and using Go Micro options
---

Go Micro uses a variadic options model for the design of passing arguments for the creation and initialisation of 
packages and also as optional params for methods. This offers flexibility in power for extending our option usage 
across plugins.

## Overview

When create a new service you have option of passing additional parameters such as setting the name, version, 
the message broker, registry or store to use along with all the other internals. 

Options are normally defined as follows

```
type Options struct {
  Name string
  Version string
  Broker broker.Broker
  Registry registry.Registry
}

type Option func(*Options)

// set the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// set the broker
func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}
```

These can then be set as follows

```
service := micro.NewService(
	micro.Name("foobar"),
	micro.Broker(broker),
)
```

## Service Options

Within Go Micro we have a number of options that can be set including the underlying packages that will be used 
for things such as authentication, configuration and storage. You can use `service.Options()` to access these.

The packages such as auth, config, registry, store, etc will default to our zero dep plugins. Where you want 
to configure them via env vars or flags you can specify `service.Init()` to parse them.

For example, if you replace the memory store to use a file store it can be done as follows

```
## as an env var
MICRO_STORE=file go run main.go

## or as a flag
go run main.go --store=file
```

Internally the store is then accessible via the options

```
service := micro.NewService(
	micro.Name("foobar"),
)

service.Init()

store := service.Options().Store
```


