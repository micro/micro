---
title: Install Locally
keywords: install
tags: [install, local]
sidebar: home_sidebar
permalink: /reference/install
summary: 
---

## Local Install

Micro can be installed locally in the following way. We assume for the most part a Linux env with Go and Git installed.

### Go Get

```
go get github.com/micro/micro/v3
```

### Docker

```
docker pull micro/micro
```

### Release Binaries

```
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

## Server Usage

To start the server simply run

```
micro server
```

This will boot the entire system and services including a http api on :8080 and grpc proxy on :8081

## Verify Install

Check help text is output with no errors
```
micro --help
```

Run helloworld

```
micro env	# should point to local
micro status	# returns empty response
micro services	# returns empty response
micro run github.com/micro/services/helloworld # run helloworld
micro status 	# wait for status running
```

Call the service and verify output

```shell
$ micro helloworld --name=John
{
        "msg": "Hello John"
}
```

Remove the service

```
micro kill helloworld
```
