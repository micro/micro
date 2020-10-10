---
title: Environments
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference/environments
---

## Environments

Micro is built with a federated and multi-environment model in mind. Our development normally maps through local, staging and production, so Micro takes 
this forward looking view and builds in the notion of environments which are completely isolated micro environments you can interact with through the CLI. 
This reference explains environments.

## Usage

Environments can be displayed using the `micro env` command.

```sh
$ micro env
* local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
```

There are three builtin environments, `local` being the default, and two [`m3o` specific](m3o.com) offerings; dev and platform.
These exist for convenience and speed of development. Additional environments can be created using `micro env add [name] [host:port]`. 
Environment addresses point to the micro proxy which defaults to :8081.

### Add Environment

The command `micro env --help` provides a summary of usage. Here's an example of how to add an environment.

```sh
$ micro env add myown stunningproject.com
$ micro env
* local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
  foobar    example.com
```

### Set Environment

The `*` marks wich environment is selected. Let's select the newly added:

```sh
$ micro env set myown
$ micro env
$ micro env
  local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
* foobar     example.com
```

### Login to an Environment

Each environment is effectively an isolated deployment with its own authentication, storage, etc. So each env requires signup and login. 
At this point we have to log in to the `example` env with `micro login`. If you don't have the credentials to the environment, you have to ask the admin.

If your `example` environment does not exist yet, you might find the [self hosting](/reference/self-hosting) reference guide useful.
