# Kubernetes Deployment

This repo serves as the kubernetes deployment for the platform.

## Overview

The platform consists of the following

- **resource** - shared resources that must be run to support the platform
- **network** - the micro runtime run on top of the shared infra as distributed systems
- **server.yaml** - experimental single yaml deployment of a self managed micro server for dev

## Dependencies

We have dependencies to get started

- Kubectl
- Helm

## Usage

For production

1. Spin up managed k8s on scaleway
2. Spin up the shared infra in resource (./install.sh)
3. Get the secrets needed for micro and install
4. kubectl apply -f network

## DNS Records

All 443 with certs managed by certmagic/acme/letsencrypt. Cloudflare used for DNS.

- api.m3o.com -> micro api
- web.m3o.com -> micro web
- proxy.m3o.com => micro proxy

