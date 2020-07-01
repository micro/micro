# Configuration

This document serves as the place for extended configuration

## Overview

The platform is automated through terraform and requires certain environmental config before it 
can be used, including the configuration for the underlying services and their access to resources 
like github and cloudflare.

## Dependencies

- Terraform >= 0.12.7
- Github

Image Pull Credentials: The default serviceaccount needs "Image pull secrets" set to a GitHub token.

## Usage

0. Obtain a working kubernetes cluster, and ensure that your default context is pointing to the cluster you want to deploy
  ```shell
  $ kubectl cluster-info
  Kubernetes master is running at https://127.0.0.1:46523
  KubeDNS is running at https://127.0.0.1:46523/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

  To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
  ```
1. Clone the micro respository
2. Create a `tfvars` file containing your configurations, see the [example](example.tfvars) 
3. Deploy the shared resources: 
  ```shell
  $ cd platform/kubernetes/micro-resource
  $ terraform apply -var-file=example.tfvars
  ```
4. Deploy the m3o platform
  ```shell
  $ cd platform/kubernetes/micro-platform
  $ terraform apply -var-file=example.tfvars
  ```

## Advanced usage:
The default configuration is that all m3o services talk to the infrastructure directly. These can be overridden using the per_service_overrides directive
