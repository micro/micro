# Kubernetes Deployment

This repo serves as the kubernetes deployment for the platform.

## Overview

The platform consists of the following

- **resource** - shared resources that must be run to support the platform
- **service** - the micro services run on top of the shared infra as a platform
- **server.yaml** - experimental single yaml deployment of a self managed micro server for dev

## Dependencies

We have dependencies to get started

- Kubectl
- Helm

## Usage

For production

1. Spin up managed k8s on scaleway
2. Switch to the k8s env
3. ./install platform
3. Install as micro-secrets (auth keys, cf token)

## DNS Records

All 443 with certs managed by certmagic/acme/letsencrypt. Cloudflare used for DNS.

- api.m3o.com -> micro api
- web.m3o.com -> micro web
- proxy.m3o.com -> micro proxy

