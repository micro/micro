# Micro Platform

This directory serves as platform bootstrapping for Micro.

## Overview

The platform provides **Micro as a Service** as a fully managed solution. The platform is 
bootstrapped onto Kubernetes on any major cloud provider, including load balancing and 
dns management. This repository serves as the entrypoint and single location for all bootstrapping
related source code and documentation.

## Contents

- [main.go](main.go) - The main program that includes the installer
- [kubernetes](kubernetes) - include the config to deploy to k8s

## Usage

To install the platform on an existing Kubernetes cluster use the following commands. 
The installer assumes the kubernetes directory is in the current directory.

```
platform install
```

To uninstall

```
platform uninstall
```

Each build of micro/platform is tagged with a snapshot, e.g. `micro/platform:20200810104423b10609`. To update the platform
to use a new tag (with zero downtime), run the following command: 

```
platform update 20200810104423b10609
```

If that fails use

```
kubectl set image deployments micro=micro/micro:20200810104423b10609 -l micro=runtime
```

The -l flag indicates we only want to do this to deployments with the micro=runtime label. 
The `micro=` part of the argument indicates the container name.

The platform binary bakes in micro with the platform profile. Type help for the micro commands.

```
platform --help
```

## TODO

- Add config - post deployment bootstrapping config
- Use the platform profile and define this as a micro platform binary
