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
