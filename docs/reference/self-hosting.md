---
title: Self hosting
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference/self-hosting
---

# Self-hosting

## Basic setup

Self hosting Micro backed only by a file store for persistence and using `mdns` for service discovery is very easy.
In fact it happens automatically every time one runs `micro server` locally.

By default the `/tmp/micro` folder will store all files (@TODO describe how to change this) the `micro server` is using.
Once the server is running and its 8081 port is exposed to the outside world, it can be connected to with the Micro CLI by [setting up the environment](environments).

## Using etcd for service discovery

TODO

## Using kubernetes for runtime

TODO