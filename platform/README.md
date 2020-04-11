# Platform

The micro platform is a fully managed platform for microservices development.

## Overview

The platform provides **Micro as a Service** as a fully managed solution. The platform is 
bootstrapped onto Kubernetes on the major cloud providers, including load balancing and 
dns management. This repository serves as the entrypoint and single location for all platform related source 
code and documentation.

The platform builds on the [Micro](https://github.com/micro/micro) runtime and includes the features defined below.

## Features

The features which will be included in the platform

- **Cloud Automation** - Full terraform automation
- **Kubernetes Native** - Built to run on Kubernetes
- **Multi-Region** - Global deployments of the platform
- **Multi-Cloud** - Deploy across multiple clouds

## Usage

Install the platform binary

```
go get github.com/micro/platform
```

To bootstrap the platform, create a [config.yaml](./config-test.yaml), and prepare a AWS S3 bucket
for [terraform state storage](https://www.terraform.io/docs/backends/types/s3.html).

Then run

```
platform infra plan -c config.yaml
platform infra apply -c config.yaml
```

To destroy the cluster

```
platform infra destroy -c config.yaml
```

Configuration options can be set with viper, for example
[the state-store flag](https://github.com/micro/platform/blob/cc27173/cmd/infra.go#L44) can be set by
setting the environment variable `MICRO_STATE_STORE`.

See the [docs](docs) for more info.

