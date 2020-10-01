---
title: Docker Deployment
keywords: docker
tags: [docker]
sidebar: home_sidebar
permalink: /deploy-docker
summary: 
---

Micro easily runs inside docker containers

### Install Micro

```
docker pull micro/micro
```

## Compose

Run a local server using docker compose

```
server:
  command: server
  build: .
  ports:
  - "8080:8080"
  - "8081:8081"
  - "8082:8082"
```

## Build from scratch

A Dockerfile is included in the repo

```
## checkout the repo
git clone https://github.com/micro/micro

## build the image
cd micro && docker build -t micro .
```

