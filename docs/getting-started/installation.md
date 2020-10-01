---
title: Install Guide
keywords: install
tags: [install]
sidebar: home_sidebar
permalink: /installation
summary: 
---

## Framework

Micro is a framework for cloud native development

### Dependencies

You will need protoc-gen-micro for code generation

- [protobuf](https://github.com/golang/protobuf)
- [protoc-gen-go](https://github.com/golang/protobuf/tree/master/protoc-gen-go)
- [protoc-gen-micro](https://github.com/micro/micro/tree/master/cmd/protoc-gen-micro)

```
# Download latest proto releaes
# https://github.com/protocolbuffers/protobuf/releases
go get github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/micro/v3/cmd/protoc-gen-micro
```

### Install

From source

```
go get github.com/micro/micro/v3
```

Docker image

```
docker pull micro/micro
```

Latest release binaries

```
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

### Usage

Start the server

```shell
micro server
```

Set your env to local

```shell
micro env set local
```

Run an example helloworld service

```shell
micro run github.com/micro/services/helloworld
```

List services

```shell
micro services
```

Get status

```shell
micro status
```

Call service

```shell
micro helloworld --name="John"
```

Output

```shell
{
	"msg": "Hello John"
}
```

