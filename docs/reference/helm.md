---
title: Install Micro using Helm
keywords: install
tags: [install, helm, kubernetes]
sidebar: home_sidebar
permalink: /reference/helm
summary: 
---

## Helm

Micro can be installed onto a Kubernetes cluster using helm. Micro will be deployed in full and leverage zero-dep implementations designed for Kubernetes. For example, micro store will internally leverage a file store on a persistant volume, meaning there are no infrastructure dependancies required.

### Dependencies

You will need to be connected to a Kubernetes cluster

### Install

Install micro with the following commands:

```shell
helm repo add micro https://ben-toogood.github.io/micro-helm
helm install micro micro/micro --set image.repo=localhost:5000/micro
```

### Uninstall

Uninstall micro with the following commands:

```shell
helm uninstall micro
helm repo remove micro
```