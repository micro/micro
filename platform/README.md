# Platform

This directory serves as platform bootstrapping for Micro.

## Overview

The platform provides **Micro as a Service** as a fully managed solution. The platform is 
bootstrapped onto Kubernetes on any major cloud provider, including load balancing and 
dns management. This repository serves as the entrypoint and single location for all bootstrapping
related source code and documentation.

## Contents

- [kubernetes](kubernetes) - include the config to deploy to k8s
- [runbook](runbook) - a directory dedicated to platform operations

## TODO

- Add config - post deployment bootstrapping config
- Add command - turn into a `micro env {create, update, delete, list}` command
- Document the runbook - add a list of commands / docs / expected outcomes

### Updates

Each build of micro is tagged with a snapshot, e.g. `micro/micro:20200810104423b10609`. To update the platform
to use a new tag (with zero downtime), run the following command: 

```
kubectl set image deployments micro=micro/micro:20200810104423b10609 -l micro=runtime
```. 

The -l flag indicates we only want to do this to deployments with the micro=runtime label. 
The `micro=` part of the argument indicates the container name.

