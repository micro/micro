---
title: Local Development
keywords: runtime
tags: [micro, runtime, dev]
sidebar: home_sidebar
permalink: /local-development
summary: Run a service with the micro runtime
---

The micro runtime provides a way to manage the lifecycle of services without the complexity of orchestration systems. 

You can start the service with one command

```
micro run [service]
```

## Usage

To run and manage your service locally do the following.

```
# Start the server
micro server

# run the service
micro run --server examples/greeter/srv
```

