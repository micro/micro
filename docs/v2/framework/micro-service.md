---
title: Micro Service
keywords: micro, service
tags: [micro, service]
sidebar: home_sidebar
permalink: /micro-service
summary: Turn anything into a micro service
---

Turn anything into a micro service. Micro provides a way of encapsulating anything to become a service.

## Overview

Micro is a runtime which manages microservices. The command line `micro service` encapsulates any app or service 
making it accessible within the micro ecosystem. The below example is for a basic http app.

## HTTP App

Here's a simple http hello world app

```
package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})
	http.ListenAndServe(":9090", nil)
}
```


Start the service using micro

```
micro service --name helloworld --endpoint http://localhost:9090 go run main.go
```

Query the service via the cli

```
micro call -o raw helloworld /
```

## File Server

Serve a file back to the caller

The file /tmp/helloworld.txt

```
helloworld
```

Run the service

```
micro service --name helloworld --endpoint file:///tmp/helloworld.txt
```

Get the file

```
micro call -o raw helloworld .
```

## Exec script

Execute a script or command remotely

```
#!bin/bash

echo `date` hello world
```

```
micro service --name helloworld --endpoint exec:///tmp/hellworld.sh
```

```
micro call -o raw helloworld .
```

