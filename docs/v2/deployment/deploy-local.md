---
title: Local Deployment
keywords: local
tags: [local]
sidebar: home_sidebar
permalink: /deploy-local
summary: 
---

Micro is incredibly simple to spin up locally

## Install

From source

```
go get github.com/micro/micro/v2
```

Release binary

```
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

## Run

Running micro is as simple as typing `micro`.

```
# Display the commands
micro --help
```

To run the server

```
micro server
```

