---
title: Environments
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference/environments
---

# Environments

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

There are three builtin environments, `local` being the default one, and two [`m3o` specific](m3o.com) ones; dev and platform.
These exist for convenience and speed of development.

The beauty of the Micro envs are however that users can add their own. This is extremely useful when one wants to interact with one's own self hosted Micro server instance.

`micro env --help` provides a succint summary of usage, but let's walk through adding an environment:

```sh
$ micro env add myown stunningproject.com
$ micro env
* local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
  myown      stunningproject.com
```

The `*` marks wich environment is selected. Let's select the newly added:

```sh
$ micro env set myown
$ micro env
$ micro env
  local      127.0.0.1:8081
  dev        proxy.m3o.dev
  platform   proxy.m3o.com
* myown      stunningproject.com
```

At this point we have to log in to the `myown` env with `micro login`.
If you don't have the credentials to the environment, you have to ask the admin.

If your `myown` environment does not exist yet, you might find the [self hosting](self-hosting) reference guide useful.
