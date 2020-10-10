---
title: CLI usage
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference/cli
summary: A CLI usage guide
---

Micro is driven entirely through a CLI experience. This reference highlights the CLI design.

## Overview

The CLI speaks to the `micro server` through the gRPC running locally by default on :8081. All requests are proxied based on your environment 
configuration. The CLI provides the sole interaction for controlling services and environments.

# Dynamic commands

When issuing a command to the Micro CLI (ie. `micro command`), if the command is not a builtin, Micro will try to dynamically resolve this command and call
a service running. Let's take the `micro registry` command, because although t he registry is core service that's running by default on a local Micro setup,
the `registry` command is not a builtin one.

With the `--help` flag, we can get information about available subcommands and flags

```sh
$ micro registry --help
NAME:
	micro registry

VERSION:
	latest

USAGE:
	micro registry [command]

COMMANDS:
	deregister
	getService
	listServices
	register
	watch
```

The commands listed are endpoints of the `registry` service (see `micro services`).

To see the flags (which are essentially endpoint request parameters) for a subcommand:

```sh
$ micro registry getService --help
NAME:
	micro registry getService

USAGE:
	micro registry getService [flags]

FLAGS:
	--service string
	--options_ttl int64
	--options_domain string

```

At this point it is useful to have a look at the proto of the [registry service here](https://github.com/micro/micro/blob/master/proto/registry/registry.proto).

In particular, let's see the `GetService` endpoint definition to understand how request parameters map to flags:

```proto
message Options {
	int64 ttl = 1;
	string domain = 2;
}

message GetRequest {
	string service = 1;
	Options options = 2;
}
```

As the above definition tells us, the request of `GetService` has the field `service` at the top level, and fields `ttl` and `domain` in an options structure.
The dynamic CLI maps the underscored flagnames (ie. `options_domain`) to request fields, so the following request JSON:

```js
{
    "service": "serviceName",
    "options": {
        "domain": "domainExample"
    }
}
```

is equivalent to the following flags:

```sh
micro registry getService --service=serviceName --options_domain=domainExample
```
